package model

import "time"

type User struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name     string `json:"name"`
	Password string `json:"-"`

	CreatedAt time.Time `json:"created_at"`
}
