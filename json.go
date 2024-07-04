package render

import (
	"encoding/json"
	"io"
	"net/http"
)

// JSON encodes the provided data as JSON and writes it to the provided io.Writer.
//
// This method takes an io.Writer (typically an http.ResponseWriter), an HTTP status code,
// and the data to be encoded as JSON. It writes the JSON-encoded data to the writer and sets
// the HTTP status code on the http.ResponseWriter.
//
// Parameters:
// - w: io.Writer to which the JSON-encoded data will be written. This is usually an http.ResponseWriter.
// - status: HTTP status code to set on the http.ResponseWriter.
// - data: The data to be JSON-encoded and written to the writer.
//
// Example usage:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    data := map[string]string{
//	        "message": "Hello, World!",
//	    }
//	    JSON(w, http.StatusOK, data)
//	}
func JSON(w io.Writer, status int, data any) {
	if hw, ok := w.(http.ResponseWriter); ok {
		if b, err := json.Marshal(data); err != nil {
			http.Error(hw, err.Error(), http.StatusInternalServerError)
		} else {
			hw.Header().Set("Content-Type", "application/json")
			hw.WriteHeader(status)
			hw.Write(b)
		}
	}
}
