package main

import (
	"log"
	"net/http"

	"github.com/IanHanna/CRUD-to-DB-in-GO/internal"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/db"
	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()

	router := mux.NewRouter()

	router.HandleFunc("/items", handlers.CreateItemHandler).Methods("POST")
	router.HandleFunc("/items", handlers.GetAllItemsHandler).Methods("GET")
	router.HandleFunc("/items/{id}", handlers.GetItemByIDHandler).Methods("GET")
	router.HandleFunc("/items/{id}", handlers.UpdateItemHandler).Methods("PUT")
	router.HandleFunc("/items/{id}", handlers.DeleteItemByIDHandler).Methods("DELETE")
	router.HandleFunc("/items", handlers.DeleteAllItemsHandler).Methods("DELETE")

	log.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
