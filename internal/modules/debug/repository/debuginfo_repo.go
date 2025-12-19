package repository

import (
	"context"

	"github.com/myapp/internal/modules/debug/models"
	"gorm.io/gorm"
)

type DebugInfoRepo interface {
	FindByAttIds(context.Context, []uint) ([]models.DebugInfo, error)
	BatchDelete(context.Context, []uint) error
	FindById(context.Context, uint) (*models.DebugInfo, error)
	FindByIds(context.Context, []uint) ([]models.DebugInfo, error)
	FindByPgmIds(context.Context, []uint) ([]models.DebugInfo, error)
	BatchDeleteByPgmIds(context.Context, []uint) error
	Create(context.Context, *models.DebugInfo) (*models.DebugInfo, error)
	Update(context.Context, *models.DebugInfo) (*models.DebugInfo, error)
	Delete(context.Context, uint) error
}

type debugInfoRepo struct {
	DB *gorm.DB
}

func NewDebugInfoRepo(DB *gorm.DB) DebugInfoRepo {
	DB.AutoMigrate(&models.DebugInfo{})
	return &debugInfoRepo{DB}
}

func (r *debugInfoRepo) FindByAttIds(ctx context.Context, ids []uint) ([]models.DebugInfo, error) {
	var debugarr []models.DebugInfo
	result := r.DB.WithContext(ctx).Where("attachmentid IN ?", ids).Find(&debugarr)
	if result.Error != nil {
		return nil, result.Error
	}
	return debugarr, nil
}

// Batch Query debuginfo By id arry
func (r *debugInfoRepo) FindByIds(ctx context.Context, ids []uint) ([]models.DebugInfo, error) {
	var debugarr []models.DebugInfo
	result := r.DB.WithContext(ctx).Where("id IN ?", ids).Find(&debugarr)
	if result.Error != nil {
		return nil, result.Error
	}
	return debugarr, nil
}

func (r *debugInfoRepo) FindById(ctx context.Context, id uint) (*models.DebugInfo, error) {
	var debug *models.DebugInfo
	result := r.DB.WithContext(ctx).Where("id = ?", id).Find(&debug)
	if result.Error != nil {
		return nil, result.Error
	}
	return debug, nil
}

func (r *debugInfoRepo) FindByPgmIds(ctx context.Context, ids []uint) ([]models.DebugInfo, error) {
	var res []models.DebugInfo
	result := r.DB.WithContext(ctx).Where("program_id IN ?", ids).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return res, nil
}

func (r *debugInfoRepo) BatchDeleteByPgmIds(ctx context.Context, d []uint) error {
	res := r.DB.WithContext(ctx).Where("program_id IN ?", d).Delete(&models.DebugInfo{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *debugInfoRepo) BatchDelete(ctx context.Context, d []uint) error {
	res := r.DB.WithContext(ctx).Where("id in ?", d).Delete(&models.DebugInfo{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// create data
func (r *debugInfoRepo) Create(ctx context.Context, info *models.DebugInfo) (*models.DebugInfo, error) {
	err := r.DB.WithContext(ctx).Create(info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

// update data
func (r *debugInfoRepo) Update(ctx context.Context, info *models.DebugInfo) (*models.DebugInfo, error) {
	var debugInfoTmp *models.DebugInfo
	r.DB.WithContext(ctx).Find(&debugInfoTmp, info.Id)

	debugInfoTmp.ProgramId = info.ProgramId
	debugInfoTmp.LineNo = info.LineNo
	debugInfoTmp.Attachmentid = info.Attachmentid

	r.DB.WithContext(ctx).Save(&debugInfoTmp)
	return debugInfoTmp, nil
}

// delete data
func (r *debugInfoRepo) Delete(ctx context.Context, id uint) error {
	r.DB.WithContext(ctx).Delete(&models.DebugInfo{}, id)
	return nil
}
