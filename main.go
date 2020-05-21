package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	Addr     string `json:"addr"`
	MongoURL string `json:"mongo_url"`
}

func main() {
	c := new(config)
	flag.StringVar(&c.Addr, "addr", envString("ADDR", ":80"), "address for listening")
	flag.StringVar(&c.MongoURL, "mongo_url", envString("MONGO_URL", "default mongo URL"), "mongoDB connection string")
	flag.Parse()

	log.Printf("config %+v", c)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		_ = json.NewEncoder(w).Encode(c)
	})

	log.Fatal(http.ListenAndServe(c.Addr, nil))
}

func envString(name, value string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}

	return value
}
