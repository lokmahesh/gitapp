package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	http.HandleFunc("/webhook", handleWebhook)

	log.Println("GitHub App listening on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
