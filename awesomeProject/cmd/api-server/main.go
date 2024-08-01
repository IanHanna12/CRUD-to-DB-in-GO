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
	router.ServeFiles("/static/*filepath", http.Dir("./frontend/static"))

	router.POST("/login", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handlers.LoginHandler(w, r)
	})
	router.GET("/items/all", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handlers.AuthMiddleware(false)(handlers.GetAllItemsHandler)(w, r)
	})
	router.POST("/items/create", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handlers.AuthMiddleware(false)(handlers.CreateItemHandler)(w, r)
	})
	router.GET("/items/single/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handlers.AuthMiddleware(false)(func(w http.ResponseWriter, r *http.Request) {
			handlers.GetItemByIDHandler(w, r, ps)
		})(w, r)
	})
	router.PUT("/items/update/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handlers.AuthMiddleware(false)(func(w http.ResponseWriter, r *http.Request) {
			handlers.UpdateItemHandler(w, r, ps)
		})(w, r)
	})
	router.DELETE("/items/delete/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handlers.AuthMiddleware(true)(func(w http.ResponseWriter, r *http.Request) {
			handlers.DeleteItemByIDHandler(w, r, ps)
		})(w, r)
	})
	router.DELETE("/items/all", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handlers.AuthMiddleware(true)(handlers.DeleteAllItemsHandler)(w, r)
	})
	router.GET("/items/prefetch", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handlers.AuthMiddleware(false)(handlers.PrefetchItemsHandler)(w, r)
	})
	router.GET("/validate-session", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handlers.ValidateSessionHandler(w, r)
	})

	serverConf := &http.Server{
		Handler:      globalCORS(router),
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

func globalCORS(router *httprouter.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		router.ServeHTTP(w, r)
	})
}
