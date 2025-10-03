package entities

import (
	"fmt"
	"strings"
	"time"
)

type Vote struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ThreadID  *int      `json:"thread_id,omitempty"`
	ReplyID   *int      `json:"reply_id,omitempty"`
	Value     int       `json:"value"` // 1 or -1
	CreatedAt time.Time `json:"created_at"`
}

// IsUp returns true if the vote represents an upvote.
func (v *Vote) IsUp() bool { return v.Value == 1 }

// IsDown returns true if the vote represents a downvote.
func (v *Vote) IsDown() bool { return v.Value == -1 }

// ValueString returns a human-readable string for the vote value.
func (v *Vote) ValueString() string {
	if v.IsUp() {
		return "up"
	}
	if v.IsDown() {
		return "down"
	}
	return "unknown"
}

// ParseVoteValueString parses common representations into the integer value used by the DB.
// Accepts: "up", "down", "1", "-1" (case-insensitive)
func ParseVoteValueString(s string) (int, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "up", "1":
		return 1, nil
	case "down", "-1":
		return -1, nil
	default:
		return 0, fmt.Errorf("invalid vote value: %q", s)
	}
}
