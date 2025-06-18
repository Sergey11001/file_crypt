package pgclient

import (
	"net/url"
)

// NewURL initializes and returns a new connection URL.
func NewURL(config Config) (*url.URL, error) {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.User, config.Password),
		Host:   config.Addr,
		Path:   config.Database,
	}

	q := make(url.Values)
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()

	return u, nil
}
