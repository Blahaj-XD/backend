package main

import (
	"context"
	"os"

	"github.com/BlahajXD/backend/api"
	"github.com/BlahajXD/backend/backend"
	"github.com/BlahajXD/backend/config"
	"github.com/BlahajXD/backend/platform"
	"github.com/BlahajXD/backend/repo"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func init() {
	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.With().Caller().Stack().Logger()
	if os.Getenv("ENVIRONMENT") == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Starting HTTP Service")
	config.Init()
}

func main() {
	ctx := context.Background()

	postgres, err := platform.NewPostgreSQLClient(ctx, config.DatabaseURL())
	if err != nil {
		log.Fatal().Err(err).Msg("initialize postgres")
	}

	repo := repo.New(postgres)
	backend := backend.New(repo)
	api := api.New(backend)

	signal := <-api.Start(config.Host(), config.Port())
	api.Shutdown(ctx, signal)
}
