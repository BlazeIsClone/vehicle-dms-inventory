package vehicle

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Handler struct {
	repo VehicleRepository
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		repo: NewMySQLVehicleRepository(db),
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	vehicle, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error finding vehicle: %v", err), http.StatusInternalServerError)
		return
	}

	if vehicle == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicle)
}

func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	vehicle := &Vehicle{
		Name:        payload.Name,
		Description: payload.Description,
	}

	if err := h.repo.Create(vehicle); err != nil {
		http.Error(w, fmt.Sprintf("Error creating vehicle: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vehicle)
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	vehicle, err := h.repo.FindByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error finding vehicle: %v", err), http.StatusInternalServerError)
		return
	}

	if vehicle == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicle)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	vehicle := &Vehicle{
		ID:          payload.ID,
		Name:        payload.Name,
		Description: payload.Description,
	}

	err = h.repo.UpdateByID(id, vehicle)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating vehicle: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Vehicle updated successfully")
}

func (h *Handler) Destroy(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	err = h.repo.DeleteByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting vehicle: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Vehicle deleted successfully")
}
