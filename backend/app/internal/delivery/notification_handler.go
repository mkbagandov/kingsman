package delivery

import (
	"encoding/json"
	"net/http"
	// "github.com/go-chi/chi/v5"
	"github.com/mkbagandov/kingsman/backend/app/internal/domain" // Import the domain package to access UserContextKey
	"github.com/mkbagandov/kingsman/backend/app/internal/usecase"
)

type NotificationHandler struct {
	notificationUseCase *usecase.NotificationUseCase
}

func NewNotificationHandler(notificationUseCase *usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{notificationUseCase: notificationUseCase}
}

func (h *NotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {

	var reqBody usecase.SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new SendNotificationRequest with UserID from context
	req := &usecase.SendNotificationRequest{
		UserID:  reqBody.UserID,
		Type:    reqBody.Type,
		Title:   reqBody.Title,
		Message: reqBody.Message,
	}

	if err := h.notificationUseCase.SendNotification(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	ctxUserID := r.Context().Value(domain.UserContextKey)
	if ctxUserID == nil {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	userID := ctxUserID.(string)

	resp, err := h.notificationUseCase.GetNotifications(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
