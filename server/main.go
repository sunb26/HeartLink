// main.go file

package main

import (
	"context"                      // package is used to pass context between functions
	"errors"                       // package is used to handle errors
	"fmt"                          // package is used to print to console in a specific format
	"heartlinkServer/deviceServer" // import the deviceServer package from within the heartlinkServer project
	"heartlinkServer/endpoint2Pkg" // TESTING
	"log"                          // "log" package is used to log errors
	"net/http"                     // "net/http" package provides HTTP client/server implementations
)

func main() {

	log.Println("Begin main function")

	mux := http.NewServeMux() // create custom multiplexer to handle incoming requests

	// each individual HandleFunc is used to handle a specific endpoint
	mux.HandleFunc("/POSTRawAudioFile", deviceServer.POSTRawAudioFile)
	mux.HandleFunc("/ReceiveMultipartForm", deviceServer.ReceiveMultipartForm) // TESTING
	mux.HandleFunc("/endpoint1", endpoint2Pkg.GetEndpoint1)                    // TESTING
	mux.HandleFunc("/endpoint2_1", endpoint2Pkg.Endpoint2Function1)            // TESTING
	mux.HandleFunc("/endpoint2_2", endpoint2Pkg.Endpoint2Function2)            // TESTING

	ctx := context.Background()

	// define the http server
	server := &http.Server{
		Addr: ":80", // can use this as address to run locally or on hosted server
		// Addr:    "192.168.137.108:8080", // TESTING (for local network)
		Handler: mux,
	}

	error := server.ListenAndServe() // starts http server, saves any resulting errors

	// check if error is that server is closed or some other error
	if errors.Is(error, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if error != nil {
		fmt.Printf("error listening for server: %s\n", error)
	}

	<-ctx.Done() // waiting indefinitely for ctx to be cancelled (should never happen)

	log.Println("End main function") // this should never be printed

}
