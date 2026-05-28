package handlers

import (
	"net/http"
	"strconv"

	"pantrypal/backend/internal/services"
	"pantrypal/backend/internal/transport/http/dto"
	"pantrypal/backend/internal/transport/http/middleware"
)

type ChatHandler struct {
	chat *services.ChatService
}

func NewChatHandler(chat *services.ChatService) *ChatHandler {
	return &ChatHandler{chat: chat}
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	var req dto.ChatSendRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if req.Message == "" {
		WriteAPIError(w, http.StatusBadRequest, "validation_error", "message is required")
		return
	}

	res, err := h.chat.SendMessage(r.Context(), userID, req)
	if err != nil {
		WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to send message")
		return
	}

	WriteJSON(w, http.StatusCreated, res)
}

func (h *ChatHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	res, err := h.chat.GetHistory(r.Context(), userID, limit)
	if err != nil {
		WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to load chat history")
		return
	}

	WriteJSON(w, http.StatusOK, res)
}
