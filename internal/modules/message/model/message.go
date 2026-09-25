package model

import (
	"github.com/google/uuid"
	domain "github.com/PritamDasTharu/conduit/internal/domain/model"
)

type Message struct {
	domain.Base
	ChannelID uuid.UUID `gorm:"type:uuid;not null;index" json:"channel_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
}
