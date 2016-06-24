package main

import (
	"log"
	"net/http"
	"strconv"

	"gopkg.in/redis.v3"
)

var client *redis.Client

func internalError(e error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(e.Error()))
}

func inc(w http.ResponseWriter, r *http.Request) {
	if err := client.Incr("counter").Err(); err != nil {
		internalError(err, w)
		return
	}

	n, err := client.Get("counter").Int64()
	if err != nil {
		internalError(err, w)
		return
	}

	w.Write([]byte(strconv.FormatInt(n, 10)))
}

func healthz(w http.ResponseWriter, r *http.Request) {
	pong, err := client.Ping().Result()

	if err != nil {
		internalError(err, w)
		return
	}

	w.Write([]byte(pong))
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/", inc)

	log.Println("Starting redis service...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
