package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/contextart/al/api/db/queries"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Server struct {
	db  *pgxpool.Pool
	dbq *queries.Queries
}

func New(db *pgxpool.Pool) *Server {
	return &Server{
		db:  db,
		dbq: queries.New(db),
	}
}

func (s *Server) Handler(env, gitSha string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/tree", s.TreeHandler)
	mux.HandleFunc("/api/v1/proof", s.GetProof)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, gitSha)
	})

	h := http.Handler(mux)
	h = versionHandler(h, gitSha)
	h = tracingHandler(env, gitSha, h)
	h = hlog.NewHandler(log.Logger)(h)
	h = hlog.UserAgentHandler("user_agent")(h)
	h = hlog.RefererHandler("referer")(h)
	h = hlog.RequestIDHandler("req_id", "Request-Id")(h)
	h = hlog.URLHandler("path")(h)
	h = hlog.RequestHandler("req")(h)

	if env == "production" {
		return h
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	h = c.Handler(h)

	return h
}

func versionHandler(h http.Handler, sha string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("server-version", sha)
		h.ServeHTTP(w, r)
	})
}

func tracingHandler(env, sha string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ctx := tracer.StartSpanFromContext(r.Context(), fmt.Sprintf("req.%s", r.URL.Path))
		defer span.Finish()

		log := zerolog.Ctx(r.Context())
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Uint64("dd.trace_id", span.Context().TraceID())
		})
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			const serviceName = "allow-list-api"
			return c.Str("dd.service", serviceName)
		})
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("dd.env", env)
		})
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("dd.version", sha)
		})

		span.SetTag(ext.ResourceName, r.URL.Path)
		span.SetTag(ext.SpanType, ext.SpanTypeWeb)
		span.SetTag(ext.HTTPMethod, r.Method)

		sc := &statusCapture{ResponseWriter: w}
		h.ServeHTTP(sc, r.WithContext(ctx))
		span.SetTag(ext.HTTPCode, sc.status)
	})
}

type statusCapture struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
}

func (s *statusCapture) WriteHeader(c int) {
	s.status = c
	s.wroteHeader = true
	s.ResponseWriter.WriteHeader(c)
}

func (s *statusCapture) Write(b []byte) (int, error) {
	if !s.wroteHeader {
		s.WriteHeader(http.StatusOK)
	}
	return s.ResponseWriter.Write(b)
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
