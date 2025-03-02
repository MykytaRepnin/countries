package main

import (
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	// Initialize ClickHouse connection
	InitClickHouse()

	// Initialize router
	r := InitRouter()

	log.Println("Server running on http://localhost:8080")
	if err := fasthttp.ListenAndServe(":8080", r.Handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
