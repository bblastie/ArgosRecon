package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	fmt.Printf("Hello, world\n")
	cmd := exec.Command(
		"docker", "exec", "subfinder", "subfinder", "-d", "floqast.app", "-o", "/recondata/floqast.txt",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
