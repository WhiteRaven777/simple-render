package render

import (
	"encoding/json"
	"io"
	"net/http"
)

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
