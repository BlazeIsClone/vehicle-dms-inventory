package vehicle

import (
	"database/sql"
	"net/http"
	"os"
)

func Routes(router *http.ServeMux, db *sql.DB) {
	handler := NewHandler(db)

	basePath := os.Getenv("BASE_PATH")

	router.HandleFunc("GET "+basePath+"/vehicles", handler.Index)
	router.HandleFunc("POST "+basePath+"/vehicles", handler.Store)
	router.HandleFunc("GET "+basePath+"/vehicles/{id}", handler.Show)
	router.HandleFunc("PUT "+basePath+"/vehicles/{id}", handler.Update)
	router.HandleFunc("DELETE "+basePath+"/vehicles/{id}", handler.Destroy)
}
