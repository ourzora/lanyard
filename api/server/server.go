package server

import (
	"encoding/json"
	"net/http"

	"github.com/contextart/al/api/config"
	"github.com/contextart/al/api/db/queries"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	settings *config.Settings
	writeDB  *pgxpool.Pool
	writeQ   *queries.Queries
}

func Open(settings *config.Settings, writeDB *pgxpool.Pool, writeQ *queries.Queries) *Server {
	return &Server{
		settings: settings,
		writeDB:  writeDB,
		writeQ:   writeQ,
	}
}

func (s *Server) InstallMiddleware(next http.Handler) http.Handler {
	c := alice.New()

	c = c.Append(hlog.NewHandler(log.Logger))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))
	c = c.Append(hlog.URLHandler("path"))
	c = c.Append(hlog.RequestHandler("req"))
	c = c.Append(withVersion(s.settings.GitSHA))
	c = c.Append(s.withTraces(s.settings.GitSHA))
	if s.settings.CurrentEnv != "production" {
		c = c.Append(cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowCredentials: true,
		}).Handler)
	}

	return c.Then(next)
}

func (s *Server) sendJSON(r *http.Request, w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func (s *Server) sendJSONError(
	r *http.Request,
	w http.ResponseWriter,
	err error,
	code int,
	customMessage string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err != nil {
		log.Ctx(r.Context()).Err(err).Send()
	}

	message := http.StatusText(code)
	if customMessage != "" {
		message = customMessage
	}

	json.NewEncoder(w).Encode(map[string]any{
		"error":   true,
		"message": message,
	})
}

func (s *Server) withTraces(GitSHA string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := zerolog.Ctx(r.Context())
			log.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("dd.service", s.settings.DatadogService)
			})
			log = zerolog.Ctx(r.Context())
			log.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("dd.env", s.settings.DatadogEnvironment)
			})
			log = zerolog.Ctx(r.Context())
			log.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("dd.version", GitSHA)
			})

			next.ServeHTTP(w, r)
		})
	}
}

func withVersion(GitSHA string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("server-version", GitSHA)
			next.ServeHTTP(w, r)
		})
	}
}
