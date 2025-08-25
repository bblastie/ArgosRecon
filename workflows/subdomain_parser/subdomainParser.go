package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type SubdomainJson struct {
	Apex       string   `json:"apex"`
	Subdomains []string `json:"subdomains"`
}

func apexExtract(apexDomains, subdomains []string) map[string][]string {
	result := make(map[string][]string)

	for _, apex := range apexDomains {
		result[apex] = []string{}
	}

	for _, subdomain := range subdomains {
		cleanSubdomain := strings.TrimPrefix(subdomain, "*.")
		cleanSubdomain = strings.TrimPrefix(cleanSubdomain, ".")

		for _, apex := range apexDomains {
			if strings.HasSuffix(cleanSubdomain, "."+apex) {
				result[apex] = append(result[apex], cleanSubdomain)
				log.Printf("[+] apex: %s [+] sub: %s", apex, subdomain)
				break
			}
		}
	}

	return result
}

func main() {
	chaosFile, err := os.Open("/argos_data/chaos.txt")
	if err != nil {
		log.Printf("Cannot open chaos file: %s", err)
		return
	}

	subfinderFile, err := os.Open("/argos_data/subfinder.txt")
	if err != nil {
		log.Printf("Cannot open subfinder file: %s", err)
		return
	}

	domainsFile, err := os.Open("/argos_data/domains.txt")
	if err != nil {
		log.Printf("Cannot open domains file: %s", err)
		return
	}

	var apexDomains []string
	var subdomains []string

	domainsScanner := bufio.NewScanner(domainsFile)
	subfinderScanner := bufio.NewScanner(subfinderFile)
	chaosScanner := bufio.NewScanner(chaosFile)

	for domainsScanner.Scan() {
		apexDomains = append(apexDomains, domainsScanner.Text())
		fmt.Printf("%s\n", domainsScanner.Text())
	}

	for subfinderScanner.Scan() {
		subdomains = append(subdomains, subfinderScanner.Text())
	}

	for chaosScanner.Scan() {
		subdomains = append(subdomains, chaosScanner.Text())
	}

	organizedList := apexExtract(apexDomains, subdomains)

	subJsonPayload := []SubdomainJson{}

	for apex, subdomains := range organizedList {
		subJsonPayload = append(subJsonPayload, SubdomainJson{
			Apex:       apex,
			Subdomains: subdomains,
		})
	}

	jsonData, err := json.Marshal(subJsonPayload)
	if err != nil {
		log.Printf("Failed to marshall json: %s", err)
		return
	}

	req, err := http.NewRequest("PUT", "http://argos_main:8080/domains/subdomains", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to send request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Client failed to send request: %s", err)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Request body: %s", string(body))
		log.Fatalf("Failed to read response body: %s", err)
	}

	log.Printf("Request result: %s", res.Status)
}
