package service

import (
	"context"
	"fmt"

	attmodels "github.com/myapp/internal/modules/attachment/models"
	"github.com/myapp/internal/modules/attachment/service"
	"github.com/myapp/internal/modules/debug/models"
	"github.com/myapp/internal/modules/debug/repository"
	"gorm.io/gorm"
)

type DebugInfoService interface {
	BatchUpdate(context.Context, uint, []models.DbgInfo, *gorm.DB) ([]models.DbgInfo, error)
	DeleteAttByAttIds(context.Context, []uint, *gorm.DB) error
	FindByAttIds(context.Context, []uint) ([]models.DbgInfo, error)
	BatchDelete(context.Context, []uint, *gorm.DB) error
	FindByPgmIds(context.Context, []uint, *gorm.DB) ([]models.DbgInfo, error)
	FindByPgmId(context.Context, uint, *gorm.DB) ([]models.DbgInfo, error)
	BatchCreate(context.Context, uint, []models.DbgInfo, *gorm.DB) ([]models.DbgInfo, error)
	Delete(context.Context, uint) error
}

type debugInfoService struct {
	Repo   repository.DebugInfoRepo
	Attsvc service.AttachmentService
}

func NewDebugInfoService(Repo repository.DebugInfoRepo, Attsvc service.AttachmentService) DebugInfoService {
	return &debugInfoService{Repo, Attsvc}
}

func (s *debugInfoService) BatchUpdate(ctx context.Context, pgmid uint, dbgarray []models.DbgInfo, tx *gorm.DB) ([]models.DbgInfo, error) {
	res := make([]models.DbgInfo, 0)
	// batch create dbginfo
	for _, dbg := range dbgarray {
		d := &models.DebugInfo{
			Id:           dbg.Id,
			LineNo:       dbg.LineNo,
			Attachmentid: dbg.Attid,
			Sort:         dbg.Sort,
			ProgramId:    pgmid,
		}
		dbginfo, err := s.Repo.Update(ctx, d, tx)
		if err != nil {
			return nil, err
		}
		tmpDbg := &models.DbgInfo{
			Id:      dbginfo.Id,
			LineNo:  dbginfo.LineNo,
			Sort:    dbginfo.Sort,
			Attid:   dbg.Id,
			AttName: dbg.AttName,
			AttType: dbg.AttType,
			AttUrl:  dbg.AttUrl,
			PgmId:   dbginfo.ProgramId,
		}
		res = append(res, *tmpDbg)
	}
	return res, nil
}

func (d *debugInfoService) DeleteAttByAttIds(ctx context.Context, attids []uint, tx *gorm.DB) error {
	return d.Attsvc.BatchDeleteByIds(ctx, attids, tx)
}

func (d *debugInfoService) FindByAttIds(ctx context.Context, ids []uint) ([]models.DbgInfo, error) {
	var res []models.DbgInfo
	dbgarray, err := d.Repo.FindByAttIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	attarray, err := d.Attsvc.FindByIds(ctx, ids, nil)
	if err != nil {
		return nil, err
	}
	for _, dbg := range dbgarray {
		att, err := matchAtt(dbg.Attachmentid, attarray)
		if err != nil {
			continue
		}
		tmpDbg := &models.DbgInfo{
			Id:      dbg.Id,
			LineNo:  dbg.LineNo,
			Sort:    dbg.Sort,
			Attid:   att.Id,
			AttName: att.Name,
			AttType: att.Type,
			AttUrl:  att.Url,
		}
		res = append(res, *tmpDbg)
	}
	return res, nil
}

func (s *debugInfoService) FindByPgmIds(ctx context.Context, ids []uint, tx *gorm.DB) ([]models.DbgInfo, error) {
	dbgarr, err := s.Repo.FindByPgmIds(ctx, ids, tx)
	if err != nil {
		return nil, err
	}
	attids := make([]uint, 0)
	for _, dbg := range dbgarr {
		attids = append(attids, dbg.Attachmentid)
	}
	attarr, err := s.Attsvc.FindByIds(ctx, attids, tx)
	if err != nil {
		return nil, err
	}
	return wrapperDbgInfo(dbgarr, attarr)
}

