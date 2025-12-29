package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/myapp/internal/common"
	"github.com/myapp/internal/modules/envzip/models"
	"github.com/myapp/internal/modules/envzip/service"
)

type EnvZipHandler struct {
	s service.EnvZipService
}

func NewEnvZipHandler(s service.EnvZipService) *EnvZipHandler {
	return &EnvZipHandler{s}
}

// add envdata
func (h *EnvZipHandler) CreateEnv(c *gin.Context) {
	var data *models.EnvDataDTO
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		common.Fail(c, 50000, "parse param error")
		return
	}
	_, err = h.s.Create(c.Request.Context(), data)
	if err != nil {
		common.Fail(c, 50000, "handler error")
		return
	}
	common.Success().Append("isSuccess", true).End(c)
}

// remover envzip
func (h *EnvZipHandler) DeleteEnv(c *gin.Context) {
	tmpid := c.Param("id")
	id, err := strconv.ParseUint(tmpid, 10, 64)
	if err != nil {
		common.Fail(c, 50000, "id can't covert to uint, please check param!")
		return
	}
	err = h.s.Delete(c.Request.Context(), uint(id))
	if err != nil {
		common.Fail(c, 50000, "Delete error:"+err.Error())
		return
	}
	common.Success().Append("isSuccess", true).End(c)
}

// update envdata
func (h *EnvZipHandler) UpdateEnv(c *gin.Context) {
	var data *models.EnvDataDTO
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		common.Fail(c, 50000, "parse param error")
		return
	}
	_, err = h.s.Update(c.Request.Context(), data)
	if err != nil {
		common.Fail(c, 50000, "handler error")
		return
	}
	common.Success().Append("isSuccess", true).End(c)
}

// find envdata by parentid
func (h *EnvZipHandler) FindByPId(c *gin.Context) {
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
	data, err := h.s.FindByPId(c.Request.Context(), ppId)
	if err != nil {
		common.Fail(c, 50000, "Query error: "+err.Error())
		return
	}
	common.Success().Append("data", data).End(c)
}

// find all env dirdata
func (h *EnvZipHandler) FindDir(c *gin.Context) {
	data, err := h.s.FindDir(c.Request.Context())
	if err != nil {
		common.Fail(c, 50000, "Query error: "+err.Error())
		return
	}
	common.Success().Append("data", data).End(c)
}
