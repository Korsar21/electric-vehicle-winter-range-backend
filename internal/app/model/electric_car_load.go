package model

import "time"

type ElectricCarLoad struct {
	ID uint `gorm:"primaryKey"`

	Name string

	Description string

	Category string

	PowerDrawW int

	Priority string

	Status string

	LikeCount int

	ImageKey string

	VideoKey string

	CreatedAt time.Time

	UpdatedAt time.Time
}
