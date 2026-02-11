package repositories

import (
	"errors"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/models/scopes"
	"golang_starter_kit_2025/app/repositories/interfaces"

	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) interfaces.PermissionRepositoryInterface {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) Update(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

func (r *permissionRepository) FindByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) List(page, limit int) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	if err := r.db.Model(&models.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Scopes(scopes.PaginateSimple(page, limit)).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

func (r *permissionRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Permission{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *permissionRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Permission{}).Count(&count).Error
	return count, err
}

func (r *permissionRepository) GetAll() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) FindByIDs(ids []uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Find(&permissions, ids).Error
	return permissions, err
}
