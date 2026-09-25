package model

import (
	"github.com/google/uuid"
	domain "github.com/PritamDasTharu/conduit/internal/domain/model"
)

type WorkspaceMember struct {
	domain.Base
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	WorkspaceID uuid.UUID `gorm:"type:uuid;not null" json:"workspace_id"`
	RoleID      uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`

	Workspace Workspace `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:RESTRICT" json:"workspace,omitempty"`
}
