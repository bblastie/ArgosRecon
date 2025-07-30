package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) addDomains(w http.ResponseWriter, r *http.Request) {
	type DomainResponse struct {
		ID        uuid.UUID `json:"id"`
		NAME      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	file, err := os.Open("/argos_data/domains.txt")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to open domains.txt file", err)
		return
	}

	defer file.Close()

	var fileDomains []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fileDomains = append(fileDomains, scanner.Text())
	}

	if len(fileDomains) == 0 {
		respondWithError(w, http.StatusBadRequest, "Domains.txt is empty, please add domains", nil)
		return
	}

	var insertErrors []error
	for _, domain := range fileDomains {
		_, err := cfg.DB.InsertDomain(r.Context(), domain)
		if err != nil {
			log.Printf("Error inserting domain %s: %v\n", domain, err)
			insertErrors = append(insertErrors, err)
		}
	}
	if len(insertErrors) > 0 {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%d domains could not be added to database", len(insertErrors)), nil)
		return
	}

	log.Printf("Successfully inserted %d domains into databse", len(fileDomains))

	domains, err := cfg.DB.AllDomains(r.Context())

	var response []DomainResponse
	for _, domain := range domains {
		response = append(response, DomainResponse{
			ID:        domain.ID,
			NAME:      domain.Name,
			CreatedAt: domain.CreatedAt,
			UpdatedAt: domain.UpdatedAt,
		})
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No domains found in DB", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}
