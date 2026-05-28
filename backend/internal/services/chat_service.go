package services

import (
	"context"

	"pantrypal/backend/internal/repositories"
	"pantrypal/backend/internal/transport/http/dto"
)

type ChatService struct {
	chat *repositories.ChatRepository
}

func NewChatService(chat *repositories.ChatRepository) *ChatService {
	return &ChatService{chat: chat}
}

func (s *ChatService) SendMessage(ctx context.Context, userID string, req dto.ChatSendRequest) (dto.ChatMessageResponse, error) {
	msg, err := s.chat.InsertUserMessage(ctx, userID, req.Message, req.Action)
	if err != nil {
		return dto.ChatMessageResponse{}, err
	}
	return dto.ChatMessageResponse{
		ID:        msg.ID,
		Role:      msg.Role,
		Action:    msg.Action,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
	}, nil
}

func (s *ChatService) GetHistory(ctx context.Context, userID string, limit int) (dto.ChatHistoryResponse, error) {
	messages, err := s.chat.ListRecent(ctx, userID, limit)
	if err != nil {
		return dto.ChatHistoryResponse{}, err
	}
	out := make([]dto.ChatMessageResponse, 0, len(messages))
	for _, msg := range messages {
		out = append(out, dto.ChatMessageResponse{
			ID:        msg.ID,
			Role:      msg.Role,
			Action:    msg.Action,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		})
	}
	return dto.ChatHistoryResponse{Messages: out}, nil
}
