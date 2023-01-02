package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mchmarny/disco/cmd/server/disco"
	"github.com/mchmarny/disco/pkg/metric"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	serviceName    = "disco"
	addressDefault = ":8080"
	logLevelEnvVar = "debug"

	closeTimeout = 3
	readTimeout  = 10
	writeTimeout = 600
)

var (
	// Version is set at build time.
	version = "v0.0.1-default"

	projectID  = os.Getenv("PROJECT_ID")
	bucketName = os.Getenv("GCS_BUCKET")
)

type key int

func main() {
	initLogging()
	log.Info().Msgf("starting %s server (version: %s, project: %s, bucket: %s)",
		serviceName, version, projectID, bucketName)

	if projectID == "" || bucketName == "" {
		log.Fatal().Msg("either PROJECT_ID or GCS_BUCKET env var not defined")
	}

	ctx := context.Background()
	counter, err := metric.NewAPICounter(projectID)
	if err != nil {
		log.Fatal().Msgf("error while creating counter: %v", err)
	}

	if err := counter.Count(ctx, metric.MakeMetricType("server/start"), 1, nil); err != nil {
		log.Fatal().Msgf("unable to write metrics: %v", err)
	}

	h, err := disco.NewHandler(version, bucketName, counter)
	if err != nil {
		log.Fatal().Msgf("error while creating event handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HandlerDefault)
	mux.HandleFunc("/disco", h.DiscoHandler)

	address := addressDefault
	if val, ok := os.LookupEnv("PORT"); ok {
		address = fmt.Sprintf(":%s", val)
	}

	run(ctx, mux, address)
}

var contextKey key

// run starts the server and waits for termination signal.
func run(ctx context.Context, mux *http.ServeMux, address string) {
	server := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: readTimeout * time.Second,
		WriteTimeout:      writeTimeout * time.Second,
		BaseContext: func(l net.Listener) context.Context {
			// adding server address to ctx handler functions receives
			return context.WithValue(ctx, contextKey, l.Addr().String())
		},
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("error listening for server: %v", err)
		}
	}()
	log.Debug().Msg("server started")

	<-done
	log.Debug().Msg("server stopped")

	downCtx, cancel := context.WithTimeout(ctx, closeTimeout*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(downCtx); err != nil {
		log.Fatal().Msgf("error shuting server down: %v", err)
	}
}

func initLogging() {
	level := zerolog.InfoLevel
	levStr := os.Getenv(logLevelEnvVar)
	if levStr == "true" {
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)

	out := zerolog.ConsoleWriter{
		Out: os.Stdout,
		PartsExclude: []string{
			zerolog.TimestampFieldName,
		},
	}

	log.Logger = zerolog.New(out)
}
