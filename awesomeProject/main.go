package main

import (
	"fmt"
	"log"
	"my-awesome-project/app"
	"net/http"
	"os"
)

func main() {
	app.InitDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	http.HandleFunc("/create-item", app.HandleCreateItem)
	http.HandleFunc("/customers", app.HandleGetAllItems)
	http.HandleFunc("/customer", app.HandleGetItemByID)
	http.HandleFunc("/update", app.HandleUpdateItem)
	http.HandleFunc("/delete", app.HandleDeleteItemByID)
	http.HandleFunc("/deleteAll", app.HandleDeleteAllItems)

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
