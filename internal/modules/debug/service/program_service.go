package service

import (
	"context"
	"slices"

	"github.com/myapp/internal/modules/debug/models"
	"github.com/myapp/internal/modules/debug/repository"
	"gorm.io/gorm"
)

type ProgramService interface {
	DeleteAttByAttIds(context.Context, []uint, *gorm.DB) error
	FindUsedAttByPgmids(context.Context, []uint, *gorm.DB) ([]uint, error)
	FindPgmInfoByPkgId(context.Context, uint, *gorm.DB) ([]models.PgmDataInfo, error)
	BatchCreate(context.Context, uint, []models.PgmDataInfo, *gorm.DB) ([]models.PgmDataInfo, error)
	BatchUpdate(context.Context, uint, []models.PgmDataInfo, *gorm.DB) ([]models.PgmDataInfo, error)
	BatchDelete(context.Context, []uint, *gorm.DB) error
	BatchDeleteByPkgId(context.Context, uint) error
	FindByPkgId(context.Context, uint, *gorm.DB) ([]models.Program, error)
	Create(context.Context, *models.Program) (*models.Program, error)
	Update(context.Context, *models.Program, *gorm.DB) (*models.Program, error)
	Delete(context.Context, uint) error
}

type programService struct {
	Repo   repository.ProgramRepo
	Dbgsvc DebugInfoService
}

func NewProgramService(Repo repository.ProgramRepo, Dbgsvc DebugInfoService) ProgramService {
	return &programService{Repo, Dbgsvc}
}

func (s *programService) DeleteAttByAttIds(ctx context.Context, attids []uint, tx *gorm.DB) error {
	return s.Dbgsvc.DeleteAttByAttIds(ctx, attids, tx)
}

func (s *programService) FindUsedAttByPgmids(ctx context.Context, pgmids []uint, tx *gorm.DB) ([]uint, error) {
	dbgarr, err := s.Dbgsvc.FindByPgmIds(ctx, pgmids, tx)
	if err != nil {
		return nil, err
	}
	res := make([]uint, 0)
	for _, dbg := range dbgarr {
		res = append(res, dbg.Attid)
	}
	return res, nil
}

func (s *programService) FindPgmInfoByPkgId(ctx context.Context, pkgid uint, tx *gorm.DB) ([]models.PgmDataInfo, error) {
	var res []models.PgmDataInfo
	// depend on pkgid query all pgminfo
	pgmres, err := s.Repo.FindByPkgId(ctx, pkgid, tx)
	if err != nil {
		return nil, err
	}
	pgmids := make([]uint, 0)
	for _, pgm := range pgmres {
		pgmids = append(pgmids, pgm.Id)
	}
	dbgres, err := s.Dbgsvc.FindByPgmIds(ctx, pgmids, tx)
	if err != nil {
		return nil, err
	}
	// wrapper program
	for _, pgm := range pgmres {
		dbgres := matchDbgByPgmId(pgm.Id, dbgres)
		tmpPgmData := &models.PgmDataInfo{
			Id:       pgm.Id,
			Name:     pgm.Name,
			Code:     pgm.Code,
			Sort:     pgm.Sort,
			PkgId:    pgm.PackageId,
			DbgArray: dbgres,
		}
		res = append(res, *tmpPgmData)
	}
	slices.SortFunc(res, func(x, y models.PgmDataInfo) int {
		return int(x.Sort) - int(y.Sort)
	})
	return res, nil
}

func matchDbgByPgmId(pgmid uint, dbgarr []models.DbgInfo) []models.DbgInfo {
	var res []models.DbgInfo
	for _, dbg := range dbgarr {
		if dbg.PgmId == pgmid {
			res = append(res, dbg)
		}
	}
	return res
}

