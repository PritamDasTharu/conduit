package model

import (
	"github.com/google/uuid"
	domain "github.com/PritamDasTharu/conduit/internal/domain/model"
)

type Role struct {
	domain.Base
	WorkspaceID uuid.UUID    `gorm:"type:uuid;not null" json:"workspace_id"`
	Name        string       `gorm:"type:text;not null" json:"name"`
	Workspace   Workspace    `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:RESTRICT" json:"workspace,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
}

type RolePermission struct {
	RoleID       uuid.UUID  `gorm:"type:uuid;primaryKey" json:"role_id"`
	PermissionID uuid.UUID  `gorm:"type:uuid;primaryKey" json:"permission_id"`
	Role         Role       `gorm:"foreignKey:RoleID;constraint:OnDelete:RESTRICT" json:"-"`
	Permission   Permission `gorm:"foreignKey:PermissionID;constraint:OnDelete:RESTRICT" json:"-"`
}
