package main

import (
	"fmt"
	"log"
	"net/http"

	libredis "github.com/redis/go-redis/v9"

	limiter "turbo-dollop"

	mhttp "turbo-dollop/drivers/middleware/stdlib"

	sredis "turbo-dollop/drivers/store/redis"
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

	// Create a store with the redis client.
	store, err := sredis.NewStorageWithOptions(client, limiter.StorageOptions{
		Prefix:   "limiter_http_example",
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a new middleware with the limiter instance.
	middleware := mhttp.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))

	// Launch a simple server.
	http.Handle("/", middleware.Handler(http.HandlerFunc(index)))
	fmt.Println("Server is running on port 7777...")
	log.Fatal(http.ListenAndServe(":7777", nil))

}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write([]byte(`{"message": "ok"}`))
	if err != nil {
		log.Fatal(err)
	}
}
