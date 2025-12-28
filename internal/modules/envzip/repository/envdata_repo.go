package repository

import (
	"context"

	"github.com/myapp/internal/modules/envzip/models"
	"gorm.io/gorm"
)

type EnvZipRepo interface {
	FindDir(context.Context) ([]models.EnvDataDTO, error)
	FindByPId(context.Context, uint) ([]models.EnvDataDTO, error)
	FindById(context.Context, uint) (*models.EnvData, error)
	Create(context.Context, *models.EnvData) (*models.EnvData, error)
	Update(context.Context, *models.EnvData) (*models.EnvData, error)
	Delete(context.Context, uint) error
}

type envZipRepo struct {
	DB *gorm.DB
}

func NewEnvZipRepo(DB *gorm.DB) EnvZipRepo {
	DB.AutoMigrate(&models.EnvData{})
	return &envZipRepo{DB}
}

// Create implements [EnvZipRepo].
func (r *envZipRepo) Create(ctx context.Context, data *models.EnvData) (*models.EnvData, error) {
	err := r.DB.WithContext(ctx).Create(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Delete implements [EnvZipRepo].
func (r *envZipRepo) Delete(ctx context.Context, id uint) error {
	var res *models.EnvData
	result := r.DB.WithContext(ctx).Delete(&res, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByPId implements [EnvZipRepo].
func (r *envZipRepo) FindByPId(ctx context.Context, pid uint) ([]models.EnvDataDTO, error) {
	var res []models.EnvDataDTO
	result := r.DB.WithContext(ctx).Raw(`SELECT t1.id, t1.name, t1.remark, t1.attachmentid, t2.url attachmenturl, t1.created_at, t1.updated_at, t1.pid
	FROM env_data t1
	LEFT JOIN attachments t2 ON t1.attachmentid = t2.id
	WHERE t1.pid = ?
	`, pid).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return res, nil
}

func (r *envZipRepo) FindById(ctx context.Context, id uint) (*models.EnvData, error) {
	var res *models.EnvData
	result := r.DB.WithContext(ctx).Where("id = ?", id).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return res, nil
}

// FindDir implements [EnvZipRepo].
func (r *envZipRepo) FindDir(ctx context.Context) ([]models.EnvDataDTO, error) {
	var res []models.EnvDataDTO
	result := r.DB.WithContext(ctx).Raw(`SELECT t1.id, t1.name, t1.remark, t1.attachmentid, t2.url attachmenturl, t1.created_at, t1.updated_at, t1.pid
	FROM env_data t1
	LEFT JOIN attachments t2 ON t1.attachmentid = t2.id
	WHERE t1.attachmentid = ?
	`, 0).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return res, nil
}

// Update implements [EnvZipRepo].
func (r *envZipRepo) Update(ctx context.Context, data *models.EnvData) (*models.EnvData, error) {
	err := r.DB.WithContext(ctx).Updates(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
