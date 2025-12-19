package repository

import (
	"context"

	"github.com/myapp/internal/modules/debug/models"
	"gorm.io/gorm"
)

type PackageRepo interface {
	FindById(context.Context, uint) (*models.Package, error)
	QueryDataByDirId(context.Context, uint) ([]models.Package, error)
	Create(context.Context, *models.Package) (*models.Package, error)
	Update(context.Context, *models.Package) (*models.Package, error)
	Delete(context.Context, uint) error
}

type packageRepo struct {
	DB *gorm.DB
}

func NewPackageRepo(DB *gorm.DB) PackageRepo {
	DB.AutoMigrate(&models.Package{})
	return &packageRepo{DB}
}

func (r *packageRepo) FindById(ctx context.Context, id uint) (*models.Package, error) {
	var p *models.Package
	res := r.DB.WithContext(ctx).Where("id = ?", id).Find(&p)
	if res.Error != nil {
		return nil, res.Error
	}
	return p, nil
}

// QueryDataByDirId
func (r *packageRepo) QueryDataByDirId(ctx context.Context, pId uint) ([]models.Package, error) {
	// SELECT * FROM directory WHERE directoryid = pId;
	var packages []models.Package
	r.DB.WithContext(ctx).Where("directory_id = ?", pId).Find(&packages)
	return packages, nil
}

// create data
func (r *packageRepo) Create(ctx context.Context, p *models.Package) (*models.Package, error) {
	err := r.DB.WithContext(ctx).Create(p).Error
	if err != nil {
		return nil, err
	}
	return p, nil
}

// update data
func (r *packageRepo) Update(ctx context.Context, pack *models.Package) (*models.Package, error) {
	var packTmp *models.Package
	r.DB.WithContext(ctx).Find(&packTmp, pack.Id)

	packTmp.Name = pack.Name
	packTmp.Description = pack.Description
	packTmp.DirectoryId = pack.DirectoryId

	r.DB.WithContext(ctx).Save(&packTmp)
	return packTmp, nil
}

func (r *packageRepo) Delete(ctx context.Context, id uint) error {
	var packageTmp *models.Package
	r.DB.WithContext(ctx).Delete(&packageTmp, id)
	return nil
}
