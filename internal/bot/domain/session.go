package domain

import "time"

type UserSession struct {
	UserID        int64                  `json:"user_id"`
	CurrentState  ConversationState      `json:"current_state"`
	ActiveCommand string                 `json:"pending_cmd,omitempty"`
	TempData      map[string]interface{} `json:"temp_data,omitempty"`
	UpdatedAt     time.Time              `json:"updated_at"`
}
