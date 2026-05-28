package repositories

import (
	"context"
	"database/sql"
	"time"

	"pantrypal/backend/internal/platform/id"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

type StoredChatMessage struct {
	ID        string
	UserID    string
	Role      string
	Action    string
	Content   string
	CreatedAt string
}

func (r *ChatRepository) InsertBotMessage(ctx context.Context, userID, content string) (StoredChatMessage, error) {
	msgID, err := id.New("msg")
	if err != nil {
		return StoredChatMessage{}, err
	}

	now := time.Now().UTC().Format(time.RFC3339)

	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO chat_messages (id, user_id, role, action, content, metadata_json, created_at)
		 VALUES (?, ?, 'assistant', '', ?, '{}', ?)`,
		msgID,
		userID,
		content,
		now,
	)
	if err != nil {
		return StoredChatMessage{}, err
	}

	return StoredChatMessage{
		ID:        msgID,
		UserID:    userID,
		Role:      "assistant",
		Content:   content,
		CreatedAt: now,
	}, nil
}

func (r *ChatRepository) InsertUserMessage(ctx context.Context, userID, message, action string) (StoredChatMessage, error) {
	msgID, err := id.New("msg")
	if err != nil {
		return StoredChatMessage{}, err
	}

	now := time.Now().UTC().Format(time.RFC3339)

	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO chat_messages (id, user_id, role, action, content, metadata_json, created_at)
		 VALUES (?, ?, 'user', NULLIF(?, ''), ?, '{}', ?)`,
		msgID,
		userID,
		action,
		message,
		now,
	)
	if err != nil {
		return StoredChatMessage{}, err
	}

	return StoredChatMessage{
		ID:        msgID,
		UserID:    userID,
		Role:      "user",
		Action:    action,
		Content:   message,
		CreatedAt: now,
	}, nil
}

func (r *ChatRepository) ListRecent(ctx context.Context, userID string, limit int) ([]StoredChatMessage, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, role, COALESCE(action, ''), content, created_at
		 FROM chat_messages
		 WHERE user_id = ?
		 ORDER BY created_at DESC
		 LIMIT ?`,
		userID,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]StoredChatMessage, 0, limit)
	for rows.Next() {
		var msg StoredChatMessage
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Role, &msg.Action, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
