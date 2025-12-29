package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/myapp/internal/modules/attachment/models"
	"github.com/myapp/internal/modules/attachment/repository"
)

type AttachmentService interface {
	FindByIds(context.Context, []uint) ([]models.Attachment, error)
	BatchDeleteByIds(context.Context, []uint) error
	DeleteById(context.Context, uint) error
	Create(context.Context, *models.Attachment) (*models.Attachment, error)
}

type attachmentService struct {
	Repo repository.AttachmentRepo
}

func NewAttachmentService(Repo repository.AttachmentRepo) AttachmentService {
	return &attachmentService{Repo}
}

func (s *attachmentService) FindByIds(ctx context.Context, ids []uint) ([]models.Attachment, error) {
	return s.Repo.FindByIds(ctx, ids)
}

// Batch Remove file and attachment
func (s *attachmentService) BatchDeleteByIds(ctx context.Context, ids []uint) error {
	attarr, err := s.Repo.FindByIds(ctx, ids)
	if err != nil {
		return err
	}
	filepaths := make([]string, 0)
	// attarr
	for _, att := range attarr {
		filepaths = append(filepaths,
			filepath.Join("uploads", att.Url))
	}
	err = removeFiles(filepaths)
	if err != nil {
		return err
	}
	return nil
}

// Remove file and Delete
func (s *attachmentService) DeleteById(ctx context.Context, id uint) error {
	att, err := s.Repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if (len(att.Url) > 0) {
		path := filepath.Join("uploads", att.Url)
		_ = removeFile(path)
	}
	err = s.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func removeFiles(filePaths []string) error {
	for _, filePath := range filePaths {
		err := os.Remove(filePath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove file %s: %v", filePath, err)
		}
		if err == nil {
			fmt.Printf("File %s deleted\n", filePath)
		}
	}
	return nil
}

func removeFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("File %s does not exist\n", filePath)
			return nil // Or handle as needed
		}
		return fmt.Errorf("failed to remove file %s: %v", filePath, err)
	}
	fmt.Printf("File %s deleted successfully\n", filePath)
	return nil
}

func (s *attachmentService) Create(ctx context.Context, a *models.Attachment) (*models.Attachment, error) {
	return s.Repo.Create(ctx, a)
}
