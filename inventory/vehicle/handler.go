package vehicle

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type VehicleService interface {
	Create(name, description string) (*Vehicle, error)
	GetAll() ([]Vehicle, error)
	FindByID(id int) (*Vehicle, error)
	Update(id int, name, description string) (*Vehicle, error)
	Delete(id int) error
}

type Handler struct {
	svc VehicleService
}

func NewHandler(svc VehicleService) *Handler {
	return &Handler{svc: svc}
}

func Routes(mux *http.ServeMux, svc VehicleService) {
	h := NewHandler(svc)
	const prefix = "/api/v1"

	mux.HandleFunc("GET "+prefix+"/vehicles", h.Index)
	mux.HandleFunc("POST "+prefix+"/vehicles", h.Store)
	mux.HandleFunc("GET "+prefix+"/vehicles/{id}", h.Show)
	mux.HandleFunc("PUT "+prefix+"/vehicles/{id}", h.Update)
	mux.HandleFunc("DELETE "+prefix+"/vehicles/{id}", h.Destroy)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.svc.GetAll()
	if err != nil {
		jsonError(w, "failed to fetch vehicles", http.StatusInternalServerError)
		return
	}
	if vehicles == nil {
		vehicles = []Vehicle{}
	}
	jsonResponse(w, vehicles, http.StatusOK)
}

type vehicleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	var req vehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		jsonError(w, "name is required", http.StatusUnprocessableEntity)
		return
	}
	v, err := h.svc.Create(req.Name, req.Description)
	if err != nil {
		jsonError(w, "failed to create vehicle", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, v, http.StatusCreated)
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	v, err := h.svc.FindByID(id)
	if errors.Is(err, ErrNotFound) {
		jsonError(w, "vehicle not found", http.StatusNotFound)
		return
	}
	if err != nil {
		jsonError(w, "failed to fetch vehicle", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, v, http.StatusOK)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var req vehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		jsonError(w, "name is required", http.StatusUnprocessableEntity)
		return
	}
	v, err := h.svc.Update(id, req.Name, req.Description)
	if errors.Is(err, ErrNotFound) {
		jsonError(w, "vehicle not found", http.StatusNotFound)
		return
	}
	if err != nil {
		jsonError(w, "failed to update vehicle", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, v, http.StatusOK)
}

func (h *Handler) Destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	err := h.svc.Delete(id)
	if errors.Is(err, ErrNotFound) {
		jsonError(w, "vehicle not found", http.StatusNotFound)
		return
	}
	if err != nil {
		jsonError(w, "failed to delete vehicle", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func jsonResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data) //nolint:errcheck
}

func jsonError(w http.ResponseWriter, msg string, status int) {
	jsonResponse(w, map[string]string{"error": msg}, status)
}
