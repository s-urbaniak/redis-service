package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"gopkg.in/redis.v3"
)

var client *redis.Client

func errHandler(w http.ResponseWriter, err error) {
	log.Printf("error occured:\n%+v\n", err)

	http.Error(
		w,
		fmt.Sprintf("%s", err),
		http.StatusInternalServerError,
	)
}

func inc(w http.ResponseWriter, r *http.Request) {
	if err := client.Incr("counter").Err(); err != nil {
		errHandler(w, errors.Wrap(err, "counter increment failed"))
		return
	}

	n, err := client.Get("counter").Int64()
	if err != nil {
		errHandler(w, errors.Wrap(err, "reading counter failed"))
		return
	}

	fmt.Fprintf(w, "count: %d\n", n)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "healthy")
}

func main() {
	addr := flag.String("addr", ":8080", "the address to listen on")
	redisAddr := flag.String("redis-addr", "localhost:6379", "address of the redis server")
	flag.Parse()

	client = redis.NewClient(&redis.Options{
		Addr:     *redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/", inc)

	log.Println("Starting redis service...")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
