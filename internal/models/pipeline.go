package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Pipeline struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey"`
	Name         string    `json:"name"`
	RepositoryId string    `json:"repository_id" gorm:"type:uuid"`
	Path         string    `json:"path"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (p *Pipeline) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New().String()
	return
}
