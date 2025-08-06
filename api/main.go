package main

import (
	"argos/internal/database"
	"bufio"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func loadDomains(cfg *apiConfig) error {
	file, err := os.Open("/argos_data/domains.txt")
	if err != nil {
		log.Printf("Failed to open domains.txt. Error: %s", err)
		return err
	}

	scanner := bufio.NewScanner(file)

	var domains []string

	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}

	for _, domain := range domains {
		dbDomain, err := cfg.DB.InsertDomain(context.Background(), domain)
		if err != nil {
			log.Printf("Error inserting %s into database: %s", domain, err)
			return err
		}
		log.Printf("Successfully added %s with ID of: %s", dbDomain.Name, dbDomain.ID)
	}
	return nil
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

	domainErr := loadDomains(&apiCfg)
	if domainErr != nil {
		log.Fatalf("Could not load domains from domains.txt: %s", domainErr)
		return
	}

	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("GET /domains", apiCfg.getDomains)
	// tod do - add a GET /domains/{id} that returns all gathered data for the domain
	mux.HandleFunc("DELETE /domains/{id}", apiCfg.deleteDomain)
	mux.HandleFunc("GET /domains/{id}/subdomains", apiCfg.getSubdomainsByDomainId)
	mux.HandleFunc("PUT /domains/subdomains", apiCfg.addSubdomains)
	mux.HandleFunc("DELETE /domains/subdomains/{id}", apiCfg.deleteSubdomain)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
