package main

import (
	handlers "github.com/IanHanna/CRUD-to-DB-in-GO/internal"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/db"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func globalCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db := db.InitDB()
	handlers.InitHandlers(db)

	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("./frontend/static"))

	router.POST("/items", wrapHandler(handlers.CreateItemHandler))
	router.GET("/items", wrapHandler(handlers.WithAuthentication(handlers.GetallitemsHandler)))
	router.GET("/items/:id", wrapHandler(handlers.WithAuthentication(handlers.GetitembyIDHandler)))
	router.PUT("/items/:id", wrapHandler(handlers.WithAuthentication(handlers.UpdateitemHandler)))
	router.DELETE("/items/:id", wrapHandler(handlers.WithAdminAuthentication(handlers.DeleteitembyIDHandler)))
	router.DELETE("/items", wrapHandler(handlers.WithAdminAuthentication(handlers.DeleteallitemsHandler)))
	router.POST("/login", wrapHandler(handlers.LoginHandler))

	go func() {
		for {
			time.Sleep(1 * time.Hour)
			handlers.CleanupSessions()
		}
	}()

	handler := globalCORS(router)

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

func wrapHandler(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h(w, r)
	}
}
