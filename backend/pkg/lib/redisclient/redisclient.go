package redisclient

import (
	"github.com/redis/go-redis/v9"
)

// Config configures client.
type Config struct {

	// Addr in form "host:port" to connect to.
	Addr string `default:"localhost:6379"`

	// User to authenticate with.
	User string `default:""`

	// Password to authenticate with.
	Password string `default:""`

	// Database to be selected after connecting.
	Database int `default:"0"`
}

// New initializes and returns a new client.
func New(config Config) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Username: config.User,
		Password: config.Password,
		DB:       config.Database,
	}), nil
}
