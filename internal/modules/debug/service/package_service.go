package service

import (
	"context"

	"github.com/myapp/internal/modules/debug/models"
	"github.com/myapp/internal/modules/debug/repository"
	"gorm.io/gorm"
)

type PackageService interface {
	RemoveUnusedAtt(context.Context, uint, []uint, *gorm.DB) error
	LoadSingle(context.Context, uint) (*models.PkgDataInfo, error)
	DeleteById(context.Context, uint) error
	UpdatePkg(context.Context, *models.PkgDataInfo, *gorm.DB) (*models.PkgDataInfo, error)
	CreatePkg(context.Context, *models.PkgDataInfo, *gorm.DB) (*models.PkgDataInfo, error)
	Create(context.Context, *models.Package, *gorm.DB) (*models.Package, error)
	QueryDataByDirId(context.Context, uint) ([]models.Package, error)
}

type packageService struct {
	DB *gorm.DB

	Repo   repository.PackageRepo
	PgmSvc ProgramService
}

func NewPackageService(DB *gorm.DB, Repo repository.PackageRepo, PgmSvc ProgramService) PackageService {
	return &packageService{DB, Repo, PgmSvc}
}

func (s *packageService) RemoveUnusedAtt(ctx context.Context, id uint, allattids []uint, tx *gorm.DB) error {
	var database *gorm.DB
	if tx != nil {
		database = tx
	} else {
		database = s.DB
	}
	if len(allattids) == 0 {
		return nil
	}
	err := database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		removeids := make([]uint, 0)
		if id == 0 {
			removeids = allattids
		} else {
			pgmarr, err := s.PgmSvc.FindByPkgId(ctx, id, tx)
			if err != nil {
				return err
			}
			pgmids := make([]uint, 0)
			for _, pgm := range pgmarr {
				pgmids = append(pgmids, pgm.Id)
			}
			usedAttids, err := s.PgmSvc.FindUsedAttByPgmids(ctx, pgmids, tx)
			if err != nil {
				return err
			}
			for _, attid := range allattids {
				if !containAttId(usedAttids, attid) {
					removeids = append(removeids, attid)
				}
			}
		}
		// removeids
		if len(removeids) > 0 {
			err := s.PgmSvc.DeleteAttByAttIds(ctx, removeids, tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func containAttId(arr []uint, id uint) bool {
	for _, tmp := range arr {
		if id == tmp {
			return true
		}
	}
	return false
}

func (s *packageService) LoadSingle(ctx context.Context, id uint) (*models.PkgDataInfo, error) {
	p, err := s.Repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	pgmarr, err := s.PgmSvc.FindPgmInfoByPkgId(ctx, p.Id, nil)
	if err != nil {
		return nil, err
	}
	res := &models.PkgDataInfo{
		Id:           p.Id,
		Name:         p.Name,
		Description:  p.Description,
		DirectoryId:  p.DirectoryId,
		ProgramArray: pgmarr,
	}
	return res, nil
}

func (s *packageService) DeleteById(ctx context.Context, id uint) error {
	// batch delete pgm
	err := s.PgmSvc.BatchDeleteByPkgId(ctx, id)
	if err != nil {
		return err
	}
	err = s.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *packageService) CreatePkg(ctx context.Context, d *models.PkgDataInfo, tx *gorm.DB) (*models.PkgDataInfo, error) {
	// create package
	var pkg = &models.Package{
		Id:          d.Id,
		Name:        d.Name,
		Description: d.Description,
		DirectoryId: d.DirectoryId,
	}
	pkgres, err := s.Create(ctx, pkg, tx)
	if err != nil {
		return nil, err
	}
	d.Id = pkg.Id
	// calc need delete, update prorgaminfo
	_, err = s.PgmSvc.FindByPkgId(ctx, pkgres.Id, tx)
	if err != nil {
		return nil, err
	}
	// need create
	s.PgmSvc.BatchCreate(ctx, pkgres.Id, d.ProgramArray, tx)
	return d, nil
}

// Update Pkg
func (s *packageService) UpdatePkg(ctx context.Context, d *models.PkgDataInfo, tx *gorm.DB) (*models.PkgDataInfo, error) {
	// calc need delete, update prorgaminfo
	pgmres, err := s.PgmSvc.FindByPkgId(ctx, d.Id, tx)
	if err != nil {
		return nil, err
	}

	pkg := &models.Package{
		Id:          d.Id,
		Name:        d.Name,
		Description: d.Description,
		DirectoryId: d.DirectoryId,
	}
	// pkg save
	pkgres, err := s.Repo.Update(ctx, pkg, tx)
	if err != nil {
		return nil, err
	}
	d.Id = pkgres.Id
	d.Name = pkgres.Name
	d.Description = pkgres.Description
	d.DirectoryId = pkgres.DirectoryId

	// query diff by id
	addRes, updatedRes, deleteRes := calcDiffPgmData(d, pgmres)

	if len(addRes) > 0 {
		// need create
		s.PgmSvc.BatchCreate(ctx, d.Id, addRes, tx)
	}

	if len(updatedRes) > 0 {
		// need updated
		s.PgmSvc.BatchUpdate(ctx, d.Id, updatedRes, tx)
	}

	if len(deleteRes) > 0 {
		// need delete and batch delete files
		s.PgmSvc.BatchDelete(ctx, deleteRes, tx)
	}
	return d, nil
}

// query need update info
func calcDiffPgmData(d *models.PkgDataInfo, arr []models.Program) (addRes []models.PgmDataInfo, updatedRes []models.PgmDataInfo, deleteRes []uint) {
	addRes = make([]models.PgmDataInfo, 0)
	updatedRes = make([]models.PgmDataInfo, 0)
	deleteRes = make([]uint, 0)
	for _, pgm := range d.ProgramArray {
		if pgm.Id == 0 {
			// needCreate
			addRes = append(addRes, pgm)
		} else {
			// current In database need update
			if queryPgmIdIsExistDB(pgm.Id, arr) {
				// update
				updatedRes = append(updatedRes, pgm)
			}
		}
	}

	// exist db but not memory
	for _, pp := range arr {
		if !queryPgmIdIsExistMemory(pp.Id, d.ProgramArray) {
			// need delete data
			deleteRes = append(deleteRes, pp.Id)
		}
	}

	return addRes, updatedRes, deleteRes
}

// exist in d and dbginfo
func queryPgmIdIsExistDB(id uint, pgmarr []models.Program) bool {
	for _, pgm := range pgmarr {
		if id == pgm.Id {
			return true
		}
	}
	return false
}

// exist in d and dbginfo
func queryPgmIdIsExistMemory(id uint, pgmarr []models.PgmDataInfo) bool {
	for _, pgm := range pgmarr {
		if id == pgm.Id {
			return true
		}
	}
	return false
}

// Create Package
func (s *packageService) Create(ctx context.Context, pkg *models.Package, tx *gorm.DB) (*models.Package, error) {
	return s.Repo.Create(ctx, pkg, tx)
}

// QueryByDirId
func (s *packageService) QueryDataByDirId(ctx context.Context, id uint) ([]models.Package, error) {
	return s.Repo.QueryDataByDirId(ctx, id)
}
