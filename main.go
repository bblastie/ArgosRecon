package main

import (
	"fmt"
)

// func subfinder_run() error {
// 	db, err := sql.Open("postgres", fmt.Sprintf(
// 		"host=%s user=%s password=%s dbname=%s sslmode=disable",
// 		os.Getenv("db"),
// 		os.Getenv("POSTGRES_USER"),
// 		os.Getenv("POSTGRES_PASS"),
// 		os.Getenv("POSTGRES_DB"),
// 	))

// 	if err != nil {
// 		return fmt.Errorf("subfinder did not connect to db: %s", err)
// 	}

// 	defer db.Close()

// 	domain := os.Getenv("TARGET_DOMAIN")

// 	cmd := exec.Command("subfinder", "-d", domain)
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return fmt.Errorf("subfinder error: %s", err)
// 	}

// 	subdomains := strings.Split(string(output), "\n")

// 	for _, subdomain := range subdomains {
// 		if subdomain != "" {
// 			fmt.Print(subdomain)
// 		}
// 	}
// 	return nil
// }

func main() {
	fmt.Printf("Hello, world\n")
	// err := subfinder_run()
	// if err != nil {
	// 	panic(err)
	// }
}