func (s *programService) BatchDeleteByPkgId(ctx context.Context, pkgid uint) error {
	// depend on pkgid query all pgminfo
	pgmres, err := s.Repo.FindByPkgId(ctx, pkgid, nil)
	if err != nil {
		return err
	}
	pgmids := make([]uint, 0)
	for _, pgm := range pgmres {
		pgmids = append(pgmids, pgm.Id)
	}
	dbgres, err := s.Dbgsvc.FindByPgmIds(ctx, pgmids, nil)
	if err != nil {
		return err
	}
	// 1. remove all dbginfo
	var (
		dbgids = make([]uint, 0)
	)
	for _, dbg := range dbgres {
		dbgids = append(dbgids, dbg.Id)
	}
	err = s.Dbgsvc.BatchDelete(ctx, dbgids, nil)
	if err != nil {
		return err
	}
	// 2. remove all pgminfo
	err = s.Repo.BatchDelete(ctx, pgmids, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *programService) BatchCreate(ctx context.Context, pkgid uint, d []models.PgmDataInfo, tx *gorm.DB) ([]models.PgmDataInfo, error) {
	// create program, after create program fetch id and create debuginfo
	for _, pgm := range d {
		p := &models.Program{
			Name:      pgm.Name,
			Code:      pgm.Code,
			Sort:      pgm.Sort,
			PackageId: pkgid,
		}
		res, err := s.Repo.Create(ctx, p, tx)
		if err != nil {
			return nil, err
		}
		pgm.Id = res.Id
		// batch create debuginfo
		_, err = s.Dbgsvc.BatchCreate(ctx, pgm.Id, pgm.DbgArray, tx)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (s *programService) BatchUpdate(ctx context.Context, pkgid uint, d []models.PgmDataInfo, tx *gorm.DB) ([]models.PgmDataInfo, error) {
	// create program, after create program fetch id and create debuginfo
	for _, pgm := range d {
		// update program info
		p := &models.Program{
			Id:        pgm.Id,
			Name:      pgm.Name,
			Code:      pgm.Code,
			Sort:      pgm.Sort,
			PackageId: pkgid,
		}
		_, err := s.Update(ctx, p, tx)
		if err != nil {
			return nil, err
		}
		// depend on pgmid query all dbginfo
		res, err := s.Dbgsvc.FindByPgmId(ctx, pgm.Id, tx)
		if err != nil {
			return nil, err
		}
		// calc and handler debuginfo
		addRes, updateRes, deleteRes := calcDiffDbgData(pgm, res)
		if len(addRes) > 0 {
			// create debuginfo
			_, err := s.Dbgsvc.BatchCreate(ctx, pgm.Id, addRes, tx)
			if err != nil {
				return nil, err
			}
		}

		if len(updateRes) > 0 {
			// create debuginfo
			_, err := s.Dbgsvc.BatchUpdate(ctx, pgm.Id, updateRes, tx)
			if err != nil {
				return nil, err
			}
		}

		if len(deleteRes) > 0 {
			// delete debuginfo
			s.Dbgsvc.BatchDelete(ctx, deleteRes, tx)
		}
	}
	return d, nil
}

func (s *programService) BatchDelete(ctx context.Context, pgmids []uint, tx *gorm.DB) error {
	// Query all DebugInfo
	res, err := s.Dbgsvc.FindByPgmIds(ctx, pgmids, tx)
	if err != nil {
		return err
	}
	dbgids := make([]uint, 0)
	for _, dbg := range res {
		dbgids = append(dbgids, dbg.Id)
	}
	// remove all dbginfo by dbgids
	err = s.Dbgsvc.BatchDelete(ctx, dbgids, tx)
	if err != nil {
		return err
	}
	return s.Repo.BatchDelete(ctx, pgmids, tx)
}

func calcDiffDbgData(d models.PgmDataInfo, arr []models.DbgInfo) (addRes []models.DbgInfo, updateRes []models.DbgInfo, deletedRes []uint) {
	addRes = make([]models.DbgInfo, 0)
	updateRes = make([]models.DbgInfo, 0)
	deletedRes = make([]uint, 0)
	for _, dbg := range d.DbgArray {
		if dbg.Id == 0 {
			// create
			addRes = append(addRes, dbg)
		} else {
			if queryDbgIdIsExistDB(dbg.Id, arr) {
				updateRes = append(updateRes, dbg)
			}
		}
	}
	for _, dbg := range arr {
		if !queryDbgIsExistMemory(dbg.Id, d.DbgArray) {
			deletedRes = append(deletedRes, dbg.Id)
		}
	}
	return addRes, updateRes, deletedRes
}

func queryDbgIdIsExistDB(id uint, pgmarr []models.DbgInfo) bool {
	for _, pgm := range pgmarr {
		if id == pgm.Id {
			return true
		}
	}
	return false
}

func queryDbgIsExistMemory(dbgid uint, dbgarray []models.DbgInfo) bool {
	for _, dbg := range dbgarray {
		if dbgid == dbg.Id {
			return true
		}
	}
	return false
}

// find by pkgid
func (s *programService) FindByPkgId(ctx context.Context, id uint, tx *gorm.DB) ([]models.Program, error) {
	return s.Repo.FindByPkgId(ctx, id, tx)
}

func (s *programService) Create(ctx context.Context, d *models.Program) (*models.Program, error) {
	return s.Repo.Create(ctx, d, nil)
}

func (s *programService) Update(ctx context.Context, d *models.Program, tx *gorm.DB) (*models.Program, error) {
	return s.Repo.Update(ctx, d, tx)
}

func (s *programService) Delete(ctx context.Context, id uint) error {
	return s.Repo.Delete(ctx, id)
}
