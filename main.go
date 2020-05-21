package main

import (
	"context"
	"flag"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"log"
	"os"
	"os/signal"
	"time"
)

type config struct {
	Addr     string `json:"addr"`
	MongoURL string `json:"mongo_url"`
}

func main() {
	c := new(config)
	flag.StringVar(&c.Addr, "addr", envString("ADDR", ":80"), "address for listening")
	flag.StringVar(&c.MongoURL, "mongo_url", envString("MONGO_URL", ""), "mongoDB connection string")
	flag.Parse()

	ctx := context.Background()
	client, err := connect(ctx, c.MongoURL)
	if err != nil {
		log.Fatalf("connect to MongoDB failed: URI = %s, error = %v", c.MongoURL, err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("ping MongoDB failed: URI = %s, error = %v", c.MongoURL, err)
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Renderer = &renderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	db := client.Database("vquiz")

	// add sample question
	for _, q := range []Question{
		{
			Content: " Explain what is GO?",
			Answer:  `GO is an open source programming language which makes it easy to build simple, reliable and efficient software. Programs are constructed from packages, whose properties allow efficient management of dependencies.`,
		},
		{
			Content: "Explain what is string types?",
			Answer:  "A string type represents the set of string values, and string values are sequence of bytes.  Strings once created is not possible to change.",
		},
		{
			Content: "Explain how arrays in GO works differently then C?",
			Answer: `In GO Array works differently than it works in C

- Arrays are values, assigning one array to another copies all the elements
- If you pass an array to a function, it will receive a copy of the array, not a pointer to it
- The size of an array is part of its type. The types [10] int and [20] int are distinct`,
		}} {
		err := insertQuestion(context.TODO(), db, q)
		if err != nil {
			e.Logger.Errorf("insert a question failed: %v", err)
		}
	}

	s := server{db: db, e: e}
	s.routes()

	// Start server
	go func() {
		if err := e.Start(c.Addr); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// drop local database
	if err := db.Drop(context.Background()); err != nil {
		e.Logger.Errorf("drop database failed %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func envString(name, value string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}

	return value
}

type renderer struct {
	templates *template.Template
}

func (t *renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
