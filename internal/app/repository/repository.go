package repository

import (
	"errors"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/model"
	"gorm.io/gorm"
)

type ElectricCarLoadStatus string

const (
	StatusDraft     ElectricCarLoadStatus = "draft"
	StatusPublished ElectricCarLoadStatus = "published"
	StatusDeleted   ElectricCarLoadStatus = "deleted"
)

var ErrElectricCarLoadNotFound = errors.New("electric car load not found")
var ErrElectricCarLoadForbidden = errors.New("electric car load belongs to another user")
var ErrUserNotFound = errors.New("user not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	return &Repository{
		db: db,
	}, nil
}

// =========================
// ELECTRIC CAR LOADS
// =========================

func (r *Repository) GetPublishedElectricCarLoads(maxPowerW int) ([]model.ElectricCarLoad, error) {
	var loads []model.ElectricCarLoad

	err := r.db.
		Where(
			"status = ? AND power_draw_w <= ?",
			StatusPublished,
			maxPowerW,
		).
		Order("id ASC").
		Find(&loads).
		Error

	return loads, err
}

func (r *Repository) GetPublishedElectricCarLoadByID(id int) (model.ElectricCarLoad, error) {
	var load model.ElectricCarLoad

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

func (r *Repository) GetFirstPublishedElectricCarLoad() (model.ElectricCarLoad, error) {
	var load model.ElectricCarLoad

	err := r.db.
		Where("status = ?", StatusPublished).
		Order("id ASC").
		First(&load).
		Error

	return load, err
}

func (r *Repository) GetNextPublishedElectricCarLoad(id int) (model.ElectricCarLoad, error) {
	var load model.ElectricCarLoad

	err := r.db.
		Where(
			"id > ? AND status = ?",
			id,
			StatusPublished,
		).
		Order("id ASC").
		First(&load).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.GetFirstPublishedElectricCarLoad()
	}

	return load, err
}

func (r *Repository) GetDraftElectricCarLoad(userID uint) (model.ElectricCarLoad, error) {
	var load model.ElectricCarLoad

	err := r.db.
		Where(
			"creator_id = ? AND status = ?",
			userID,
			StatusDraft,
		).
		Order("id ASC").
		First(&load).
		Error

	return load, err
}

func (r *Repository) CreateElectricCarLoad(load *model.ElectricCarLoad) error {
	load.Status = string(StatusDraft)

	return r.db.Create(load).Error
}

func (r *Repository) PublishElectricCarLoad(
	id int,
	userID uint,
	description string,
	category string,
	powerDrawW int,
	priority string,
) error {
	var load model.ElectricCarLoad

	err := r.db.
		Where(
			"id = ? AND creator_id = ? AND status = ?",
			id,
			userID,
			StatusDraft,
		).
		First(&load).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrElectricCarLoadNotFound
	}

	if err != nil {
		return err
	}

	return r.db.
		Model(&load).
		Updates(map[string]interface{}{
			"description":  description,
			"category":     category,
			"power_draw_w": powerDrawW,
			"priority":     priority,
			"status":       StatusPublished,
			"formed_at":    gorm.Expr("NOW()"),
		}).
		Error
}

func (r *Repository) DeleteElectricCarLoad(id int, userID uint) error {
	var load model.ElectricCarLoad

	err := r.db.
		Where(
			"id = ? AND creator_id = ? AND status != ?",
			id,
			userID,
			StatusDeleted,
		).
		First(&load).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrElectricCarLoadNotFound
	}

	if err != nil {
		return err
	}

	return r.db.
		Model(&load).
		Update("status", StatusDeleted).
		Error
}

func (r *Repository) GetElectricCarLoadByID(id int) (model.ElectricCarLoad, error) {
	var load model.ElectricCarLoad

	err := r.db.
		Where("id = ?", id).
		First(&load).
		Error

	return load, err
}

// =========================
// LIKES
// =========================

func (r *Repository) IsLiked(
	userID uint,
	electricCarLoadID uint,
) (bool, error) {
	var count int64

	err := r.db.
		Model(&model.Like{}).
		Where(
			"user_id = ? AND electric_car_load_id = ?",
			userID,
			electricCarLoadID,
		).
		Count(&count).
		Error

	return count > 0, err
}

func (r *Repository) SetLike(
	userID uint,
	electricCarLoadID uint,
	like int,
) error {
	if like != 0 && like != 1 {
		return errors.New("like must be 0 or 1")
	}

	tx := r.db.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	var existing model.Like

	err := tx.
		Where(
			"user_id = ? AND electric_car_load_id = ?",
			userID,
			electricCarLoadID,
		).
		First(&existing).
		Error

	if like == 1 {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&model.Like{
				UserID:            userID,
				ElectricCarLoadID: electricCarLoadID,
			}).Error; err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.
				Model(&model.ElectricCarLoad{}).
				Where("id = ?", electricCarLoadID).
				UpdateColumn(
					"like_count",
					gorm.Expr("like_count + 1"),
				).
				Error; err != nil {
				tx.Rollback()
				return err
			}
		} else if err != nil {
			tx.Rollback()
			return err
		}
	}

	if like == 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return tx.Commit().Error
		}

		if err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Delete(&existing).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.
			Model(&model.ElectricCarLoad{}).
			Where("id = ? AND like_count > 0", electricCarLoadID).
			UpdateColumn(
				"like_count",
				gorm.Expr("like_count - 1"),
			).
			Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// =========================
// USERS
// =========================

func (r *Repository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByName(name string) (model.User, error) {
	var user model.User

	err := r.db.
		Where("name = ?", name).
		First(&user).
		Error

	return user, err
}
