package models

import "time"

// CommonModelFields ...
type CommonModelFields struct {
	ID        uint       `gorm:"primary_key;auto_increment" json:"id"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
