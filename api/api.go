package api

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/BlahajXD/backend/backend"
	"github.com/BlahajXD/backend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

type Server struct {
	app     *fiber.App
	backend *backend.Dependency
}

func New(backend *backend.Dependency) *Server {
	server := &Server{
		app: fiber.New(fiber.Config{
			AppName:       config.AppName(),
			WriteTimeout:  30 * time.Second,
			ReadTimeout:   30 * time.Second,
			ErrorHandler:  Error,
			CaseSensitive: true,
		}),
		backend: backend,
	}

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPatch,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodHead,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders: strings.Join([]string{
			fiber.HeaderAuthorization,
			fiber.HeaderContentType,
			fiber.HeaderAccept,
		}, ","),
	})

	server.app.Use(helmet.New())
	server.app.Use(recover.New())
	server.app.Use(corsMiddleware)

	server.app.Get("/", server.Health)

	return server
}

func (s *Server) Start(host, port string) <-chan os.Signal {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		address := net.JoinHostPort(host, port)
		log.Info().Msgf("Listening on %s", address)
		err := s.app.Listen(address)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	return exitSignal
}

func (s *Server) Shutdown(ctx context.Context, signal os.Signal) {
	timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	shutdownChan := make(chan error, 1)

	go func() {
		log.Warn().Any("signal", signal.String()).Msg("received signal, shutting down...")
		shutdownChan <- s.app.Shutdown()
	}()

	select {
	case <-timeout.Done():
		log.Warn().Msg("shutdown timed out, forcing exit")
		os.Exit(1)
	case err := <-shutdownChan:
		if err != nil {
			log.Fatal().Err(err).Msg("there was an error shutting down")
		} else {
			log.Info().Msg("shutdown complete")
		}
	}
}
