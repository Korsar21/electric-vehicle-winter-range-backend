package model

type User struct {
	ID uint `gorm:"primaryKey"`

	Name     string
	Password string
}
