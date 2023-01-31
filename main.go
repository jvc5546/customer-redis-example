package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"os"
	"context"

	"github.com/redis/go-redis/v9"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPW := os.Getenv("REDIS_PASSWORD")

	redisAddress := fmt.Sprintf("%s:%s", redisHost, redisPort)
	log.Printf(redisAddress)
	log.Printf("Version %s", "1.0.0")

	client := redis.NewClient(&redis.Options{
		Addr:	  redisAddress,
		Password: redisPW, // no password set
		DB:		  0,  // use default DB
	})

	key := randSeq(8)
	log.Printf("key set as %s", key)

	pong, err := client.Ping(ctx).Result()
	log.Println(pong, err)

	handlers := CounterHandlers{
		client: client,
		key:    key,
	}

	http.HandleFunc("/", hello)
	http.HandleFunc("/increment", func(w http.ResponseWriter, r *http.Request) {
		handlers.Increment(ctx, w, r)
	})
	http.HandleFunc("/decrement", func(w http.ResponseWriter, r *http.Request) {
		handlers.Decrement(ctx, w, r)
	})
	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		handlers.Count(ctx, w, r)
	})

	log.Println("Starting http server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type CounterHandlers struct {
	client *redis.Client
	key    string
}

func (h CounterHandlers) Increment(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	val, err := h.client.Incr(ctx, h.key).Result()
	if err != nil {
		log.Printf("error incrementing %v", err)
		http.Error(w, "error incrementing", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Incremented count to %d", val)
}

func (h CounterHandlers) Decrement(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	val, err := h.client.Decr(ctx, h.key).Result()
	if err != nil {
		log.Printf("error deccrementing %v", err)
		http.Error(w, "error deccrementing", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Decremented count to %d", val)
}

func (h CounterHandlers) Count(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	val, err := h.client.Get(ctx, h.key).Result()
	if err != nil {
		log.Printf("error retreiving value %v", err)
		http.Error(w, "error retrieving value", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Current count is %s", val)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
	<!doctype html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>Welcome</title>
	</head>
	<body>
		<h1>Welcome to the Redis Example</h1>
		<p>You can increment the stored count at <a href="/increment">/increment</a></p>
		<p>You can decrement the stored count at <a href="/decrement">/decrement</a></p>
		<p>You can retrieve the stored count at <a href="/count">/count</a></p>
	</body>
	</html>
`)
}
