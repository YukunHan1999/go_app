package service

import (
	"context"

	"github.com/myapp/internal/modules/attachment/service"
	"github.com/myapp/internal/modules/envzip/models"
	"github.com/myapp/internal/modules/envzip/repository"
)

type EnvZipService interface {
	FindDir(context.Context) ([]models.EnvDataDTO, error)
	FindByPId(context.Context, uint) ([]models.EnvDataDTO, error)
	Create(context.Context, *models.EnvDataDTO) (*models.EnvData, error)
	Update(context.Context, *models.EnvDataDTO) (*models.EnvData, error)
	Delete(context.Context, uint) error
}

type envZipService struct {
	Repo   repository.EnvZipRepo
	Attsvc service.AttachmentService
}

func NewEnvZipService(repo repository.EnvZipRepo, Attsvc service.AttachmentService) EnvZipService {
	return &envZipService{repo, Attsvc}
}

// Create implements [envZipService].
func (s *envZipService) Create(ctx context.Context, data *models.EnvDataDTO) (*models.EnvData, error) {
	envData := &models.EnvData{
		Name:         data.Name,
		Remark:       data.Remark,
		Attachmentid: data.Attachmentid,
		Pid:          data.Pid,
	}
	res, err := s.Repo.Create(ctx, envData)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Delete implements [envZipService].
func (s *envZipService) Delete(ctx context.Context, id uint) error {
	data, err := s.Repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if data.Attachmentid != 0 {
		// delete file and remove attachment
		err := s.Attsvc.DeleteById(ctx, uint(data.Attachmentid))
		if err != nil {
			return err
		}
	}
	err = s.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// FindByPId implements [envZipService].
func (s *envZipService) FindByPId(ctx context.Context, pid uint) ([]models.EnvDataDTO, error) {
	res, err := s.Repo.FindByPId(ctx, pid)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// FindDir implements [envZipService].
func (s *envZipService) FindDir(ctx context.Context) ([]models.EnvDataDTO, error) {
	res, err := s.Repo.FindDir(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Update implements [envZipService].
func (s *envZipService) Update(ctx context.Context, data *models.EnvDataDTO) (*models.EnvData, error) {
	envData, err := s.Repo.FindById(ctx, data.Id)
	if err != nil {
		return nil, err
	}
	if data.Attachmentid != envData.Attachmentid {
		// remove old env attachment
		err := s.Attsvc.DeleteById(ctx, envData.Attachmentid)
		if err != nil {
			return nil, err
		}
	}
	// envData
	envData.Name = data.Name
	envData.Remark = data.Remark
	envData.Attachmentid = data.Attachmentid
	envData.Pid = data.Pid
	// update new env attachment
	res, err := s.Repo.Update(ctx, envData)
	if err != nil {
		return nil, err
	}
	return res, nil
}
