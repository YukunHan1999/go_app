package service

import (
	"context"

	"github.com/myapp/internal/modules/debug/models"
	"github.com/myapp/internal/modules/debug/repository"
)

type DirectoryService interface {
	DeleteById(context.Context, uint) error
	Find(context.Context) ([]models.Directory, error)
	Create(context.Context, *models.Directory) (*models.Directory, error)
	Update(context.Context, *models.Directory) (*models.Directory, error)
	QueryByDirId(context.Context, uint) ([]models.Directory, error)
}

type directoryService struct {
	Repo repository.DirectoryRepo
}

func NewDirectoryService(Repo repository.DirectoryRepo) DirectoryService {
	return &directoryService{Repo}
}

func (s *directoryService) DeleteById(ctx context.Context, id uint) error {
	return s.Repo.Delete(ctx, id)
}

// Create dir
func (s *directoryService) Create(ctx context.Context, d *models.Directory) (*models.Directory, error) {
	return s.Repo.Create(ctx, d)
}

// Update dir
func (s *directoryService) Update(ctx context.Context, d *models.Directory) (*models.Directory, error) {
	return s.Repo.Update(ctx, d)
}

// Query All Dir
func (s *directoryService) Find(ctx context.Context) ([]models.Directory, error) {
	return s.Repo.Find(ctx)
}

// QueryByDirId
func (s *directoryService) QueryByDirId(ctx context.Context, pid uint) ([]models.Directory, error) {
	return s.Repo.QueryByDirId(ctx, pid)
}
