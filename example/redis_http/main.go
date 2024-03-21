package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	libredis "github.com/redis/go-redis/v9"

	limiter "turbo-dollop"

	mhttp "turbo-dollop/drivers/middleware/stdlib"

	sredis "turbo-dollop/drivers/storage/redis"
)

func main() {

	// Define a limit rate to 4 requests per hour.
	rate, err := limiter.NewRateCreator("4|H")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a redis client.
	option, err := libredis.ParseURL("redis://localhost:6379/0")
	if err != nil {
		log.Fatal(err)
		return
	}
	client := libredis.NewClient(option)

	// Create a storage with the redis client.
	storage, err := sredis.NewStorageWithOptions(client, limiter.StorageOptions{
		Prefix:   "limiter_http_example",
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer storage.Close(context.Background())

	// Create a new middleware with the limiter instance.
	middleware := mhttp.NewMiddleware(limiter.New(storage, rate, limiter.WithTrustForwardHeader(true)))

	http.Handle("/", middleware.Handler(http.HandlerFunc(index)))
	httpServer := http.Server{
		Addr: ":7777",
	}

	// Launch a simple server.
	go func() {
		// service connections
		fmt.Println("Server is running on port 7777...")
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe Error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")

}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write([]byte(`{"message": "ok"}`))
	if err != nil {
		log.Fatal(err)
	}
}
