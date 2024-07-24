package main

import (
	handlers "github.com/IanHanna/CRUD-to-DB-in-GO/internal"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/db"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	db := db.InitDB()
	handlers.InitHandlers(db)

	router := mux.NewRouter()

	router.HandleFunc("/items", handlers.CreateItemHandler).Methods("POST")
	router.HandleFunc("/items", handlers.GetAllItemsHandler).Methods("GET")
	router.HandleFunc("/items/{id}", handlers.GetItemByIDHandler).Methods("GET")
	router.HandleFunc("/items/{id}", handlers.UpdateItemHandler).Methods("PUT")
	router.HandleFunc("/items/{id}", handlers.DeleteItemByIDHandler).Methods("DELETE")
	router.HandleFunc("/items", handlers.DeleteAllItemsHandler).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
