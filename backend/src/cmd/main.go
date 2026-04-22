package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/OGZKTeBmj/forum/internal/app"
	"github.com/OGZKTeBmj/forum/internal/config"
	"github.com/OGZKTeBmj/forum/internal/handler"
	"github.com/OGZKTeBmj/forum/internal/service"
	"github.com/OGZKTeBmj/forum/internal/storage/postgres"
	"github.com/OGZKTeBmj/forum/internal/storage/s3storage"
	"github.com/OGZKTeBmj/forum/utils"
	"github.com/OGZKTeBmj/forum/utils/flagandenv"
)

func main() {
	//Config and logger
	getter := flagandenv.EnvGetter{}

	authJWTSecret := getter.Get("AUTH_JWT_SECRET")

	pgsUser := getter.Get("POSTGRES_USER")
	pgsPassword := getter.Get("POSTGRES_PASSWORD")
	pgsDBName := getter.Get("POSTGRES_DB_NAME")

	s3Endpoint := getter.Get("S3_ENDPOINT_URL")
	s3Region := getter.Get("S3_REGION")
	s3Bucket := getter.Get("S3_BUCKET")
	s3AccessKey := getter.Get("S3_ACCESS_KEY")
	s3SecretKey := getter.Get("S3_SECRET_KEY")
	s3PublicBase := getter.Get("S3_PUBLIC_BASE")

	if err := getter.EmptiesValues(); err != nil {
		panic(err)
	}

	parser := flagandenv.NewFlagParser()

	cfgPath := parser.String("cfg-path", "", "path to config file")

	if err := parser.Parse(); err != nil {
		panic(err)
	}

	cfg := config.MustLoad(*cfgPath)
	log := utils.SetupLoger(cfg.Env)
	ctx := context.Background()

	//Postgres Storage
	postgresStorage := postgres.Storage{}
	postgresStorage.MustConnect(ctx, cfg.PostgresURL, pgsUser, pgsPassword, pgsDBName)

	defer postgresStorage.Stop(ctx)

	if err := postgresStorage.Init(ctx); err != nil {
		log.Error("postgres init error", utils.SlogErr(err))
		panic(err)
	}

	//S3 Storage
	s3Storage, err := s3storage.New(ctx, s3Endpoint, s3Region, s3Bucket, s3AccessKey, s3SecretKey, s3PublicBase)
	if err != nil {
		log.Error("s3 init error", utils.SlogErr(err))
	}

	//Services
	postsService := &service.PostsService{
		Provider: &postgresStorage,
		Log:      log,
	}

	authService := &service.AuthService{
		UserProvider: &postgresStorage,
		TokenTTL:     cfg.AuthTokenTTL,
		Log:          log,
		Secret:       []byte(authJWTSecret),
	}

	imageService := &service.ImageService{
		Log:      log,
		Provider: s3Storage,
	}

	//Handler
	handler := handler.New(
		log,
		authService,
		postsService,
		imageService,
	)
	handler.Init()

	//Application run
	application := app.New(
		handler,
		log,
		cfg.HTTPPort,
	)
	go application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("application is stopping", slog.String("Signal", sign.String()))

	postgresStorage.Stop(ctx)

	log.Info("application stopped")
}
