package main

import (
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
