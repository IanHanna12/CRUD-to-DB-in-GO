package main

import (
	handlers "github.com/IanHanna/CRUD-to-DB-in-GO/internal"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/db"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"time"
)

func main() {
	db := db.InitDB()
	handlers.InitHandlers(db)

	router := mux.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("./frontend/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.PathPrefix("/static/").HandlerFunc(handlers.ServeWithProperMIME)

	// Main page handler
	router.HandleFunc("/", handlers.MainPageHandler)

	// API routes
	router.HandleFunc("/items", handlers.CreateItemHandler).Methods("POST")
	router.HandleFunc("/items", handlers.GetAllItemsHandler).Methods("GET")
	router.HandleFunc("/items/{id}", handlers.GetItemByIDHandler).Methods("GET")
	router.HandleFunc("/items/{id}", handlers.UpdateItemHandler).Methods("PUT")
	router.HandleFunc("/items/{id}", handlers.DeleteItemByIDHandler).Methods("DELETE")
	router.HandleFunc("/items", handlers.DeleteAllItemsHandler).Methods("DELETE")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:63832"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := corsHandler.Handler(router)

	serverConf := &http.Server{
		Handler:      handler,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Starting server on :8080")
	if err := serverConf.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
