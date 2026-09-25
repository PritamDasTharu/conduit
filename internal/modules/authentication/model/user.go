package model

import (
	domain "github.com/PritamDasTharu/conduit/internal/domain/model"
)

type User struct {
	domain.Base
	Name         string `gorm:"type:text;not null" json:"name"`
	Email        string `gorm:"type:text;not null" json:"email"`
	PasswordHash string `gorm:"type:text;not null" json:"-"`
}
