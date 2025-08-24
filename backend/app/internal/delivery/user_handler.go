package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	// "github.com/go-chi/chi/v5"
	"github.com/mkbagandov/kingsman/backend/app/internal/domain" // Import the domain package to access UserContextKey
	"github.com/mkbagandov/kingsman/backend/app/internal/usecase"
)

type UserHandler struct {
	userUseCase    *usecase.UserUseCase
	loyaltyUseCase *usecase.LoyaltyUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase, loyaltyUseCase *usecase.LoyaltyUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase, loyaltyUseCase: loyaltyUseCase}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req usecase.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.userUseCase.RegisterUser(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req usecase.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.userUseCase.LoginUser(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctxUserID := r.Context().Value(domain.UserContextKey)
	if ctxUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	userID := ctxUserID.(string)

	resp, err := h.userUseCase.GetUserProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetUserLoyaltyProfile handles the request to get a user's loyalty profile.
func (h *UserHandler) GetUserLoyaltyProfile(w http.ResponseWriter, r *http.Request) {
	ctxUserID := r.Context().Value(domain.UserContextKey)
	if ctxUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	userIDStr := ctxUserID.(string)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID format in JWT", http.StatusBadRequest)
		return
	}

	resp, err := h.loyaltyUseCase.GetUserLoyaltyProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type AddLoyaltyPointsRequest struct {
	Points int    `json:"points"`
	Type   string `json:"type"`
}

// AddLoyaltyPoints handles the request to add loyalty points to a user.
func (h *UserHandler) AddLoyaltyPoints(w http.ResponseWriter, r *http.Request) {
	ctxUserID := r.Context().Value(domain.UserContextKey)
	if ctxUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	userIDStr := ctxUserID.(string)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID format in JWT", http.StatusBadRequest)
		return
	}

	var req AddLoyaltyPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.loyaltyUseCase.AddLoyaltyPoints(r.Context(), userID, req.Points, req.Type); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type AddLoyaltyActivityRequest struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// AddLoyaltyActivity handles the request to record a loyalty activity for a user.
func (h *UserHandler) AddLoyaltyActivity(w http.ResponseWriter, r *http.Request) {
	ctxUserID := r.Context().Value(domain.UserContextKey)
	if ctxUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	userIDStr := ctxUserID.(string)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID format in JWT", http.StatusBadRequest)
		return
	}

	var req AddLoyaltyActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.loyaltyUseCase.AddLoyaltyActivity(r.Context(), userID, req.Type, req.Description); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetLoyaltyTiers handles the request to get all loyalty tiers.
func (h *UserHandler) GetLoyaltyTiers(w http.ResponseWriter, r *http.Request) {
	resp, err := h.loyaltyUseCase.GetLoyaltyTiers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetUserDiscountCard handles the request to get a user's discount card information.
func (h *UserHandler) GetUserDiscountCard(w http.ResponseWriter, r *http.Request) {
	ctxUserID := r.Context().Value(domain.UserContextKey)
	if ctxUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	userID := ctxUserID.(string)

	resp, err := h.userUseCase.GetUserDiscountCard(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type UpdateUserDiscountCardRequest struct {
	DiscountLevel       int     `json:"discount_level"`
	ProgressToNextLevel float64 `json:"progress_to_next_level"`
}

// UpdateUserDiscountCard handles the request to update a user's discount card information.
func (h *UserHandler) UpdateUserDiscountCard(w http.ResponseWriter, r *http.Request) {
	ctxUserID := r.Context().Value(domain.UserContextKey)
	if ctxUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	userID := ctxUserID.(string)

	var req UpdateUserDiscountCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.UpdateUserDiscountCard(r.Context(), userID, req.DiscountLevel, req.ProgressToNextLevel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetUserQRCode handles the request to get a user's QR code image.
func (h *UserHandler) GetUserQRCode(w http.ResponseWriter, r *http.Request) {
	ctxUserID := r.Context().Value(domain.UserContextKey)
	if ctxUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	userID := ctxUserID.(string)

	qrCodeImage, err := h.userUseCase.GenerateQRCodeImage(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(qrCodeImage)
}
