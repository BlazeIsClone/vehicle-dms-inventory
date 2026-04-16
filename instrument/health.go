package instrument

import (
	"net/http"
	"os"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func hostNameHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	hostname, err := os.Hostname()
	if err != nil {
		http.Error(w, "Error retrieving hostname", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(hostname))
}

func Routes(router *http.ServeMux) {
	router.HandleFunc("GET /health", healthCheckHandler)
	router.HandleFunc("GET /hostname", hostNameHandler)
}
