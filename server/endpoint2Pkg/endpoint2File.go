// endpoint2File.go file

package endpoint2Pkg

import (
	"fmt"
	"io"
	"net/http"
)

func Endpoint2Function1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "heartlinkServer endpoint2 (function 1)\n")
	io.WriteString(w, "Written via io.WriteString function\n")
}

func Endpoint2Function2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "heartlinkServer endpoint2 (function 2)")
}

// GetEndpoint1 function - TESTING
func GetEndpoint1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Start heartlinkServer endpoint1\n\n") // Write start message to response writer "w"

	// r.URL.Query used to access the query parameters in the URL
	hasFirst := r.URL.Query().Has("first") // the Has method checks if a query parameter exists (returns Bool)
	first := r.URL.Query().Get("first")    // the Get method retrieves the value of a query parameter (returns String)
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	// Printf writes the arguments extracted from URL to the console
	fmt.Printf("getEndpoint1 arguments: first(%t)=%s, second(%t)=%s\n",
		// ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second)

	// Fprintf writes the arguments extracted from URL to the response writer "w"
	fmt.Fprintf(w, "getEndpoint1 arguments: first(%t)=%s, second(%t)=%s\n",
		// ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second)

	fmt.Fprint(w, "End heartlinkServer endpoint1\n") // Write end message to response writer "w"
}
