package repository

import "fmt"

type ElectricCarLoadStatus string

const (
	StatusDraft     ElectricCarLoadStatus = "draft"
	StatusPublished ElectricCarLoadStatus = "published"
	StatusDeleted   ElectricCarLoadStatus = "deleted"
)

type ElectricCarLoad struct {
	ID          int
	Name        string
	Description string
	Category    string
	PowerDrawW  int
	Priority    string
	Status      ElectricCarLoadStatus
	LikeUserIDs []int
	ImageKey    string
	VideoKey    string
}

type Repository struct {
	electricCarLoads []ElectricCarLoad
}

func generateLikeUserIDs(startID, count int) []int {
	ids := make([]int, count)

	for i := 0; i < count; i++ {
		ids[i] = startID + i
	}

	return ids
}

func NewRepository() (*Repository, error) {
	loads := []ElectricCarLoad{
		{
			ID:          1,
			Name:        "Обогреватель салона",
			Description: "Мощная система отопления салона поддерживает комфорт пассажиров при отрицательных температурах. В суровых зимних условиях она может стать одной из крупнейших дополнительных электрических нагрузок и заметно снизить доступный запас хода.", Category: "Отопление",
			PowerDrawW:  4500,
			Priority:    "Высокий",
			Status:      StatusPublished,
			LikeUserIDs: generateLikeUserIDs(1000, 128),
			ImageKey:    "cabin-heater.jpg",
			VideoKey:    "cabin-heater.mp4",
		},
		{
			ID:          2,
			Name:        "Обогреватель переднего сиденья",
			Description: "Локальный электрический обогрев встроен в подушку и спинку переднего сиденья. Он обеспечивает непосредственный комфорт пассажира и потребляет значительно меньше энергии, чем отопление всего салона.", Category: "Отопление",
			PowerDrawW:  150,
			Priority:    "Низкий",
			Status:      StatusPublished,
			LikeUserIDs: generateLikeUserIDs(2000, 84),
			ImageKey:    "front-seat-heater.jpg",
			VideoKey:    "front-seat-heater.mp4",
		},
		{
			ID:          4,
			Name:        "Обогреватель рулевого колеса",
			Description: "Локальный электрический обогрев обода рулевого колеса повышает комфорт водителя в холодную погоду. Низкое энергопотребление делает его более экономичным, чем повышение температуры во всём салоне.", Category: "Отопление",
			PowerDrawW:  50,
			Priority:    "Низкий",
			Status:      StatusPublished,
			LikeUserIDs: generateLikeUserIDs(3000, 61),
			ImageKey:    "steering-wheel-heater.jpg",
			VideoKey:    "steering-wheel-heater.mp4",
		},
		{
			ID:          7,
			Name:        "Обогреватель лобового стекла",
			Description: "Электрический обогрев лобового стекла удаляет иней и лёд и помогает сохранять обзорность при эксплуатации автомобиля зимой. Во время работы система увеличивает дополнительное потребление электроэнергии.", Category: "Обзорность",
			PowerDrawW:  540,
			Priority:    "Высокий",
			Status:      StatusPublished,
			LikeUserIDs: generateLikeUserIDs(4000, 97),
			ImageKey:    "windshield-defroster.jpg",
			VideoKey:    "windshield-defroster.mp4",
		},
		{
			ID:          9,
			Name:        "Обогреватель батареи",
			Description: "Нагреватель системы терморегулирования прогревает тяговую батарею при низкой температуре окружающей среды. Поддержание подходящей температуры батареи улучшает возможности зарядки и отдачу мощности, но потребляет дополнительную энергию.", Category: "Терморегулирование батареи",
			PowerDrawW:  1000,
			Priority:    "Высокий",
			Status:      StatusPublished,
			LikeUserIDs: generateLikeUserIDs(5000, 73),
			ImageKey:    "battery-heater.jpg",
			VideoKey:    "battery-heater.mp4",
		},
		{
			ID:          12,
			Name:        "Обогреватель заднего сиденья",
			Description: "Система обогрева заднего сиденья предназначена для локального повышения комфорта пассажиров в холодную погоду.",
			Category:    "Отопление",
			PowerDrawW:  150,
			Priority:    "Низкий",
			Status:      StatusDraft,
			LikeUserIDs: []int{},
			ImageKey:    "rear-seat-heater.jpg",
			VideoKey:    "rear-seat-heater.mp4",
		},
		{
			ID:          15,
			Name:        "Обогреватель заднего стекла",
			Description: "Electrical rear-window heating used to remove frost and condensation.",
			Category:    "Обзорность",
			PowerDrawW:  240,
			Priority:    "Высокий",
			Status:      StatusDeleted,
			LikeUserIDs: generateLikeUserIDs(6000, 34),
			ImageKey:    "rear-window-defroster.jpg",
			VideoKey:    "rear-window-defroster.mp4",
		},
	}

	if len(loads) == 0 {
		return nil, fmt.Errorf("electric car load collection is empty")
	}

	return &Repository{
		electricCarLoads: loads,
	}, nil
}

func (r *Repository) GetPublishedElectricCarLoads(maxPowerW int) ([]ElectricCarLoad, error) {
	result := make([]ElectricCarLoad, 0)

	for _, load := range r.electricCarLoads {
		if load.Status != StatusPublished {
			continue
		}

		if load.PowerDrawW > maxPowerW {
			continue
		}

		result = append(result, load)
	}

	return result, nil
}

func (r *Repository) GetPublishedElectricCarLoadByID(id int) (ElectricCarLoad, error) {
	for _, load := range r.electricCarLoads {
		if load.ID == id && load.Status == StatusPublished {
			return load, nil
		}
	}

	return ElectricCarLoad{}, fmt.Errorf(
		"published electric car load with id %d not found",
		id,
	)
}

func (r *Repository) GetDraftElectricCarLoad() (ElectricCarLoad, error) {
	for _, load := range r.electricCarLoads {
		if load.Status == StatusDraft {
			return load, nil
		}
	}

	return ElectricCarLoad{}, fmt.Errorf(
		"draft electric car load not found",
	)
}

func (r *Repository) GetNextPublishedElectricCarLoad(id int) (ElectricCarLoad, error) {
	currentFound := false

	for _, load := range r.electricCarLoads {
		if load.Status != StatusPublished {
			continue
		}

		if currentFound {
			return load, nil
		}

		if load.ID == id {
			currentFound = true
		}
	}

	if !currentFound {
		return ElectricCarLoad{}, fmt.Errorf(
			"published vehicle auxiliary load with id %d not found",
			id,
		)
	}

	return r.GetFirstPublishedElectricCarLoad()
}

func (r *Repository) GetFirstPublishedElectricCarLoad() (ElectricCarLoad, error) {
	for _, load := range r.electricCarLoads {
		if load.Status == StatusPublished {
			return load, nil
		}
	}

	return ElectricCarLoad{}, fmt.Errorf(
		"published vehicle auxiliary load not found",
	)
}

func (r *Repository) HasNextPublishedElectricCarLoad(id int) bool {
	for _, load := range r.electricCarLoads {
		if load.ID == id && load.Status == StatusPublished {
			return true
		}
	}

	return false
}
