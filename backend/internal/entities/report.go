package entities

import "time"

type Report struct {
	ID         string     `json:"id"`
	ReporterID *int       `json:"reporter_id,omitempty"`
	Kind       string     `json:"kind"` // 'thread' or 'user'
	TargetID   int        `json:"target_id"`
	Reason     string     `json:"reason,omitempty"`
	Status     string     `json:"status"` // 'open','resolved','dismissed'
	CreatedAt  time.Time  `json:"created_at"`
	ResolvedBy *int       `json:"resolved_by,omitempty"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}
