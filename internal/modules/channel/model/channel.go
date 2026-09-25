package model

import (
	"github.com/google/uuid"
	domain "github.com/PritamDasTharu/conduit/internal/domain/model"
)

type Channel struct {
	domain.Base
	WorkspaceID uuid.UUID `gorm:"type:uuid;not null" json:"workspace_id"`
	Name        string    `gorm:"type:text;not null" json:"name"`
	Slug        string    `gorm:"type:text;not null" json:"slug"`
	Visibility  string    `gorm:"type:text;not null;default:'public';check:chk_channel_visibility,visibility IN ('public','private')" json:"visibility"`
}
