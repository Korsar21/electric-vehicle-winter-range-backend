package repository

import (
	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/model"
	"gorm.io/gorm"
)

type ElectricCarLoadStatus string

const (
	StatusDraft     ElectricCarLoadStatus = "draft"
	StatusPublished ElectricCarLoadStatus = "published"
	StatusDeleted   ElectricCarLoadStatus = "deleted"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	return &Repository{
		db: db,
	}, nil
}

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

	if err != nil {
		return r.GetFirstPublishedElectricCarLoad()
	}

	return load, nil
}

func (r *Repository) GetDraftElectricCarLoad(userID uint) (model.ElectricCarLoad, error) {

	var load model.ElectricCarLoad

	err := r.db.
		Where(
			"creator_id = ? AND status = ?",
			userID,
			StatusDraft,
		).
		First(&load).
		Error

	return load, err
}

// CreateElectricCarLoad СЃРѕР·РґР°С‘С‚ РёРјРµРЅРЅРѕ Р§Р•Р РќРћР’РРљ.
// РСЃРїРѕР»СЊР·СѓРµС‚СЃСЏ ORM.
func (r *Repository) CreateElectricCarLoad(load *model.ElectricCarLoad) error {

	load.Status = string(StatusDraft)

	return r.db.Create(load).Error
}

// PublishElectricCarLoad Р·Р°РїРѕР»РЅСЏРµС‚ РїРѕР»СЏ РїСЂРµРґРјРµС‚РЅРѕР№ РѕР±Р»Р°СЃС‚Рё
// Рё РїРµСЂРµРІРѕРґРёС‚ С‡РµСЂРЅРѕРІРёРє РІ published.
// РСЃРїРѕР»СЊР·СѓРµС‚СЃСЏ ORM.
func (r *Repository) PublishElectricCarLoad(
	id int,
	description string,
	category string,
	powerDrawW int,
	priority string,
) error {

	return r.db.
		Model(&model.ElectricCarLoad{}).
		Where(
			"id = ? AND status = ?",
			id,
			StatusDraft,
		).
		Updates(map[string]interface{}{
			"description":  description,
			"category":     category,
			"power_draw_w": powerDrawW,
			"priority":     priority,
			"status":       StatusPublished,
		}).
		Error
}

// DeleteElectricCarLoad вЂ” Р•Р”РРќРЎРўР’Р•РќРќР«Р™ РјРµС‚РѕРґ СѓРґР°Р»РµРЅРёСЏ.
// Р—РґРµСЃСЊ РЅР°РјРµСЂРµРЅРЅРѕ РќР• РёСЃРїРѕР»СЊР·СѓРµС‚СЃСЏ ORM Update/Save.
// Р’С‹РїРѕР»РЅСЏРµС‚СЃСЏ РѕР±С‹С‡РЅС‹Р№ SQL UPDATE.
func (r *Repository) DeleteElectricCarLoad(id int) error {

	return r.db.Exec(
		"UPDATE electric_car_loads SET status = ? WHERE id = ?",
		StatusDeleted,
		id,
	).Error
}

func (r *Repository) GetElectricCarLoadByID(id int) (model.ElectricCarLoad, error) {

	var load model.ElectricCarLoad

	err := r.db.
		Where("id = ?", id).
		First(&load).
		Error

	return load, err
}
