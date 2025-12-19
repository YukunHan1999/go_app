package repository

import (
	"context"

	"github.com/myapp/internal/modules/debug/models"
	"gorm.io/gorm"
)

type DirectoryRepo interface {
	Find(context.Context) ([]models.Directory, error)
	QueryByDirId(context.Context, uint) ([]models.Directory, error)
	Create(context.Context, *models.Directory) (*models.Directory, error)
	Update(context.Context, *models.Directory) (*models.Directory, error)
	Delete(context.Context, uint) error
}

type directoryRepo struct {
	DB *gorm.DB
}

func NewDirectoryRepo(DB *gorm.DB) DirectoryRepo {
	DB.AutoMigrate(&models.Directory{})
	return &directoryRepo{DB}
}

// QueryDir
func (r *directoryRepo) Find(ctx context.Context) ([]models.Directory, error) {
	var dirs []models.Directory
	r.DB.Find(&dirs)
	return dirs, nil
}

// QueryByDirId
func (r *directoryRepo) QueryByDirId(ctx context.Context, pId uint) ([]models.Directory, error) {
	var dirs []models.Directory
	r.DB.Where("parent_id = ?", pId).Find(&dirs)
	return dirs, nil
}

// Create data
func (r *directoryRepo) Create(ctx context.Context, d *models.Directory) (*models.Directory, error) {
	err := r.DB.Create(d).Error
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Update Data
func (r *directoryRepo) Update(ctx context.Context, d *models.Directory) (*models.Directory, error) {
	var tmp *models.Directory
	r.DB.Find(&tmp, d.Id)

	tmp.Name = d.Name
	tmp.Remark = d.Remark
	tmp.ParentId = d.ParentId

	r.DB.Save(&tmp)
	return tmp, nil
}

func (r *directoryRepo) Delete(ctx context.Context, id uint) error {
	var tmp *models.Directory
	r.DB.Delete(&tmp, id)
	return nil
}
