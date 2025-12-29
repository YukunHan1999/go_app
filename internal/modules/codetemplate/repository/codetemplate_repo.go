package repository

import (
	"context"
	"fmt"

	"github.com/myapp/internal/modules/codetemplate/models"

	"gorm.io/gorm"
)

type CodeTemplateRepo interface {
	Count(context.Context) (int64, error)
	Create(context.Context, *models.Codetemplate) (*models.Codetemplate, error)
	Update(context.Context, string, string, string, string, string) (*models.Codetemplate, error)
	FindAll(context.Context, int, int) ([]models.Codetemplate, error)
	Delete(context.Context, string) error
}

type codeTemplateRepo struct {
	DB *gorm.DB
}

func NewCodeTemplateRepo(DB *gorm.DB) CodeTemplateRepo {
	DB.AutoMigrate(&models.Codetemplate{})
	return &codeTemplateRepo{DB}
}

func (r *codeTemplateRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.DB.WithContext(ctx).Model(&models.Codetemplate{}).Count(&count) // 统计 User 表的总行数
	if result.Error != nil {
		// 处理错误
		return 0, fmt.Errorf("query total count throw exception: %s", result.Error.Error())
	}
	return count, nil
}

// create data
func (r *codeTemplateRepo) Create(ctx context.Context, c *models.Codetemplate) (*models.Codetemplate, error) {
	err := r.DB.WithContext(ctx).Create(c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}

// update data
func (r *codeTemplateRepo) Update(ctx context.Context, id, lang, name, code, tips string) (*models.Codetemplate, error) {
	var tmp *models.Codetemplate
	r.DB.WithContext(ctx).First(&tmp, id)

	tmp.Lang = lang
	tmp.Name = name
	tmp.Code = code
	tmp.Tips = tips

	r.DB.WithContext(ctx).Save(&tmp)
	return tmp, nil
}

// remove data
func (r *codeTemplateRepo) Delete(ctx context.Context, id string) error {
	var tmp *models.Codetemplate
	r.DB.WithContext(ctx).Delete(&tmp, id)
	return nil
}

// query data
func (r *codeTemplateRepo) FindAll(ctx context.Context, size, page int) ([]models.Codetemplate, error) {
	var codeTmps []models.Codetemplate
	skipSize := (page - 1) * size
	r.DB.WithContext(ctx).Limit(size).Offset(skipSize).Find(&codeTmps)
	return codeTmps, nil
}
