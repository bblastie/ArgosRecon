package main

import (
	"argos/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Printf("Database connection failed: %s", err)
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

	mux.HandleFunc("POST /api/domains", apiCfg.addDomains)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
