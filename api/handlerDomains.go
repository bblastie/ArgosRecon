package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type DomainResponse struct {
	ID        uuid.UUID `json:"id"`
	NAME      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) getDomains(w http.ResponseWriter, r *http.Request) {
	domains, err := cfg.DB.AllDomains(r.Context())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "error retrieving domains from database", err)
		return
	}

	var domainsJson []DomainResponse

	for _, domain := range domains {
		domainsJson = append(domainsJson, DomainResponse{
			ID:        domain.ID,
			NAME:      domain.Name,
			CreatedAt: domain.CreatedAt,
			UpdatedAt: domain.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, domainsJson)
}

func (cfg *apiConfig) deleteDomain(w http.ResponseWriter, r *http.Request) {
	domainId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid ID", err)
		return
	}
	dbDomain, err := cfg.DB.LookupDomainByID(r.Context(), domainId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Domain with ID %s not found", domainId), err)
		return
	}

	log.Printf("Preparing to delete %s", dbDomain.ID)
	log.Printf("If you want this delete to be permanent, make sure to delete the domain from domains.txt too")

	delErr := cfg.DB.DeleteOneDomain(r.Context(), dbDomain.ID)
	if delErr != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting %s from db", dbDomain.Name), delErr)
		return
	}

	respondWithJSON(w, http.StatusNoContent, DomainResponse{})
}
