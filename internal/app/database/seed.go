package database

import (
	"log"
	"time"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/model"
)

func Seed() {

	var count int64

	DB.Table("electric_car_loads").Count(&count)

	if count > 0 {
		log.Println("Seed skipped: data already exists")
		return
	}

	user := model.User{
		Name:     "ладимир",
		Password: "123456",
	}

	if err := DB.Create(&user).Error; err != nil {
		log.Println("User seed error:", err)
		return
	}

	now := time.Now()

	loads := []model.ElectricCarLoad{

		{
			Name:        "богреватель салона",
			Description: "ощная система отопления салона поддерживает комфорт пассажиров при отрицательных температурах. Система быстро прогревает внутреннее пространство автомобиля и помогает поддерживать комфортную температуру во время поездки. ри этом она является одной из наиболее энергозатратных зимних нагрузок.",
			Status:      "published",
			ImageKey:    "cabin-heater.jpg",
			VideoKey:    "cabin-heater.mp4",
			Category:    "топление",
			PowerDrawW:  4500,
			Priority:    "ысокий",
			CreatedAt:   now,
			FormedAt:    now,
			CreatorID:   user.ID,
			LikeCount:   128,
		},

		{
			Name:        "богреватель переднего сиденья",
			Description: "окальный электрический обогрев встроен в подушку и спинку переднего сиденья. н позволяет быстро прогреть место водителя и пассажира без необходимости полностью прогревать салон. Такая нагрузка потребляет значительно меньше энергии по сравнению с общим отоплением.",
			Status:      "published",
			ImageKey:    "front-seat-heater.jpg",
			VideoKey:    "front-seat-heater.mp4",
			Category:    "топление",
			PowerDrawW:  150,
			Priority:    "изкий",
			CreatedAt:   now,
			FormedAt:    now,
			CreatorID:   user.ID,
			LikeCount:   84,
		},

		{
			Name:        "богреватель рулевого колеса",
			Description: "окальный обогрев рулевого колеса повышает комфорт водителя зимой. агревательный элемент быстро прогревает поверхность рулевого колеса и позволяет водителю пользоваться автомобилем без ощущения холодного руля. отребление энергии при этом остаётся небольшим.",
			Status:      "published",
			ImageKey:    "steering-wheel-heater.jpg",
			VideoKey:    "steering-wheel-heater.mp4",
			Category:    "топление",
			PowerDrawW:  50,
			Priority:    "изкий",
			CreatedAt:   now,
			FormedAt:    now,
			CreatorID:   user.ID,
			LikeCount:   61,
		},

		{
			Name:        "богреватель лобового стекла",
			Description: "лектрический обогрев удаляет иней и лёд с поверхности лобового стекла. Система помогает быстро восстановить обзорность автомобиля после стоянки при низкой температуре. агрузка относится к важным зимним функциям автомобиля.",
			Status:      "published",
			ImageKey:    "windshield-defroster.jpg",
			VideoKey:    "windshield-defroster.mp4",
			Category:    "бзорность",
			PowerDrawW:  540,
			Priority:    "ысокий",
			CreatedAt:   now,
			FormedAt:    now,
			CreatorID:   user.ID,
			LikeCount:   97,
		},

		{
			Name:        "богреватель батареи",
			Description: "агреватель системы терморегулирования прогревает тяговую батарею. оддержание оптимальной температуры аккумулятора особенно важно при эксплуатации электромобиля зимой, поскольку низкая температура влияет на доступную мощность и запас хода. абота системы помогает подготовить батарею к эффективной эксплуатации.",
			Status:      "published",
			ImageKey:    "battery-heater.jpg",
			VideoKey:    "battery-heater.mp4",
			Category:    "Терморегулирование батареи",
			PowerDrawW:  1000,
			Priority:    "ысокий",
			CreatedAt:   now,
			FormedAt:    now,
			CreatorID:   user.ID,
			LikeCount:   73,
		},

		{
			Name:        "далённая нагрузка",
			Description: "та карточка используется для проверки логического удаления услуги. осле изменения статуса на deleted карточка не должна отображаться среди опубликованных нагрузок.",
			Status:      "deleted",
			ImageKey:    "default.jpg",
			VideoKey:    "default.mp4",
			Category:    "ругое",
			PowerDrawW:  700,
			Priority:    "изкий",
			CreatedAt:   now,
			FormedAt:    now,
			CreatorID:   user.ID,
			LikeCount:   0,
		},

		{
			Name:        "овая нагрузка",
			Description: "ерновик для добавления новой нагрузки. осле заполнения необходимых данных карточку можно будет опубликовать.",
			Status:      "draft",
			ImageKey:    "",
			VideoKey:    "",
			Category:    "",
			PowerDrawW:  0,
			Priority:    "",
			CreatedAt:   now,
			FormedAt:    now,
			CreatorID:   user.ID,
			LikeCount:   0,
		},
	}

	for _, load := range loads {

		if err := DB.Create(&load).Error; err != nil {
			log.Println("Load seed error:", err)
		}
	}

	log.Println("Seed completed")
}
