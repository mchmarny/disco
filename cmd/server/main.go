package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mchmarny/disco/cmd/server/disco"
	"github.com/mchmarny/disco/pkg/metric"
)

const (
	serviceName    = "disco"
	addressDefault = ":8080"

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
	log.SetFlags(log.Lshortfile)
	log.Printf("starting %s server (%s)...\n", serviceName, version)

	if projectID == "" || bucketName == "" {
		log.Fatal("either PROJECT_ID or GCS_BUCKET env var not defined")
	}

	ctx := context.Background()
	counter, err := metric.NewAPICounter(projectID)
	if err != nil {
		log.Fatalf("error while creating counter: %v", err)
	}

	if err := counter.Count(ctx, metric.MakeMetricType("server/start"), 1, nil); err != nil {
		log.Printf("unable to write metrics: %v", err)
	}

	h, err := disco.NewHandler(bucketName, counter)
	if err != nil {
		log.Fatalf("error while creating event handler: %v", err)
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
			log.Printf("error listening for server: %v\n", err)
		}
	}()
	log.Print("server started")

	<-done
	log.Print("server stopped")

	downCtx, cancel := context.WithTimeout(ctx, closeTimeout*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(downCtx); err != nil {
		log.Fatalf("error shuting server down: %v", err)
	}
}
