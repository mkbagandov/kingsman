package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mkbagandov/kingsman/backend/app/internal/domain"
	"github.com/mkbagandov/kingsman/backend/app/internal/usecase"
)

type CartHandler struct {
	cartUseCase *usecase.CartUseCase
}

func NewCartHandler(cartUseCase *usecase.CartUseCase) *CartHandler {
	return &CartHandler{cartUseCase: cartUseCase}
}

func (h *CartHandler) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserContextKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req usecase.AddItemToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.UserID = userID

	cartItem, err := h.cartUseCase.AddItemToCart(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartItem)

}

func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserContextKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req usecase.UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.UserID = userID

	err := h.cartUseCase.UpdateCartItem(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Cart item updated successfully"})
}

func (h *CartHandler) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserContextKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	productID := chi.URLParam(r, "productID")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	req := usecase.RemoveCartItemRequest{UserID: userID, ProductID: productID}

	err := h.cartUseCase.RemoveCartItem(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Cart item updated successfully"})
}

func (h *CartHandler) GetUserCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserContextKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cart, err := h.cartUseCase.GetUserCart(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserContextKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.cartUseCase.ClearCart(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Cart cleared successfully"})
}
