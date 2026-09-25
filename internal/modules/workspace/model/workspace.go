package model

import (
	domain "github.com/PritamDasTharu/conduit/internal/domain/model"
)

type Workspace struct {
	domain.Base
	Name string `gorm:"type:text;not null" json:"name"`
	Slug string `gorm:"type:text;not null" json:"slug"`
}
