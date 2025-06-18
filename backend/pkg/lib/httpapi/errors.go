package httpapi

import (
	"net/http"
	"univer/pkg/lib/httpserver"

	"univer/pkg/lib/errs"
)

type Error struct {
	Code    string             `json:"code"`
	Details *map[string]string `json:"details,omitempty"`
	Message string             `json:"message"`
}

type ErrorResult struct {
	Error Error `json:"error"`
}

type SuccessResult struct {
	Success bool `json:"success"`
}

func WriteError(w http.ResponseWriter, err error) {
	class, code, message, details := errs.Parse(err)

	httpserver.WriteJSONResponse(w, errs.HTTPStatus(class), ErrorResult{
		Error: Error{
			Code:    code,
			Message: message,
			Details: &details,
		},
	})
}

func PropagateError(res *http.Response, body *ErrorResult) error {
	if res == nil {
		panic("http api: nil response")
	}
	if body == nil {
		panic("http api: nil body")
	}

	var detailers []errs.Detailer
	if body.Error.Details != nil && len(*body.Error.Details) > 0 {
		detailers = []errs.Detailer{errs.Details(*body.Error.Details)}
	}

	return errs.HTTPClass(res.StatusCode).New(body.Error.Code, body.Error.Message, detailers...)
}
