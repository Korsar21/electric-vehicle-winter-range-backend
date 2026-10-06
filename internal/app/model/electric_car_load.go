package model

import "time"

type ElectricCarLoad struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Status string `json:"status"`

	ImageKey string `json:"image_key"`
	VideoKey string `json:"video_key"`

	Category   string `json:"category"`
	PowerDrawW int    `json:"power_draw_w"`
	Priority   string `json:"priority"`

	LikeCount int `json:"like_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FormedAt  time.Time `json:"formed_at"`

	CreatorID uint `json:"creator_id"`
	Creator   User `json:"creator"`

	Likes []User `gorm:"many2many:likes;" json:"likes"`
}
