// main.go file

package main

import (
	"context"
	"errors"
	"fmt"
	"heartlinkServer/deviceServer"
	"heartlinkServer/endpoint2Pkg" // TESTING
	"log"
	"net/http"
)

func main() {

	log.Println("Server Running")

	mux := http.NewServeMux() // create custom multiplexer to handle incoming requests

	// each HandleFunc is used to handle a specific endpoint
	mux.HandleFunc("/POSTRawAudioFile", deviceServer.POSTRawAudioFile)
	mux.HandleFunc("/endpoint1", endpoint2Pkg.GetEndpoint1)         // TESTING
	mux.HandleFunc("/endpoint2_1", endpoint2Pkg.Endpoint2Function1) // TESTING
	mux.HandleFunc("/endpoint2_2", endpoint2Pkg.Endpoint2Function2) // TESTING

	ctx := context.Background()

	// define http server
	server := &http.Server{
		Addr: ":8080", // address to run on localHost or on google cloud server
		// Addr:    "192.168.:8080", // connect to local IP network
		Handler: mux,
	}

	error := server.ListenAndServe() // starts http server
	if errors.Is(error, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if error != nil {
		fmt.Printf("error listening for server: %s\n", error)
	}

	<-ctx.Done() // waiting indefinitely for context to be cancelled (should never happen)

	log.Println("Server Ended") // this should never be printed

}
