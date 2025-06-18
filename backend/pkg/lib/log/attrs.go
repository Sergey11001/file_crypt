package log

import (
	"log/slog"
)

func ErrorAttr(value error) slog.Attr {
	var s string
	if value == nil {
		s = "<nil>"
	} else {
		s = value.Error()
	}

	return slog.String("error", s)
}
