package repository

import (
	"gorm.io/gorm"
)

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
	LikeCount   int
	ImageKey    string
	VideoKey    string
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository, error) {

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetPublishedElectricCarLoads(maxPowerW int) ([]ElectricCarLoad, error) {

	var loads []ElectricCarLoad

	err := r.db.
		Where(
			"status = ? AND power_draw_w <= ?",
			StatusPublished,
			maxPowerW,
		).
		Find(&loads).
		Error

	return loads, err
}

func (r *Repository) GetPublishedElectricCarLoadByID(id int) (ElectricCarLoad, error) {

	var load ElectricCarLoad

	err := r.db.
		Where(
			"id = ? AND status = ?",
			id,
			StatusPublished,
		).
		First(&load).
		Error

	return load, err
}

func (r *Repository) GetDraftElectricCarLoad() (ElectricCarLoad, error) {

	var load ElectricCarLoad

	err := r.db.
		Where(
			"status = ?",
			StatusDraft,
		).
		First(&load).
		Error

	return load, err
}

func (r *Repository) GetNextPublishedElectricCarLoad(id int) (ElectricCarLoad, error) {

	var load ElectricCarLoad

	err := r.db.
		Where(
			"id > ? AND status = ?",
			id,
			StatusPublished,
		).
		First(&load).
		Error

	if err != nil {
		return r.GetFirstPublishedElectricCarLoad()
	}

	return load, nil
}

func (r *Repository) GetFirstPublishedElectricCarLoad() (ElectricCarLoad, error) {

	var load ElectricCarLoad

	err := r.db.
		Where(
			"status = ?",
			StatusPublished,
		).
		First(&load).
		Error

	return load, err
}

func (r *Repository) CreateElectricCarLoad(load ElectricCarLoad) error {

	if load.Status == "" {
		load.Status = StatusPublished
	}

	return r.db.Create(&load).Error
}

func (r *Repository) DeleteElectricCarLoad(id int) error {

	return r.db.
		Model(&ElectricCarLoad{}).
		Where("id = ?", id).
		Update(
			"status",
			StatusDeleted,
		).
		Error
}

func (r *Repository) GetElectricCarLoadByID(id int) (ElectricCarLoad, error) {

	var load ElectricCarLoad

	err := r.db.
		Where("id = ?", id).
		First(&load).
		Error

	return load, err
}

func (r *Repository) UpdateElectricCarLoad(load ElectricCarLoad) error {

	return r.db.
		Save(&load).
		Error
}

func (r *Repository) LikeElectricCarLoad(id int) error {

	return r.db.
		Model(&ElectricCarLoad{}).
		Where("id = ?", id).
		UpdateColumn(
			"like_count",
			gorm.Expr("like_count + ?", 1),
		).
		Error
}
