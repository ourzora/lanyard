package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/contextwtf/lanyard/api/tracing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/ryandotsmith/jh"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
)

type Server struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Server {
	return &Server{
		db: db,
	}
}

func jsonHandler(f any) http.Handler {
	h, err := jh.Handler(f, jh.ErrHandler)
	if err != nil {
		panic(fmt.Sprintf("setting up handler: %s", err))
	}
	return h
}

func apiError(c int, m string) error {
	return jh.Error{Code: c, Message: m}
}

func (s *Server) treeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			jsonHandler(s.CreateTree).ServeHTTP(w, r)
		case http.MethodGet:
			jsonHandler(s.GetTree).ServeHTTP(w, r)
		}
	})
}

func (s *Server) Handler(env, gitSha string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/v1/tree", s.treeHandler())
	mux.Handle("/api/v1/proof", jsonHandler(s.GetProof))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, gitSha)
	})

	h := http.Handler(mux)
	h = versionHandler(h, gitSha)
	h = tracingHandler(os.Getenv("DD_ENV"), os.Getenv("DD_SERVICE"), gitSha, h)
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

func tracingHandler(env, service, sha string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ctx := tracing.SpanFromContext(r.Context(), "http.request")

		defer span.Finish()

		log := zerolog.Ctx(r.Context())
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Uint64("dd.trace_id", span.Context().TraceID())
		})
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("dd.service", service)
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
