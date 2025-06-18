package pgclient

type Config struct {
	Addr     string `default:"localhost:5432"`
	User     string `default:"postgres"`
	Password string `default:"postgres"`
	Database string `default:"postgres"`
}
