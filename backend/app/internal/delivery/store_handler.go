package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mkbagandov/kingsman/backend/app/internal/usecase"
)

type StoreHandler struct {
	storeUseCase *usecase.StoreUseCase
}

func NewStoreHandler(storeUseCase *usecase.StoreUseCase) *StoreHandler {
	return &StoreHandler{storeUseCase: storeUseCase}
}

func (h *StoreHandler) GetStores(w http.ResponseWriter, r *http.Request) {
	resp, err := h.storeUseCase.GetStores(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *StoreHandler) GetStoreByID(w http.ResponseWriter, r *http.Request) {
	storeID := chi.URLParam(r, "storeID")
	if storeID == "" {
		http.Error(w, "Store ID is required", http.StatusBadRequest)
		return
	}

	store, err := h.storeUseCase.GetStoreByID(r.Context(), storeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store)
}
