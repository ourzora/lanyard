package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/contextart/al/api/config"
	"github.com/contextart/al/api/db"
	"github.com/contextart/al/api/db/queries"
	"github.com/contextart/al/api/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	settings := config.Init()

	ctx := setupLogger(settings)

	writeDB, err := db.WithConfig(ctx, settings.PostgresURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	writeQ := queries.New(writeDB)

	s := server.Open(settings, writeDB, writeQ)

	r := muxtrace.NewRouter()
	r.HandleFunc("/api/v1/merkle", s.CreateTree).Methods("POST")
	r.HandleFunc("/api/v1/merkle/{root}", s.RetrieveTree).Methods("GET")
	r.HandleFunc("/api/v1/merkle/{root}/proof/{address}", s.RetrieveProof).Methods("GET")
	r.HandleFunc("/health", s.Health).Methods("GET")

	router := &http.Server{
		Addr:    ":" + settings.Port,
		Handler: s.InstallMiddleware(r),
	}

	go func() {
		log.Info().Str("port", settings.Port).Msgf("HTTP server started")

		err := router.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Err(err).Msg("Error starting server")
			return
		}
	}()

	<-sigint
	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	go func() {
		<-sigint
		os.Exit(1)
	}()

	router.Shutdown(ctx)
	tracer.Stop()

	defer cancel()
}

func setupLogger(settings *config.Settings) context.Context {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var (
		logger = log.Logger.With().Logger()
		ctx    = logger.WithContext(context.Background())
	)

	if settings.CurrentEnv == "" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return ctx
}
