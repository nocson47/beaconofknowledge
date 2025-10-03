package entities

import "time"

type Reply struct {
	ID       int `json:"id"`
	ThreadID int `json:"thread_id"`
	UserID   int `json:"user_id"`
	// Author username (denormalized for read responses)
	Author    string     `json:"author,omitempty"`
	ParentID  *int       `json:"parent_id,omitempty"`
	Body      string     `json:"body"`
	IsDeleted bool       `json:"is_deleted"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
