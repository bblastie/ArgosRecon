package main

import (
	"argos/internal/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type subdomainReq struct {
	ApexDomain string   `json:"apex"`
	Subdomains []string `json:"subdomains"`
}

type subdomainResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	DomainID   uuid.UUID `json:"domain_id"`
	DomainName string    `json:"domain_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (cfg *apiConfig) addSubdomains(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var jsonData []subdomainReq
	var jsonResponse []subdomainResponse

	err := decoder.Decode(&jsonData)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid json", err)
		return
	}

	for _, obj := range jsonData {
		dbDomain, err := cfg.DB.LookupDomainByName(r.Context(), obj.ApexDomain)
		if err != nil {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("apex domain %s not in database. Please add it to domains.txt and POST /domains", obj.ApexDomain), err)
			return
		}

		var subdomainInsertErrors []error
		var insertedSubdomains []database.Subdomain

		for _, subdomain := range obj.Subdomains {
			subdomainToInsert := database.InsertSubdomainParams{
				Name:     subdomain,
				DomainID: dbDomain.ID,
			}

			dbSubdomain, err := cfg.DB.InsertSubdomain(r.Context(), subdomainToInsert)
			if err != nil {
				log.Printf("Error adding %s to database", dbSubdomain.Name)
				subdomainInsertErrors = append(subdomainInsertErrors, err)
			}
			insertedSubdomains = append(insertedSubdomains, dbSubdomain)
		}
		if len(subdomainInsertErrors) > 0 {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%d subdomains could not be added to database", len(subdomainInsertErrors)), nil)
			return
		}

		for _, sub := range insertedSubdomains {
			dbSub, err := cfg.DB.LookupSubdomainByName(r.Context(), sub.Name)
			if err != nil {
				log.Printf("error looking up %s from database", dbSub.Name)
			}
			jsonSub := subdomainResponse{
				ID:         dbSub.ID.String(),
				Name:       dbSub.Name,
				DomainID:   dbSub.DomainID,
				DomainName: dbSub.DomainName,
				CreatedAt:  dbSub.CreatedAt,
				UpdatedAt:  dbSub.UpdatedAt,
			}
			jsonResponse = append(jsonResponse, jsonSub)
		}
	}
	respondWithJSON(w, http.StatusCreated, jsonResponse)
}

func (cfg *apiConfig) deleteSubdomain(w http.ResponseWriter, r *http.Request) {
	subdomainID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid ID", err)
		return
	}

	dbSubdomain, err := cfg.DB.LookupSubdomainByID(r.Context(), subdomainID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Subdomain with ID %s not found", subdomainID), err)
		return
	}

	log.Printf("Preparing to delete %s from database", dbSubdomain.Name)

	delErr := cfg.DB.DeleteSubdomainByID(r.Context(), dbSubdomain.ID)
	if delErr != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error deleting %s from databse", dbSubdomain.Name), err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, subdomainResponse{})
}
