package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/oapi-codegen/runtime/types"

	"univer/pkg/lib/errs"
)

func handler(f func(r *http.Request) (any, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if w == nil {
			panic("http api: nil response writer")
		}
		if r == nil {
			panic("http api: nil request")
		}

		result, err := f(r)
		if err != nil {
			WriteError(w, err)

			return
		}

		WriteResult(w, result)
	})
}

func presentJSONDecodingError(err error) error {
	details := errs.Details{
		"cause": err.Error(),
	}

	if e := (*json.SyntaxError)(nil); errors.As(err, &e) {
		return errs.Invalid.New("BadRequest", "invalid syntax", details)
	}
	if e := (*json.UnmarshalTypeError)(nil); errors.As(err, &e) {
		return errs.Invalid.New("BadRequest", "invalid type", details)
	}

	return errs.Invalid.New("BadRequest", "invalid body", details)
}

// Handler returns a new handler calling the specified function to handle requests without body.
//
// Handler panics if f is nil.
func Handler(f func(ctx context.Context) (any, error)) http.Handler {
	if f == nil {
		panic("http api: nil function")
	}

	return handler(func(r *http.Request) (any, error) {
		return f(r.Context())
	})
}

// HandlerWithInput returns a new handler calling the specified function to handle requests with JSON body of type T.
//
// HandlerWithInput panics if f is nil.
func HandlerWithInput[T any](f func(ctx context.Context, input T) (any, error)) http.Handler {
	if f == nil {
		panic("http api: nil function")
	}

	return handler(func(r *http.Request) (any, error) {
		if !strings.Contains(r.Header.Get("Content-Type"), "json") { // as in generated clients
			return nil, errs.Invalid.New("BadRequest", "invalid content type")
		}

		var input T
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			return nil, presentJSONDecodingError(err)
		}

		return f(r.Context(), input)
	})
}

// HandlerWithForm returns a new handler calling the specified function to handle requests with multipart/form-data body of type T.
//
// HandlerWithForm panics if f is nil.
func HandlerWithForm[T any](f func(ctx context.Context, form T) (any, error)) http.Handler {
	if f == nil {
		panic("http api: nil function")
	}

	return handler(func(r *http.Request) (any, error) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			return nil, errs.Invalid.New("BadRequest", "invalid content type")
		}

		var form T
		err := parseForm(r, &form)
		if err != nil {
			return nil, err
		}

		return f(r.Context(), form)
	})
}

//nolint:cyclop
func parseForm(r *http.Request, v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Ptr && reflect.TypeOf(v).Elem().Kind() == reflect.Interface {
		return nil
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return errs.Invalid.New("BadRequest", "multipart form parsing failed", errs.Details{
			"cause": err.Error(),
		})
	}

	val := reflect.ValueOf(v).Elem()
	for i := range val.NumField() {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Tag.Get("json")

		//nolint:nestif
		if formValues, ok := r.MultipartForm.Value[fieldName]; ok {
			if field.Kind() == reflect.String {
				field.SetString(formValues[0])
			}
			if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
				err := parseFormArray(field, formValues[0])
				if err != nil {
					return err
				}
			}
			if field.Kind() == reflect.Struct {
				err := parseFormStruct(field, formValues[0])
				if err != nil {
					return err
				}
			}
		}
		if fileHeaders, ok := r.MultipartForm.File[fieldName]; ok {
			if len(fileHeaders) > 0 {
				f := types.File{}
				f.InitFromMultipart(fileHeaders[0])
				field.Set(reflect.ValueOf(f))
			}
		}
	}

	return nil
}

func parseFormStruct(dst reflect.Value, src string) error {
	ptr := reflect.New(dst.Type())

	err := json.Unmarshal([]byte(src), ptr.Interface())
	if err != nil {
		return presentJSONDecodingError(err)
	}

	dst.Set(ptr.Elem())

	return nil
}

func parseFormArray(dst reflect.Value, src string) error {
	typ := dst.Type().Elem()
	if typ.Kind() == reflect.Struct {
		ptr := reflect.New(reflect.SliceOf(typ))

		err := json.Unmarshal([]byte(src), ptr.Interface())
		if err != nil {
			return presentJSONDecodingError(err)
		}

		dst.Set(ptr.Elem())
	}

	return nil
}
