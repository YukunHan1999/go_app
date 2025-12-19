package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/myapp/internal/common"
	"github.com/myapp/internal/modules/codetemplate/models"
	"github.com/myapp/internal/modules/codetemplate/service"
)

type CodeTemplateHandler struct {
	svc service.CodeTemplateService
}

func NewCodeTemplateHandler(svc service.CodeTemplateService) *CodeTemplateHandler {
	return &CodeTemplateHandler{svc}
}

// query total count
func (h *CodeTemplateHandler) Count(c *gin.Context) {
	count, err := h.svc.Count(c.Request.Context())
	if err != nil {
		common.Fail(c, 500, "query count throw error"+err.Error())
		return
	}
	res := make(map[string]any)
	res["total"] = count
	common.Success().Append("total", count).End(c)
}

func (h *CodeTemplateHandler) SaveOrUpdate(c *gin.Context) {
	var req struct {
		Id   string `json:"id"`
		Lang string `json:"lang"`
		Name string `json:"name"`
		Code string `json:"code"`
		Tips string `json:"tips"`
	}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		common.Fail(c, 50000, "parse param failed!")
		return
	}

	var (
		codeTmp *models.Codetemplate
		err     error
	)
	if req.Id != "" {
		codeTmp, err = h.svc.Update(c.Request.Context(), req.Id, req.Lang, req.Name, req.Code, req.Tips)
		if err != nil {
			common.Fail(c, 50000, "update error:"+err.Error())
			return
		}
	} else {
		codeTmp, err = h.svc.Create(c.Request.Context(), req.Lang, req.Name, req.Code, req.Tips)
		if err != nil {
			common.Fail(c, 50000, "create error:"+err.Error())
			return
		}
	}
	common.Success().Append("data", codeTmp).Append("isSuccess", true).End(c)
}

func (h *CodeTemplateHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		common.Fail(c, 50000, "Delete error"+err.Error())
		return
	}
	common.Success().Append("isSuccess", true).End(c)
}

func (h *CodeTemplateHandler) FindAll(c *gin.Context) {
	pageSize := c.Query("pageSize")
	pageNo := c.Query("pageNo")

	psize, err := strconv.Atoi(pageSize)
	if err != nil {
		common.Fail(c, 50000, "parse param pageSize failed: "+err.Error())
		return
	}
	pno, err := strconv.Atoi(pageNo)
	if err != nil {
		common.Fail(c, 50000, "parse param pageNo failed: "+err.Error())
		return
	}

	codeTmps, err := h.svc.FindAll(c.Request.Context(), psize, pno)
	if err != nil {
		common.Fail(c, 50000, "search data throw error:"+err.Error())
		return
	}

	common.Success().Append("data", codeTmps).End(c)
}
