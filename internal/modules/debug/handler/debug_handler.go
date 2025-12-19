package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/myapp/internal/common"
	"github.com/myapp/internal/modules/debug/models"
	"github.com/myapp/internal/modules/debug/service"
)

type DebugHandler struct {
	s service.DebugService
}

func NewDirectoryHandler(s service.DebugService) *DebugHandler {
	return &DebugHandler{s}
}

// find pkg
func (h *DebugHandler) FindPackageById(c *gin.Context) {
	pId := c.Param("pid")
	id, err := strconv.Atoi(pId)
	if err != nil {
		common.Fail(c, 50000, "pid can't covert to uint, please check param!")
		return
	}
	if id < 0 {
		common.Fail(c, 50000, "pid not is negative!")
		return
	}
	ppId := uint(id)
	pkgData, err := h.s.FindSinglePkg(c.Request.Context(), ppId)
	if err != nil {
		common.Fail(c, 50000, "Query error: "+err.Error())
		return
	}
	common.Success().Append("data", pkgData).End(c)
}

// add pkg, program, attachment
func (h *DebugHandler) CreatePkg(c *gin.Context) {
	var data *models.AddPkgDto
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		common.Fail(c, 50000, "parse param error")
		return
	}
	_, err = h.s.CreatePkg(c.Request.Context(), data)
	if err != nil {
		common.Fail(c, 50000, "handler error")
		return
	}
	common.Success().Append("isSuccess", true).End(c)
}

// update pkg info, program, attachment, info
func (h *DebugHandler) UpdatePkg(c *gin.Context) {
	var data *models.AddPkgDto
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		common.Fail(c, 50000, "parse param error")
		return
	}
	_, err = h.s.UpdatePkg(c.Request.Context(), data)
	if err != nil {
		common.Fail(c, 50000, "handler error")
		return
	}
	common.Success().Append("isSuccess", true).End(c)
}

func (h *DebugHandler) ClearPkgInfo(c *gin.Context) {
	var data *models.CleakPkgDto
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		common.Fail(c, 50000, "parse param error")
		return
	}
	err = h.s.ClearGarbageAtt(c.Request.Context(), data.Id, data.AttIds)
	if err != nil {
		common.Fail(c, 50000, "handler error")
		return
	}
	common.Success().Append("isSuccess", true).End(c)
}

// query directory and package by pId
func (h *DebugHandler) FetchPkgOrDirData(c *gin.Context) {
	pId := c.Param("pid")
	id, err := strconv.Atoi(pId)
	if err != nil {
		common.Fail(c, 50000, "pid can't covert to uint, please check param!")
		return
	}
	if id < 0 {
		common.Fail(c, 50000, "pid not is negative!")
		return
	}
	ppId := uint(id)
	info, err := h.s.FetchPkgData(c.Request.Context(), ppId)
	if err != nil {
		common.Fail(c, 50000, "Query error: "+err.Error())
		return
	}
	common.Success().Append("data", info).End(c)
}

func (h *DebugHandler) DeletePkgOrDir(c *gin.Context) {
	var data *models.RemovePkgOrDir
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		common.Fail(c, 50000, "parse param error")
		return
	}
	if data.IsFile == 0 {
		err = h.s.DeleteDirById(c.Request.Context(), data.Id)
	} else {
		err = h.s.DeletePkgById(c.Request.Context(), data.Id)
	}
	if err != nil {
		common.Fail(c, 50000, "Delete error:"+err.Error())
		return
	}
	common.Success().Append("isSuccess", true).End(c)
}

// =================================================Dir

// query all directory
func (h *DebugHandler) FetchAllDirData(c *gin.Context) {
	r, err := h.s.FetchDirData(c.Request.Context())
	if err != nil {
		common.Fail(c, 50000, "Query error:"+err.Error())
		return
	}
	common.Success().Append("data", r).End(c)
}

// Create Directory or Package
func (h *DebugHandler) CreateDir(c *gin.Context) {
	var data *models.DirData
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		common.Fail(c, 50000, "parse param error")
		return
	}
	r, err := h.s.CreateDir(c.Request.Context(), data)
	if err != nil {
		common.Fail(c, 50000, "Create error:"+err.Error())
		return
	}
	common.Success().Append("data", r).Append("isSuccess", true).End(c)
}

// update directory info
func (h *DebugHandler) UpdateDir(c *gin.Context) {
	var data *models.DirDataInfo
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		common.Fail(c, 50000, "parse param error")
		return
	}
	dir, err := h.s.UpdateDir(c.Request.Context(), data)
	if err != nil {
		common.Fail(c, 50000, "handler error")
		return
	}
	common.Success().Append("data", dir).Append("isSuccess", true).End(c)
}
