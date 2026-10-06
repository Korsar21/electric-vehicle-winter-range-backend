package model

type Like struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID            uint `json:"user_id"`
	ElectricCarLoadID uint `json:"electric_car_load_id"`
}
