package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"univer/pkg/lib/redisclient"

	"univer/internal/httpengine"
	"univer/internal/pgmigrations"
	"univer/internal/service"
	"univer/pkg/auth"
	"univer/pkg/hasher"
	"univer/pkg/lib/config"
	"univer/pkg/lib/httpserver"
	"univer/pkg/lib/log"
	"univer/pkg/lib/pgclient"
	"univer/pkg/lib/pgmigrator"
	"univer/pkg/lib/runner"
	"univer/pkg/lib/s3client"
)

type Config struct {
	Logger             log.LoggerConfig
	HTTPServer         httpserver.Config
	PgClient           pgclient.Config
	RedisClient        redisclient.Config
	S3Buckets          []string `default:"filecrypto"`
	S3Client           s3client.ClientConfig
	Salt               string
	SingingKey         string
	UsersServiceConfig service.UsersConfig
}

func main() {
	err := run()
	if err != nil {
		_, _ = fmt.Printf("%+v", err) //nolint:forbidigo
		os.Exit(1)
	}
}

func run() (err error) {
	cfg, err := config.New[Config]()
	if err != nil {
		return err
	}

	logger, err := log.New(cfg.Logger)
	if err != nil {
		return err
	}

	runn, err := runner.NewRunner(context.Background(), logger)
	if err != nil {
		return err
	}

	httpRouter, err := httpserver.NewRouter()
	if err != nil {
		return err
	}

	httpServer, err := httpserver.New(cfg.HTTPServer, logger, httpRouter)
	if err != nil {
		return err
	}

	terminateHTTPServer := runn.RunModule(httpServer)
	defer terminateHTTPServer()

	err = runPgMigrator(cfg.PgClient, "univer", pgmigrations.FS, logger)
	if err != nil {
		return err
	}

	pgClient, err := pgclient.New(cfg.PgClient, logger)
	if err != nil {
		return err
	}
	defer pgClient.Close()

	redisClient, err := redisclient.New(cfg.RedisClient)
	if err != nil {
		return err
	}

	s3Client, err := s3client.New(cfg.S3Client, logger)
	if err != nil {
		return err
	}

	err = s3Client.InitBuckets(cfg.S3Buckets)
	if err != nil {
		return err
	}

	hashManager := hasher.NewSHA1Hasher(cfg.Salt)
	tokenManager, err := auth.NewManager(cfg.SingingKey)
	if err != nil {
		return err
	}

	usersService := service.NewUsersService(
		cfg.UsersServiceConfig,
		pgClient,
		redisClient,
		hashManager,
		tokenManager,
	)

	fileService := service.NewFilesService(pgClient, s3Client)

	httpEngine, err := httpengine.New(
		logger,
		fileService,
		usersService,
		tokenManager,
	)
	if err != nil {
		return err
	}

	httpRouter.Mount(httpEngine)

	return runn.Listen()
}

func runPgMigrator(pgConfig pgclient.Config, name string, fs embed.FS, logger *slog.Logger) error {
	logger.Debug(fmt.Sprintf("pg migrator: %s: started", name))
	defer logger.Debug(fmt.Sprintf("pg migrator: %s: finished", name))

	pgURL, err := pgclient.NewURL(pgConfig)
	if err != nil {
		return err
	}

	return pgmigrator.Run(name, logger, fs, pgURL)
}
