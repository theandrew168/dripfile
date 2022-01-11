package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/flow"
	"github.com/coreos/go-systemd/daemon"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/static"
	"github.com/theandrew168/dripfile/internal/web"
)

func main() {
	// log everything to stdout, include file name and line number
	logger := log.New(os.Stdout, "", 0)

	// check for config file flag
	conf := flag.String("conf", "dripfile.conf", "app config file")
	flag.Parse()

	// load user-defined config (if specified), else use defaults
	cfg, err := config.ReadFile(*conf)
	if err != nil {
		logger.Fatalln(err)
	}

	// open a database connection pool
	conn, err := pgxpool.Connect(context.Background(), cfg.DatabaseURI)
	if err != nil {
		logger.Fatalln(err)
	}
	defer conn.Close()

	// test connection to ensure all is well
	if err = conn.Ping(context.Background()); err != nil {
		logger.Fatalln(err)
	}

	mux := flow.New()

	// handle top-level special cases
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	})
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/webp")
		w.Write(static.Favicon)
	})
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(static.Robots)
	})

	// static files app
	staticApp := static.NewApplication(cfg, logger)
	mux.Handle("/static/...", http.StripPrefix("/static", staticApp.Router()))

	// rest api app
	//	apiApp := api.NewApplication(cfg, storage, logger)
	//	mux.Handle("/api/v1/...", http.StripPrefix("/api/v1", apiApp.Router()))

	// primary web app (last due to being a top-level catch-all)
	webApp := web.NewApplication(cfg, nil, logger)
	mux.Handle("/...", webApp.Router())

	addr := fmt.Sprintf("127.0.0.1:%s", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// open up the socket listener
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalln(err)
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	logger.Printf("started server on %s\n", addr)

	// kick off a goroutine to listen for SIGINT and SIGTERM
	shutdownError := make(chan error)
	go func() {
		// idle until a signal is caught
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		// give the web server 5 seconds to shutdown gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// shutdown the web server and track any errors
		logger.Println("stopping server")
		srv.SetKeepAlivesEnabled(false)
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		shutdownError <- nil
	}()

	// serve the app, check for ErrServerClosed (expected after shutdown)
	err = srv.Serve(l)
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalln(err)
	}

	// check for shutdown errors
	err = <-shutdownError
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Println("stopped server")
}
