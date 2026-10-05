package model

import "time"

type ElectricCarLoad struct {
	ID uint `gorm:"primaryKey"`

	Name        string
	Description string

	Status string

	ImageKey string
	VideoKey string

	Category   string
	PowerDrawW int
	Priority   string

	CreatedAt time.Time
	FormedAt  time.Time

	CreatorID uint
	Creator   User

	LikeCount int

	Likes []User `gorm:"many2many:likes;"`
}
