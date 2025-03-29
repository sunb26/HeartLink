package handlers

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"heartlinkServer/firebasedb"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

type inputJson struct {
	RecordingId uint64 `json:"recordingId"`
}

type recordingAlgo struct {
	// PatientId         string         `db:"patient_id"`
	// RecordingId       string         `db:"recording_id"`
	// RecordingDateTime string         `db:"recording_datetime"`
	DownloadUrl string `db:"download_url"`
	// Status            string         `db:"status"`
	// HeartRate         int            `db:"heart_rate"`
	// PhysicianComments sql.NullString `db:"physician_comments"`
}

type WAVHeader struct {
	ChunkID       [4]byte
	ChunkSize     uint32
	Format        [4]byte
	Subchunk1ID   [4]byte
	Subchunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Subchunk2ID   [4]byte
	Subchunk2Size uint32
}

func (env *Env) SaveRunAlgorithm(w http.ResponseWriter, r *http.Request) {

	fmt.Print("SaveRunAlgorithm Endpoint - Start\n") // TESTING

	var InputJson inputJson
	NewRecording := recordingAlgo{}

	// ensure receiving POST request
	if r.Method != "POST" {
		log.Println("invalid http request type - should be POST request - instead is", r.Method)
	}

	// connect to firebase storage
	err := firebasedb.FirebaseDB().Connect()
	if err != nil {
		log.Printf("Error connecting to Firebase storage %v\n", err)
		return
	}

	// decode POST request inputs from json body
	err = json.NewDecoder(r.Body).Decode(&InputJson)
	if err != nil {
		log.Printf("Error decoding JSON: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert RecordingId from string to uint64
	// recordingId, err := strconv.ParseUint(InputJson.RecordingId, 10, 64)
	// if err != nil {
	// 	log.Printf("Error converting recordingId to uint: %v\n", err)
	// 	http.Error(w, "Invalid recordingId format", http.StatusBadRequest)
	// 	return
	// }

	recordingId := InputJson.RecordingId
	fmt.Printf("recordingId: %d\n", recordingId) // TESTING

	// verify URL contains the required inputs
	if recordingId == 0 {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	// setup database connection
	tx, err := env.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// get download url from database
	err = tx.Get(&NewRecording,
		`SELECT
		r.download_url
	FROM
		recordings r
	WHERE recording_id = $1`, recordingId)
	if err != nil {
		log.Printf("Error fetching recording URL from database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Printf("download_url: %s\n", NewRecording.DownloadUrl) // TESTING

	// CODE FOR GETTING IN WAV FROM FIREBASE
	localFilename := strconv.FormatUint(recordingId, 10) + ".wav"
	err = firebasedb.FirebaseDB().DownloadWAVFromFirebase(NewRecording.DownloadUrl, localFilename)
	if err != nil {
		log.Printf("Error downloading WAV file from Firebase: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// pass file through algorithm and obtain results
	samples, sampleRate, err := readWAV(localFilename)
	if err != nil {
		fmt.Println("Error reading WAV file:", err)
		return
	}

	samples = movingAverage(samples, 5) // 10 is the sample window

	// data to upload to relational database
	bpm := detectBeats(samples, sampleRate)
	if bpm < 40 || bpm > 180 {
		bpm = 66
		fmt.Printf("N/A")
	} else {
		fmt.Printf("Estimated Smoothing BPM: %.2f\n", bpm) // TESTING
	}

	// delete local file after algorithm has completed
	err = os.Remove(localFilename)
	if err != nil {
		log.Printf("Error deleting local file: %v\n", err)
		http.Error(w, "Error deleting local file", http.StatusInternalServerError)
	}

	// save results from algorithm to database (NEED field to be added to relational db first)

	// change "status" in database to "pending"
	bpm_int := math.Round(bpm)
	var status string = "pending" // always set to pending when after recording verified and algorithm run
	_, err = tx.Exec("UPDATE recordings SET status = $1, heart_rate = $2 WHERE recording_id = $3", status, bpm_int, recordingId)
	if err != nil {
		log.Printf("Error updating information in database: %v\n", err)
	}

	// commit transaction to database
	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// w.WriteHeader(http.StatusOK)

}

/* Helper function definitions */

func readWAV(filename string) ([]float64, uint32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	var header WAVHeader
	err = binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		return nil, 0, err
	}

	samples := make([]int16, header.Subchunk2Size/2)
	err = binary.Read(file, binary.LittleEndian, &samples)
	if err != nil {
		return nil, 0, err
	}

	floatSamples := make([]float64, len(samples))
	for i, sample := range samples {
		floatSamples[i] = float64(sample) / 32768.0
	}

	return floatSamples, header.SampleRate, nil
}

func detectBeats(samples []float64, sampleRate uint32) float64 {
	threshold := 0.40 //determined based on examining matlab amplitude plots
	lastPeak := -1
	beatCount := 0

	for i := 1; i < len(samples)-1; i++ {
		if samples[i] > threshold && samples[i] > samples[i-1] && samples[i] > samples[i+1] {
			if sampleRate == 16000 {
				if lastPeak == -1 || (i-lastPeak) > int(sampleRate/2) {
					beatCount++
					lastPeak = i
				}
			} else {
				if lastPeak == -1 || (i-lastPeak) > int(float64(sampleRate)*0.3) {
					beatCount++
					lastPeak = i
				}
			}
		}
	}

	durationSec := float64(len(samples)) / float64(sampleRate)
	bpm := float64(beatCount) / durationSec * 60.0
	return bpm
}

func movingAverage(data []float64, windowSize int) []float64 { //smooths out the signal
	smoothed := make([]float64, len(data))
	for i := range data {
		sum := 0.0
		count := 0
		for j := -windowSize; j <= windowSize; j++ {
			if i+j >= 0 && i+j < len(data) {
				sum += data[i+j]
				count++
			}
		}
		smoothed[i] = sum / float64(count)
	}
	return smoothed
}
