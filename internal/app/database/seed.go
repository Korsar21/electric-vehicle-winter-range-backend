package database

import (
	"log"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/model"
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/repository"
)

func Seed() {

	var count int64

	DB.Table("electric_car_loads").Count(&count)

	if count > 0 {
		log.Println("Seed skipped: data already exists")
		return
	}

	user := model.User{
		Name: "Владимир",
	}

	DB.Create(&user)

	loads := []repository.ElectricCarLoad{

		{
			Name:        "Обогреватель салона",
			Description: "Мощная система отопления салона поддерживает комфорт пассажиров при отрицательных температурах.",
			Category:    "Отопление",
			PowerDrawW:  4500,
			Priority:    "Высокий",
			Status:      repository.StatusPublished,
			ImageKey:    "cabin-heater.jpg",
			VideoKey:    "cabin-heater.mp4",
		},

		{
			Name:        "Обогреватель переднего сиденья",
			Description: "Локальный электрический обогрев встроен в подушку и спинку переднего сиденья.",
			Category:    "Отопление",
			PowerDrawW:  150,
			Priority:    "Низкий",
			Status:      repository.StatusPublished,
			ImageKey:    "front-seat-heater.jpg",
			VideoKey:    "front-seat-heater.mp4",
		},

		{
			Name:        "Обогреватель рулевого колеса",
			Description: "Локальный обогрев рулевого колеса повышает комфорт водителя зимой.",
			Category:    "Отопление",
			PowerDrawW:  50,
			Priority:    "Низкий",
			Status:      repository.StatusPublished,
			ImageKey:    "steering-wheel-heater.jpg",
			VideoKey:    "steering-wheel-heater.mp4",
		},

		{
			Name:        "Обогреватель лобового стекла",
			Description: "Электрический обогрев удаляет иней и лёд.",
			Category:    "Обзорность",
			PowerDrawW:  540,
			Priority:    "Высокий",
			Status:      repository.StatusPublished,
			ImageKey:    "windshield-defroster.jpg",
			VideoKey:    "windshield-defroster.mp4",
		},

		{
			Name:        "Обогреватель батареи",
			Description: "Нагреватель системы терморегулирования прогревает тяговую батарею.",
			Category:    "Терморегулирование батареи",
			PowerDrawW:  1000,
			Priority:    "Высокий",
			Status:      repository.StatusPublished,
			ImageKey:    "battery-heater.jpg",
			VideoKey:    "battery-heater.mp4",
		},
		{
			Name:        "Новая нагрузка",
			Description: "Черновик для добавления новой нагрузки",
			Category:    "Другое",
			PowerDrawW:  500,
			Priority:    "Низкий",
			Status:      repository.StatusDraft,
			ImageKey:    "default.jpg",
			VideoKey:    "",
		},
	}

	for _, load := range loads {

		err := DB.Create(&load).Error

		if err != nil {
			log.Println("Seed error:", err)
		}

	}

	log.Println("Seed completed")
}
