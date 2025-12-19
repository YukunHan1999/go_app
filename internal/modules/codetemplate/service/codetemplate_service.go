package service

import (
	"context"

	"github.com/myapp/internal/modules/codetemplate/models"
	"github.com/myapp/internal/modules/codetemplate/repository"
)

type CodeTemplateService interface {
	Count(context.Context) (int64, error)
	Delete(context.Context, string) error
	Create(context.Context, string, string, string, string) (*models.Codetemplate, error)
	Update(context.Context, string, string, string, string, string) (*models.Codetemplate, error)
	FindAll(context.Context, int, int) ([]models.Codetemplate, error)
}

type codeTemplateService struct {
	Repo repository.CodeTemplateRepo
}

func NewCodeTemplateService(repo repository.CodeTemplateRepo) CodeTemplateService {
	return &codeTemplateService{Repo: repo}
}

func (s *codeTemplateService) Count(ctx context.Context) (int64, error) {
	count, err := s.Repo.Count(ctx)
	return count, err
}

func (s *codeTemplateService) Create(ctx context.Context, lang, name, code, tips string) (*models.Codetemplate, error) {
	c := &models.Codetemplate{Lang: lang, Name: name, Code: code, Tips: tips}
	return s.Repo.Create(ctx, c)
}

func (s *codeTemplateService) Update(ctx context.Context, id, lang, name, code, tips string) (*models.Codetemplate, error) {
	return s.Repo.Update(ctx, id, lang, name, code, tips)
}

func (s *codeTemplateService) FindAll(ctx context.Context, size, page int) ([]models.Codetemplate, error) {
	return s.Repo.FindAll(ctx, size, page)
}

func (s *codeTemplateService) Delete(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}
