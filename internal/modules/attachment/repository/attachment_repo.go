package repository

import (
	"context"

	"github.com/myapp/internal/modules/attachment/models"
	"gorm.io/gorm"
)

type AttachmentRepo interface {
	FindByIds(context.Context, []uint, *gorm.DB) ([]models.Attachment, error)
	FindById(context.Context, uint) (*models.Attachment, error)
	Create(context.Context, *models.Attachment) (*models.Attachment, error)
	Update(context.Context, *models.Attachment) (*models.Attachment, error)
	Delete(context.Context, uint) error
}

type attachmentRepo struct {
	DB *gorm.DB
}

func NewAttachmentRepo(DB *gorm.DB) AttachmentRepo {
	DB.AutoMigrate(&models.Attachment{})
	return &attachmentRepo{DB}
}

// Batch Query By Ids
func (r *attachmentRepo) FindByIds(ctx context.Context, ids []uint, tx *gorm.DB) ([]models.Attachment, error) {
	var database *gorm.DB;
	if tx!=nil {
		database = tx
    } else {
		database = r.DB
    }

	var attarr []models.Attachment
	result := database.WithContext(ctx).Where("id IN ?", ids).Find(&attarr)
	if result.Error != nil {
		return nil, result.Error
	}
	return attarr, nil
}

func (r *attachmentRepo) FindById(ctx context.Context, id uint) (*models.Attachment, error) {
	var att *models.Attachment
	result := r.DB.WithContext(ctx).Where("id = ?", id).Find(&att)
	if result.Error != nil {
		return nil, result.Error
	}
	return att, nil
}

// create data
func (r *attachmentRepo) Create(ctx context.Context, a *models.Attachment) (*models.Attachment, error) {
	err := r.DB.WithContext(ctx).Create(a).Error
	if err != nil {
		return nil, err
	}
	return a, nil
}

// update data
func (r *attachmentRepo) Update(ctx context.Context, a *models.Attachment) (*models.Attachment, error) {
	var aTmp *models.Attachment
	r.DB.WithContext(ctx).Find(&aTmp, a.Id)

	aTmp.Name = a.Name
	aTmp.Url = a.Url

	r.DB.WithContext(ctx).Save(&aTmp)
	return aTmp, nil
}

// delete data
func (r *attachmentRepo) Delete(ctx context.Context, id uint) error {
	var tmp *models.Attachment
	r.DB.WithContext(ctx).Delete(&tmp, id)
	return nil
}
