package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/retriever/fileretriever"
)

type application struct {
	logger *slog.Logger
	db     *pgxpool.Pool
}

func main() {
	port := flag.Int("port", 8080, "Port number to serve the server")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		slog.Info(".env file not detected")
	}

	var (
		host     = os.Getenv("POSTGRES_HOST")
		dbPort   = 5432
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = os.Getenv("POSTGRES_DB")
		env      = os.Getenv("ENVIRONMENT")
	)

	var logLevel = slog.LevelInfo
	if env == "DEV" {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      logLevel,
		TimeFormat: time.Kitchen,
	}))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", host, dbPort, user, dbname, password)
	dbpool, err := pgxpool.New(context.Background(), psqlInfo)
	if err != nil {
		logger.Error("unable to connect to database", "error", err)
	}
	defer dbpool.Close()
	err = dbpool.Ping(context.Background())
	if err != nil {
		logger.Error("error pinging database", "error", err)
	}

	app := application{
		logger: logger,
		db:     dbpool,
	}

	err = ffclient.Init(ffclient.Config{
		PollingInterval: 3 * time.Second,
		Retriever: &fileretriever.Retriever{
			Path: "feature-flags.yaml",
		},
	})
	if err != nil {
		logger.Error("feature-flags file not detected")
	}
	defer ffclient.Close()

	router := http.NewServeMux()

	router.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Path[len("/static/"):]
		fullPath := filepath.Join(".", "static", filePath)
		http.ServeFile(w, r, fullPath)
	})

	router.HandleFunc("GET /healthcheck", app.healthcheckHandler)
	router.HandleFunc("GET /", app.homeHandler)
	router.HandleFunc("POST /submit", app.submitHandler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}

	app.logger.Info(fmt.Sprintf("Starting server on port %s", server.Addr))
	server.ListenAndServe()
}
