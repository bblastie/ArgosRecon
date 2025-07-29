package main

import (
	"argos/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL not set in .env")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Printf("Database connection failed: %s", err)
		return
	}

	dbQueries := database.New(db)

	const port = "8080"

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /domains", apiCfg.addDomains)
	mux.HandleFunc("PUT /domains/subdomains", apiCfg.addSubdomains)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
