package vehicle

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/blazeisclone/vehicle-dms-inventory/pkg/strutils"
)

type VehicleService interface {
	Create(cmd CreateVehicleCommand) (*Vehicle, error)
	GetAll() ([]Vehicle, error)
	FindByID(id int) (*Vehicle, error)
	Update(id int, cmd UpdateVehicleCommand) (*Vehicle, error)
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

	mux.HandleFunc("GET "+strutils.APIPath("v1", "/vehicles"), h.Index)
	mux.HandleFunc("POST "+strutils.APIPath("v1", "/vehicles"), h.Store)
	mux.HandleFunc("GET "+strutils.APIPath("v1", "/vehicles/{id}"), h.Show)
	mux.HandleFunc("PUT "+strutils.APIPath("v1", "/vehicles/{id}"), h.Update)
	mux.HandleFunc("DELETE "+strutils.APIPath("v1", "/vehicles/{id}"), h.Destroy)
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

func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	var req vehicleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if errs := req.validate(); errs.HasErrors() {
		jsonValidationErrors(w, errs)
		return
	}

	vehicle, err := h.svc.Create(CreateVehicleCommand{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
	})

	if err != nil {
		jsonError(w, "failed to create vehicle", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, vehicle, http.StatusCreated)
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	vehicle, err := h.svc.FindByID(id)

	if errors.Is(err, ErrNotFound) {
		jsonError(w, "vehicle not found", http.StatusNotFound)
		return
	}

	if err != nil {
		jsonError(w, "failed to fetch vehicle", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, vehicle, http.StatusOK)
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

	if errs := req.validate(); errs.HasErrors() {
		jsonValidationErrors(w, errs)
		return
	}

	vehicle, err := h.svc.Update(id, UpdateVehicleCommand{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
	})

	if errors.Is(err, ErrNotFound) {
		jsonError(w, "vehicle not found", http.StatusNotFound)
		return
	}

	if err != nil {
		jsonError(w, "failed to update vehicle", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, vehicle, http.StatusOK)
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
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, status int) {
	jsonResponse(w, map[string]string{"error": msg}, status)
}

func jsonValidationErrors(w http.ResponseWriter, errs ValidationErrors) {
	jsonResponse(w, map[string]ValidationErrors{"errors": errs}, http.StatusUnprocessableEntity)
}
