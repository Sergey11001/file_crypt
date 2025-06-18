package httpapi

import (
	"net/http"

	"univer/pkg/lib/httpserver"
)

// WriteResult writes the specified result.
//
// WriteResult panics if w is nil.
func WriteResult(w http.ResponseWriter, result any) {
	if w == nil {
		panic("http api: nil response writer")
	}

	if result == nil {
		w.WriteHeader(http.StatusOK)

		return
	}

	httpserver.WriteJSONResponse(w, http.StatusOK, result)
}
