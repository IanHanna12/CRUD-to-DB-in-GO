package main

import (
	handlers "github.com/IanHanna/CRUD-to-DB-in-GO/internal"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/db"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func main() {
	database := db.InitDB()
	handlers.InitHandlers(database)

	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("./frontend/static/"))

	router.POST("/login", handlers.LoginHandler)
	router.GET("/items/all", handlers.AuthMiddleware(false)(handlers.GetAllItemsHandler))
	router.POST("/items/create", handlers.AuthMiddleware(false)(handlers.CreateItemHandler))
	router.GET("/items/single/:id", handlers.AuthMiddleware(false)(handlers.GetItemByIDHandler))
	router.PUT("/items/update/:id", handlers.AuthMiddleware(false)(handlers.UpdateItemHandler))
	router.DELETE("/items/delete/:id", handlers.AuthMiddleware(false)(handlers.DeleteItemByIDHandler))
	router.DELETE("/items/all", handlers.AuthMiddleware(true)(handlers.DeleteAllItemsHandler))
	router.GET("/items/prefetch", handlers.AuthMiddleware(false)(handlers.PrefetchItemsHandler))
	router.GET("/validate-session", handlers.ValidateSessionHandler)

	serverConf := &http.Server{
		Handler:      GlobalCORS(router),
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

func GlobalCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:63342")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
