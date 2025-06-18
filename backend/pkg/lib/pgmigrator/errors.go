package pgmigrator

import (
	"errors"
)

// ErrDatabaseOpeningFailed indicates an error when opening database client.
var ErrDatabaseOpeningFailed = errors.New("pg migrator: database opening failed")

// ErrDatabaseWrappingFailed indicates an error when wrapping database client.
var ErrDatabaseWrappingFailed = errors.New("pg migrator: database wrapping failed")

// ErrInstantiationFailed indicates an error when instantiating migrator.
var ErrInstantiationFailed = errors.New("pg migrator: instantiation failed")

// ErrMigrationFailed indicates an error when migrating database.
var ErrMigrationFailed = errors.New("pg migrator: migration failed")

// ErrSourceWrappingFailed indicates an error when wrapping source file system.
var ErrSourceWrappingFailed = errors.New("pg migrator: source wrapping failed")

// ErrVersionLoadingFailed indicates an error when loading database version.
var ErrVersionLoadingFailed = errors.New("pg migrator: version loading failed")
