package api

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync/atomic"

	"github.com/gorilla/mux"
	_ "github.com/nacuellar25111992/bricks-backend-poc/internal/api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// @title Bricks Backend POC API
// @version 2.0
// @description Go microservice template for Kubernetes.

// @contact.name Source Code
// @contact.url https://github.com/nacuellar25111992/bricks-backend-poc

// @license.name MIT License
// @license.url https://github.com/nacuellar25111992/bricks-backend-poc/blob/master/LICENSE

// @host localhost:9898
// @BasePath /
// @schemes http https

var (
	healthy int32
	ready   int32
)

type Server struct {
	router  *mux.Router
	logger  *zap.Logger
	config  *Config
	handler http.Handler
}

func NewServer(config *Config, logger *zap.Logger) *Server {

	srv := &Server{
		router: mux.NewRouter(),
		logger: logger,
		config: config,
	}

	return srv
}

func (s *Server) registerHandlers() {

	s.router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	s.router.HandleFunc("/", s.infoHandler).Methods("GET")
	s.router.HandleFunc("/version", s.versionHandler).Methods("GET")
	s.router.HandleFunc("/env", s.envHandler).Methods("GET", "POST")
	s.router.HandleFunc("/healthz", s.healthzHandler).Methods("GET")
	s.router.HandleFunc("/readyz", s.readyzHandler).Methods("GET")
	s.router.HandleFunc("/readyz/enable", s.enableReadyHandler).Methods("POST")
	s.router.HandleFunc("/readyz/disable", s.disableReadyHandler).Methods("POST")
	s.router.HandleFunc("/token", s.tokenGenerateHandler).Methods("POST")
	s.router.HandleFunc("/token/validate", s.tokenValidateHandler).Methods("GET")
	s.router.HandleFunc("/api/info", s.infoHandler).Methods("GET")
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	s.router.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		doc, err := swag.ReadDoc()
		if err != nil {
			s.logger.Error("swagger error", zap.Error(err), zap.String("path", "/swagger.json"))
		}
		w.Write([]byte(doc))
	})
}

func (s *Server) registerMiddlewares() {

	httpLogger := NewLoggingMiddleware(s.logger)

	s.router.Use(httpLogger.Handler)
	s.router.Use(versionMiddleware)
}

func (s *Server) ListenAndServe(stopCh <-chan struct{}) {

	s.registerHandlers()
	s.registerMiddlewares()

	if s.config.H2C {
		s.handler = h2c.NewHandler(s.router, &http2.Server{})
	} else {
		s.handler = s.router
	}

	s.printRoutes()

	// create the http server.
	srv := s.startServer()

	// signal Kubernetes the server is ready to receive traffic.
	if !s.config.Unhealthy {
		atomic.StoreInt32(&healthy, 1)
	}
	if !s.config.Unready {
		atomic.StoreInt32(&ready, 1)
	}

	// wait for SIGTERM or SIGINT.
	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), s.config.HttpServerShutdownTimeout)
	defer cancel()

	// all calls to /healthz and /readyz will fail from now on.
	atomic.StoreInt32(&healthy, 0)
	atomic.StoreInt32(&ready, 0)

	s.logger.Info("shutting down http server", zap.Duration("timeout", s.config.HttpServerShutdownTimeout))

	// determine if the http server was started.
	if srv != nil {
		if err := srv.Shutdown(ctx); err != nil {
			s.logger.Warn("http server graceful shutdown failed", zap.Error(err))
		}
	}
}

func (s *Server) startServer() *http.Server {

	// determine if the port is specified.
	if s.config.Port == "0" {

		// move on immediately
		return nil
	}

	srv := &http.Server{
		Addr:         s.config.Host + ":" + s.config.Port,
		WriteTimeout: s.config.HttpServerTimeout,
		ReadTimeout:  s.config.HttpServerTimeout,
		IdleTimeout:  2 * s.config.HttpServerTimeout,
		Handler:      s.handler,
	}

	// start the server in the background.
	go func() {
		s.logger.Info("starting http server.", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatal("http server crashed", zap.Error(err))
		}
	}()

	// return the server and routine.
	return srv
}

func (s *Server) printRoutes() {

	err := s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {

		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("route: ", pathTemplate)
		}

		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("path regexp: ", pathRegexp)
		}

		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("queries templates: ", strings.Join(queriesTemplates, ","))
		}

		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("queries regexps: ", strings.Join(queriesRegexps, ","))
		}

		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("methods: ", strings.Join(methods, ","))
		}

		fmt.Println()

		return nil
	})
	if err != nil {
		s.logger.Info("error printing routes", zap.Error(err))
	}
}

type ArrayResponse []string
type MapResponse map[string]string
