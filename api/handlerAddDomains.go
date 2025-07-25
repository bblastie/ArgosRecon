package main

import (
	"bufio"
	"net/http"
	"os"
)

func (cfg *apiConfig) addDomains(w http.ResponseWriter, r *http.Request) {
	type response struct {
		domains []string
	}

	file, err := os.Open("/argos_data/domains.txt")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to open domains.txt file", err)
		return
	}

	defer file.Close()

	var domains []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}

	for _, domain := range domains {
		_, err := cfg.DB.InsertDomain(r.Context(), domain)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Domain could not be added to database", err)
			continue
		}
	}
	respondWithJSON(w, http.StatusOK, response{domains})
}
