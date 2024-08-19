package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	ProjectId string    `json:"project_id" gorm:"type:uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (r *Repository) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New().String()
	return
}
