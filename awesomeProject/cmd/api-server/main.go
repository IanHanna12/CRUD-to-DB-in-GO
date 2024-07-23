package main

import (
	handlers "github.com/IanHanna/CRUD-to-DB-in-GO/internal"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/db"
	"log"
	"net/http"
)

func main() {
	db.InitDB()

	http.HandleFunc("/create-item", handlers.CreateItemHandler) // Register the correct endpoint
	http.HandleFunc("/items/all", handlers.GetAllItemsHandler)
	http.HandleFunc("/items/get", handlers.GetItemByIDHandler)
	http.HandleFunc("/items/update", handlers.UpdateItemHandler)
	http.HandleFunc("/items/delete", handlers.DeleteItemByIDHandler)
	http.HandleFunc("/items/deleteAll", handlers.DeleteAllItemsHandler) // Register the new handler

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
