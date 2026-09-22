package model

type Like struct {
	ID uint `gorm:"primaryKey"`

	UserID uint

	ElectricCarLoadID uint
}
