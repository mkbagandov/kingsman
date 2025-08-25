package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mkbagandov/kingsman/backend/app/internal/usecase"
)

type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
}

func NewProductHandler(productUseCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUseCase: productUseCase}
}

func (h *ProductHandler) GetProductCatalog(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category_id")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var minPrice *float64
	if minPriceStr != "" {
		p, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			http.Error(w, "Invalid min_price", http.StatusBadRequest)
			return
		}
		minPrice = &p
	}

	var maxPrice *float64
	if maxPriceStr != "" {
		p, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			http.Error(w, "Invalid max_price", http.StatusBadRequest)
			return
		}
		maxPrice = &p
	}

	limit := 100 // Default limit
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l <= 0 {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		limit = l
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err != nil || o < 0 {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		offset = o
	}

	req := &usecase.GetProductCatalogRequest{
		CategoryID: &categoryID,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		SortBy:     &sortBy,
		SortOrder:  &sortOrder,
		Limit:      limit,
		Offset:     offset,
	}

	resp, err := h.productUseCase.GetProductCatalog(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	resp, err := h.productUseCase.GetProductByID(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
