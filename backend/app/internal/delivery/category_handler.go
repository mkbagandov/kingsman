package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/mkbagandov/kingsman/backend/app/internal/usecase"
)

type CategoryHandler struct {
	categoryUseCase *usecase.CategoryUseCase
}

func NewCategoryHandler(categoryUseCase *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryUseCase: categoryUseCase}
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	resp, err := h.categoryUseCase.GetCategories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
