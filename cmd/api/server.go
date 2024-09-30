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
	"github.com/lmittmann/tint"
)

type application struct {
	logger *slog.Logger
	db     *pgxpool.Pool
}

func main() {
	port := flag.Int("port", 8080, "Port number to serve the server")
	flag.Parse()

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	}))

	const (
		host     = "localhost"
		dbPort   = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "postgres"
		sslmode  = "disable"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", host, dbPort, user, dbname, password, sslmode)
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