// batch delete
func (s *debugInfoService) BatchDelete(ctx context.Context, ids []uint, tx *gorm.DB) error {
	// need batch remove attachment
	res, err := s.Repo.FindByIds(ctx, ids, tx)
	if err != nil {
		return err
	}
	attids := make([]uint, 0)
	for _, dbg := range res {
		attids = append(attids, dbg.Attachmentid)
	}
	err = s.Attsvc.BatchDeleteByIds(ctx, attids, tx)
	if err != nil {
		return err
	}
	return s.Repo.BatchDelete(ctx, ids, tx)
}

// Update
func (s *debugInfoService) FindByPgmId(ctx context.Context, pgmid uint, tx *gorm.DB) ([]models.DbgInfo, error) {
	var pgmids = make([]uint, 0)
	pgmids = append(pgmids, pgmid)
	dbgarr, err := s.Repo.FindByPgmIds(ctx, pgmids, tx)
	if err != nil {
		return nil, err
	}
	var attids = make([]uint, 0)
	for _, dbg := range dbgarr {
		attids = append(attids, dbg.Attachmentid)
	}
	attarr, err := s.Attsvc.FindByIds(ctx, attids, tx)
	if err != nil {
		return nil, err
	}
	return wrapperDbgInfo(dbgarr, attarr)
}

func wrapperDbgInfo(dbgarr []models.DebugInfo, attarr []attmodels.Attachment) ([]models.DbgInfo, error) {
	res := make([]models.DbgInfo, 0)
	for _, dbg := range dbgarr {
		att, err := matchAtt(dbg.Attachmentid, attarr)
		tmpDbg := &models.DbgInfo{
			Id:      dbg.Id,
			LineNo:  dbg.LineNo,
			Sort:    dbg.Sort,
			Attid:   0,
			AttName: "",
			AttType: "",
			AttUrl:  "",
			PgmId:   dbg.ProgramId,
		}
		if err == nil {
			tmpDbg.Attid = att.Id
			tmpDbg.AttName = att.Name
			tmpDbg.AttType = att.Type
			tmpDbg.AttUrl = att.Url
		}
		res = append(res, *tmpDbg)
	}
	return res, nil
}

func matchAtt(id uint, attarr []attmodels.Attachment) (*attmodels.Attachment, error) {
	for _, att := range attarr {
		if id == att.Id {
			return &att, nil
		}
	}
	return nil, fmt.Errorf("attachment is not exist: %d", id)
}

func (s *debugInfoService) BatchCreate(ctx context.Context, pgmid uint, dbgarray []models.DbgInfo, tx *gorm.DB) ([]models.DbgInfo, error) {
	res := make([]models.DbgInfo, 0)
	// batch create dbginfo
	for _, dbg := range dbgarray {
		d := &models.DebugInfo{
			LineNo:       dbg.LineNo,
			Attachmentid: dbg.Attid,
			Sort:         dbg.Sort,
			ProgramId:    pgmid,
		}
		dbginfo, err := s.Repo.Create(ctx, d, tx)
		if err != nil {
			return nil, err
		}
		tmpDbg := &models.DbgInfo{
			Id:      dbginfo.Id,
			LineNo:  dbginfo.LineNo,
			Sort:    dbginfo.Sort,
			Attid:   dbg.Id,
			AttName: dbg.AttName,
			AttType: dbg.AttType,
			AttUrl:  dbg.AttUrl,
			PgmId:   dbginfo.ProgramId,
		}
		res = append(res, *tmpDbg)
	}
	return res, nil
}

func (s *debugInfoService) Update(ctx context.Context, d *models.DebugInfo) (*models.DebugInfo, error) {
	return s.Repo.Update(ctx, d, nil)
}

func (s *debugInfoService) Delete(ctx context.Context, id uint) error {
	// remove associate attachment
	dbg, err := s.Repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	err = s.Attsvc.DeleteById(ctx, dbg.Attachmentid)
	if err != nil {
		return err
	}
	return s.Repo.Delete(ctx, id)
}
