package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/contextart/al/api"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var GitSha string

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "processor error: %s", err)
		debug.PrintStack()
		os.Exit(1)
	}
}

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if env == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	var (
		logger = log.Logger.With().Caller().Logger()
		ctx    = logger.WithContext(context.Background())
	)
	const defaultPGURL = "postgres:///al"
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		dburl = defaultPGURL
	}
	dbc, err := pgxpool.ParseConfig(dburl)
	check(err)
	dbc.ConnConfig.LogLevel = pgx.LogLevelTrace
	dbc.MaxConns = 20
	db, err := pgxpool.ConnectConfig(ctx, dbc)
	check(err)

	s := api.New(db)

	const defaultListen = ":8080"
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = defaultListen
	}
	hs := &http.Server{
		Addr:    listen,
		Handler: s.Handler(env, GitSha),
	}
	log.Ctx(ctx).Info().Str("listen", listen).Str("git-sha", GitSha).Msg("http server")
	check(hs.ListenAndServe())
}
