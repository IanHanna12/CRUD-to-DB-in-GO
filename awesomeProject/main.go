package main

import (
	"fmt"
	"github.com/IanHanna/CRUD-to-DB-in-GO/app"
	"log"
	"net/http"
	"os"
)

func main() {
	app.InitDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	http.HandleFunc("/create-item", app.HandleCreateItem)
	http.HandleFunc("/items", app.HandleGetAllItems)
	http.HandleFunc("/item", app.HandleGetItemByID)
	http.HandleFunc("/update-item", app.HandleUpdateItem)
	http.HandleFunc("/delete-item", app.HandleDeleteItemByID)
	http.HandleFunc("/delete-all-items", app.HandleDeleteAllItems)

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
