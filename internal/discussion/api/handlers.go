package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"RESTAPI/internal/discussion/model"
	"RESTAPI/internal/discussion/service"
)

// Handler handles HTTP requests for messages
type Handler struct {
	service *service.MessageService
}

// NewHandler creates a new Handler
func NewHandler(service *service.MessageService) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the API routes
func (h *Handler) RegisterRoutes(r *mux.Router) {
	api := r.PathPrefix("/api/v1.0").Subrouter()
	api.HandleFunc("/messages", h.GetAllMessages).Methods(http.MethodGet)
	api.HandleFunc("/messages", h.CreateMessage).Methods(http.MethodPost)
	api.HandleFunc("/messages/{id:[0-9]+}", h.GetMessage).Methods(http.MethodGet)
	api.HandleFunc("/messages/news/{newsId:[0-9]+}", h.GetMessagesByNewsID).Methods(http.MethodGet)
	api.HandleFunc("/messages/{id:[0-9]+}", h.UpdateMessage).Methods(http.MethodPut)
	api.HandleFunc("/messages/{id:[0-9]+}", h.DeleteMessage).Methods(http.MethodDelete)
}

// GetAllMessages handles retrieving all messages
func (h *Handler) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := h.service.GetAllMessages(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if messages == nil {
		messages = []*model.Message{} // Return empty array instead of null
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// CreateMessage handles message creation
func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message model.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateMessage(r.Context(), &message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// GetMessage handles retrieving a single message
func (h *Handler) GetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	message, err := h.service.GetMessage(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if message == nil {
		http.Error(w, "message not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// GetMessagesByNewsID handles retrieving messages for a news item
func (h *Handler) GetMessagesByNewsID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newsID, _ := strconv.ParseInt(vars["newsId"], 10, 64)

	messages, err := h.service.GetMessagesByNewsID(r.Context(), newsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// UpdateMessage handles message updates
func (h *Handler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	var message model.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	message.ID = id

	if err := h.service.UpdateMessage(r.Context(), &message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(message)
}

// DeleteMessage handles message deletion
func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	if err := h.service.DeleteMessage(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
} 