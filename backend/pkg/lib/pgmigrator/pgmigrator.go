package pgmigrator

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"univer/pkg/lib/log"
)

var (
	MigrationsSchemaName = "public"
	MigrationsTableName  = ""
)

func Run(name string, logger Logger, source fs.FS, database *url.URL) error {
	if name == "" {
		panic(fmt.Sprintf("pg migrator: %s: no name", name))
	}
	if logger == nil {
		panic(fmt.Sprintf("pg migrator: %s: nil logger", name))
	}
	if source == nil {
		panic(fmt.Sprintf("pg migrator: %s: nil source", name))
	}
	if database == nil {
		panic(fmt.Sprintf("pg migrator: %s: nil database", name))
	}

	sourceName := "iofs"

	sourceInstance, err := iofs.New(nocloseFS{source}, ".")
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSourceWrappingFailed, err)
	}
	var sourceInstanceIsManaged bool
	defer func() {
		if sourceInstanceIsManaged {
			return
		}

		err := sourceInstance.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("pg migrator: %s: source closing failed", name), log.ErrorAttr(err))
		}
	}()

	// There is no way to suppress closing of *sql.DB...

	databaseName := "pgx/v5"

	databaseClient, err := sql.Open(databaseName, database.String())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDatabaseOpeningFailed, err)
	}
	var databaseClientIsManaged bool
	defer func() {
		if databaseClientIsManaged {
			return
		}

		err := databaseClient.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("pg migrator: %s: database closing failed", name), log.ErrorAttr(err))
		}
	}()

	migrationsTableName := MigrationsTableName
	if migrationsTableName == "" {
		migrationsTableName = strings.ReplaceAll(name, "-", "_") + "_migrations"
	}

	databaseInstance, err := pgx.WithInstance(databaseClient, &pgx.Config{
		MigrationsTable: migrationsTableName,
		DatabaseName:    database.Path,
		SchemaName:      MigrationsSchemaName,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDatabaseWrappingFailed, err)
	}
	databaseClientIsManaged = true
	var databaseInstanceIsManaged bool
	defer func() {
		if databaseInstanceIsManaged {
			return
		}

		err := databaseInstance.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("pg migrator: %s: database closing failed", name), log.ErrorAttr(err))
		}
	}()

	m, err := migrate.NewWithInstance(sourceName, sourceInstance, databaseName, databaseInstance)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInstantiationFailed, err)
	}
	sourceInstanceIsManaged = true
	databaseInstanceIsManaged = true
	defer func() {
		sourceErr, databaseErr := m.Close()
		if sourceErr != nil {
			logger.Error(fmt.Sprintf("pg migrator: %s: source closing failed", name), log.ErrorAttr(sourceErr))
		}
		if databaseErr != nil {
			logger.Error(fmt.Sprintf("pg migrator: %s: database closing failed", name), log.ErrorAttr(databaseErr))
		}
	}()

	version, dirty, err := m.Version()
	if err != nil {
		if !errors.Is(err, migrate.ErrNilVersion) {
			return fmt.Errorf("%w: %w", ErrVersionLoadingFailed, err)
		}

		logger.Debug(fmt.Sprintf("pg migrator: %s: no version detected", name))
	} else {
		logger.Debug(fmt.Sprintf("pg migrator: %s: detected version is %s", name, formatVersion(version, dirty))) //nolint:perfsprint
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Debug(fmt.Sprintf("pg migrator: %s: already up to date", name))

			return nil
		}

		return fmt.Errorf("%w: %w", ErrMigrationFailed, err)
	}

	version, dirty, err = m.Version()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrVersionLoadingFailed, err)
	}
	logger.Debug(fmt.Sprintf("pg migrator: %s: updated version is %s", name, formatVersion(version, dirty))) //nolint:perfsprint

	return nil
}

type nocloseFS struct {
	fs.FS
}

func (nocloseFS) Close() error {
	return nil
}

func formatVersion(version uint, dirty bool) string {
	s := strconv.FormatUint(uint64(version), 10)
	if dirty {
		s += " (dirty)"
	}

	return s
}
