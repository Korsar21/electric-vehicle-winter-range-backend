package repository

import "fmt"

type AuxiliaryLoadStatus string

const (
	StatusDraft     AuxiliaryLoadStatus = "draft"
	StatusPublished AuxiliaryLoadStatus = "published"
	StatusDeleted   AuxiliaryLoadStatus = "deleted"
)

type VehicleAuxiliaryLoad struct {
	ID          int
	Name        string
	Description string
	Category    string
	PowerDrawW  int
	Status      AuxiliaryLoadStatus
	LikeUserIDs []int
	ImageKey    string
	VideoKey    string
}

type Repository struct {
	vehicleAuxiliaryLoads []VehicleAuxiliaryLoad
}

func generateLikeUserIDs(startID, count int) []int {
	ids := make([]int, count)

	for i := 0; i < count; i++ {
		ids[i] = startID + i
	}

	return ids
}

func NewRepository() (*Repository, error) {
	loads := []VehicleAuxiliaryLoad{
		{
			ID:          1,
			Name:        "Cabin Heater",
			Description: "High-power cabin heating system used to maintain passenger comfort in sub-zero temperatures. In severe winter conditions it can become one of the largest auxiliary electrical loads and noticeably reduce the available driving range.", Category: "Heating",
			PowerDrawW:  4500,
			Status:      StatusPublished,
			LikeUserIDs: generateLikeUserIDs(1000, 128),
			ImageKey:    "cabin-heater.jpg",
			VideoKey:    "cabin-heater.mp4",
		},
		{
			ID:          2,
			Name:        "Front Seat Heater",
			Description: "Localized electric heating built into the front seat cushion and backrest. It provides direct passenger comfort while consuming substantially less electrical power than heating the entire vehicle cabin.", Category: "Heating",
			PowerDrawW:  150,
			Status:      StatusPublished,
			LikeUserIDs: generateLikeUserIDs(2000, 84),
			ImageKey:    "front-seat-heater.jpg",
			VideoKey:    "front-seat-heater.mp4",
		},
		{
			ID:          4,
			Name:        "Steering Wheel Heater",
			Description: "Localized electric heating around the steering wheel rim improves driver comfort in cold weather. Its relatively low power demand makes it more energy-efficient than increasing whole-cabin temperature.", Category: "Heating",
			PowerDrawW:  50,
			Status:      StatusPublished,
			LikeUserIDs: generateLikeUserIDs(3000, 61),
			ImageKey:    "steering-wheel-heater.jpg",
			VideoKey:    "steering-wheel-heater.mp4",
		},
		{
			ID:          7,
			Name:        "Windshield Defroster",
			Description: "Electrical windshield heating removes frost and ice and helps maintain forward visibility during winter operation. The system increases auxiliary electrical consumption while it is active.", Category: "Visibility",
			PowerDrawW:  540,
			Status:      StatusPublished,
			LikeUserIDs: generateLikeUserIDs(4000, 97),
			ImageKey:    "windshield-defroster.jpg",
			VideoKey:    "windshield-defroster.mp4",
		},
		{
			ID:          9,
			Name:        "Battery Heater",
			Description: "Battery thermal-management heater warms the traction battery at low ambient temperatures. Maintaining a suitable battery temperature supports charging and power performance but consumes additional stored energy.", Category: "Battery Thermal Management",
			PowerDrawW:  1000,
			Status:      StatusPublished,
			LikeUserIDs: generateLikeUserIDs(5000, 73),
			ImageKey:    "battery-heater.jpg",
			VideoKey:    "battery-heater.mp4",
		},
		{
			ID:          12,
			Name:        "Rear Seat Heater",
			Description: "Rear seat heating system intended to provide localized passenger comfort in cold weather.",
			Category:    "Heating",
			PowerDrawW:  150,
			Status:      StatusDraft,
			LikeUserIDs: []int{},
			ImageKey:    "rear-seat-heater.jpg",
			VideoKey:    "rear-seat-heater.mp4",
		},
		{
			ID:          15,
			Name:        "Rear Window Defroster",
			Description: "Electrical rear-window heating used to remove frost and condensation.",
			Category:    "Visibility",
			PowerDrawW:  240,
			Status:      StatusDeleted,
			LikeUserIDs: generateLikeUserIDs(6000, 34),
			ImageKey:    "rear-window-defroster.jpg",
			VideoKey:    "rear-window-defroster.mp4",
		},
	}

	if len(loads) == 0 {
		return nil, fmt.Errorf("vehicle auxiliary load collection is empty")
	}

	return &Repository{
		vehicleAuxiliaryLoads: loads,
	}, nil
}

func (r *Repository) GetPublishedVehicleAuxiliaryLoads(maxPowerW int) ([]VehicleAuxiliaryLoad, error) {
	result := make([]VehicleAuxiliaryLoad, 0)

	for _, load := range r.vehicleAuxiliaryLoads {
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

func (r *Repository) GetPublishedVehicleAuxiliaryLoadByID(id int) (VehicleAuxiliaryLoad, error) {
	for _, load := range r.vehicleAuxiliaryLoads {
		if load.ID == id && load.Status == StatusPublished {
			return load, nil
		}
	}

	return VehicleAuxiliaryLoad{}, fmt.Errorf(
		"published vehicle auxiliary load with id %d not found",
		id,
	)
}

func (r *Repository) GetDraftVehicleAuxiliaryLoad() (VehicleAuxiliaryLoad, error) {
	for _, load := range r.vehicleAuxiliaryLoads {
		if load.Status == StatusDraft {
			return load, nil
		}
	}

	return VehicleAuxiliaryLoad{}, fmt.Errorf(
		"draft vehicle auxiliary load not found",
	)
}

func (r *Repository) GetNextPublishedVehicleAuxiliaryLoad(id int) (VehicleAuxiliaryLoad, error) {
	currentFound := false

	for _, load := range r.vehicleAuxiliaryLoads {
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
		return VehicleAuxiliaryLoad{}, fmt.Errorf(
			"published vehicle auxiliary load with id %d not found",
			id,
		)
	}

	return VehicleAuxiliaryLoad{}, fmt.Errorf(
		"there is no next published vehicle auxiliary load after id %d",
		id,
	)
}

func (r *Repository) GetFirstPublishedVehicleAuxiliaryLoad() (VehicleAuxiliaryLoad, error) {
	for _, load := range r.vehicleAuxiliaryLoads {
		if load.Status == StatusPublished {
			return load, nil
		}
	}

	return VehicleAuxiliaryLoad{}, fmt.Errorf(
		"published vehicle auxiliary load not found",
	)
}

func (r *Repository) HasNextPublishedVehicleAuxiliaryLoad(id int) bool {
	currentFound := false

	for _, load := range r.vehicleAuxiliaryLoads {
		if load.Status != StatusPublished {
			continue
		}

		if currentFound {
			return true
		}

		if load.ID == id {
			currentFound = true
		}
	}

	return false
}
