package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelMember struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ChannelID uuid.UUID  `gorm:"type:uuid;not null" json:"channel_id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	JoinedAt  time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"joined_at"`
	LeftAt    *time.Time `json:"left_at,omitempty"`
}
