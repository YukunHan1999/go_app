package repository

import (
	"context"

	"github.com/myapp/internal/modules/debug/models"
	"gorm.io/gorm"
)

type ProgramRepo interface {
	BatchDelete(context.Context, []uint) error
	FindByPkgId(context.Context, uint) ([]models.Program, error)
	Create(context.Context, *models.Program) (*models.Program, error)
	Update(context.Context, *models.Program) (*models.Program, error)
	Delete(context.Context, uint) error
}

type programRepo struct {
	DB *gorm.DB
}

func NewProgramRepo(DB *gorm.DB) ProgramRepo {
	DB.AutoMigrate(&models.Program{})
	return &programRepo{DB}
}

func (r *programRepo) BatchDelete(ctx context.Context, d []uint) error {
	var arr []models.Program
	res := r.DB.WithContext(ctx).Where("id IN ?", d).Delete(&arr)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *programRepo) FindByPkgId(ctx context.Context, id uint) ([]models.Program, error) {
	var pgms []models.Program
	res := r.DB.WithContext(ctx).Where("package_id = ?", id).Find(&pgms)
	if res.Error != nil {
		return nil, res.Error
	}
	return pgms, nil
}

// create a new program
func (r *programRepo) Create(ctx context.Context, p *models.Program) (*models.Program, error) {
	err := r.DB.WithContext(ctx).Create(p).Error
	if err != nil {
		return nil, err
	}
	return p, nil
}

// update program info
func (r *programRepo) Update(ctx context.Context, m *models.Program) (*models.Program, error) {
	var programTmp *models.Program
	r.DB.WithContext(ctx).First(&programTmp, m.Id)

	programTmp.Name = m.Name
	programTmp.Code = m.Code
	programTmp.PackageId = m.PackageId
	programTmp.Sort = m.Sort

	r.DB.WithContext(ctx).Save(&programTmp)
	return programTmp, nil
}

// delete program info
func (r *programRepo) Delete(ctx context.Context, id uint) error {
	var codeTmp *models.Program
	r.DB.WithContext(ctx).Delete(&codeTmp, id)
	return nil
}
