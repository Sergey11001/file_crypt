package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// WriteJSONResponse writes JSON response with content specified by v and the specified status code.
//
// WriteJSONResponse panics if w is nil.
func WriteJSONResponse(w http.ResponseWriter, statusCode int, v any) {
	if w == nil {
		panic("http server: nil response writer")
	}

	data, err := json.Marshal(v)
	if err != nil {
		statusCode = http.StatusInternalServerError
		http.Error(w, fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)), statusCode)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(statusCode)
	_, _ = w.Write(data)
}
