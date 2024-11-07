package main

import (
	"log"
	"net/http"
	"students/routes"
)

func main() {
    router := routes.SetupRoutes()
    
    log.Println("Server starting on port 8080...")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal(err)
    }
}