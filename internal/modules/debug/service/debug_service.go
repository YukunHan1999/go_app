package service

import (
	"context"

	attsvc "github.com/myapp/internal/modules/attachment/service"
	"github.com/myapp/internal/modules/debug/models"
	"gorm.io/gorm"
)

type DebugService interface {
	ClearGarbageAtt(context.Context, uint, []uint) error
	FindSinglePkg(context.Context, uint) (*models.PkgDataInfo, error)
	CreatePkg(context.Context, *models.AddPkgDto) (*models.AddPkgDto, error)
	DeletePkgById(context.Context, uint) error
	DeleteDirById(context.Context, uint) error
	UpdatePkg(context.Context, *models.AddPkgDto) (*models.AddPkgDto, error)
	UpdateDir(context.Context, *models.DirDataInfo) (*models.DirDataInfo, error)
	CreateDir(context.Context, *models.DirData) (*models.DirData, error)
	FetchDirData(context.Context) ([]models.DirDataInfo, error)
	FetchPkgData(context.Context, uint) ([]models.DirDataInfo, error)
}

type debugService struct {
	DB     *gorm.DB
	Dbgsvc DebugInfoService
	Dirsvc DirectoryService
	Pkgsvc PackageService
	AttSvc attsvc.AttachmentService
}

func NewDebugService(DB *gorm.DB, Dbgsvc DebugInfoService, Dirsvc DirectoryService, Pkgsvc PackageService, AttSvc attsvc.AttachmentService) DebugService {
	return &debugService{DB, Dbgsvc, Dirsvc, Pkgsvc, AttSvc}
}

func (s *debugService) ClearGarbageAtt(ctx context.Context, pkgid uint, attGarbage []uint) error {
	return s.Pkgsvc.RemoveUnusedAtt(ctx, pkgid, attGarbage)
}

func (s *debugService) FindSinglePkg(ctx context.Context, id uint) (*models.PkgDataInfo, error) {
	// query pkg data
	return s.Pkgsvc.LoadSingle(ctx, id)
}

// create package info
func (s *debugService) CreatePkg(ctx context.Context, d *models.AddPkgDto) (*models.AddPkgDto, error) {
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// create pkg
		pkg, err := s.Pkgsvc.CreatePkg(ctx, d.Pkg)
		if err != nil {
			return err
		}
		d.Pkg = pkg

		return s.Pkgsvc.RemoveUnusedAtt(ctx, d.Pkg.Id, d.AttIds)
	})
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (s *debugService) DeleteDirById(ctx context.Context, id uint) error {
	return s.Dirsvc.DeleteById(ctx, id)
}

// delete pkg
func (s *debugService) DeletePkgById(ctx context.Context, id uint) error {
	return s.Pkgsvc.DeleteById(ctx, id)
}

// update pkginfo
func (s *debugService) UpdatePkg(ctx context.Context, d *models.AddPkgDto) (*models.AddPkgDto, error) {
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updatedPkg, err := s.Pkgsvc.UpdatePkg(ctx, d.Pkg)
		if err != nil {
			return err
		}
		d.Pkg = updatedPkg
		return s.Pkgsvc.RemoveUnusedAtt(ctx, d.Pkg.Id, d.AttIds)
	})
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Update Directory
func (s *debugService) UpdateDir(ctx context.Context, d *models.DirDataInfo) (*models.DirDataInfo, error) {
	var r *models.DirDataInfo
	dir := &models.Directory{
		Id:       d.Id,
		Name:     d.Name,
		Remark:   d.Description,
		ParentId: d.ParentId,
	}
	// Create Dir
	dir, err := s.Dirsvc.Update(ctx, dir)
	if err != nil {
		return nil, err
	}
	r = &models.DirDataInfo{
		Id:          dir.Id,
		Name:        dir.Name,
		Description: dir.Remark,
		ParentId:    dir.ParentId,
		CreatedAt:   dir.CreatedAt,
		UpdatedAt:   dir.UpdatedAt,
	}
	return r, nil
}

// Query package and directory Data
func (s *debugService) FetchPkgData(ctx context.Context, pId uint) ([]models.DirDataInfo, error) {
	// Query program and directory by parentId
	pArray, err := s.Pkgsvc.QueryDataByDirId(ctx, pId)
	if err != nil {
		return nil, err
	}
	dArray, err := s.Dirsvc.QueryByDirId(ctx, pId)
	if err != nil {
		return nil, err
	}
	r := make([]models.DirDataInfo, 0, 10)
	for _, p := range pArray {
		r = append(r,
			models.DirDataInfo{Id: p.Id, Name: p.Name, Description: p.Description, IsFile: 1, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt, ParentId: p.DirectoryId},
		)
	}

	for _, d := range dArray {
		r = append(r,
			models.DirDataInfo{Id: d.Id, Name: d.Name, Description: d.Remark, IsFile: 0, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt, ParentId: d.ParentId},
		)
	}
	return r, nil
}

// Query Directory Data
func (s *debugService) FetchDirData(ctx context.Context) ([]models.DirDataInfo, error) {
	dArray, err := s.Dirsvc.Find(ctx)
	if err != nil {
		return nil, err
	}
	r := make([]models.DirDataInfo, 0, 10)
	for _, d := range dArray {
		r = append(r,
			models.DirDataInfo{Id: d.Id, Name: d.Name, Description: d.Remark, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt, ParentId: d.ParentId},
		)
	}
	return r, nil
}

func (s *debugService) CreateDir(ctx context.Context, d *models.DirData) (*models.DirData, error) {
	var r *models.DirData
	dir := &models.Directory{
		Name:     d.Name,
		Remark:   d.Description,
		ParentId: d.ParentId,
	}
	// Create Dir
	dir, err := s.Dirsvc.Create(ctx, dir)
	if err != nil {
		return nil, err
	}
	r = &models.DirData{
		Id:          dir.Id,
		Name:        dir.Name,
		Description: dir.Remark,
		ParentId:    dir.ParentId,
	}
	return r, nil
}
