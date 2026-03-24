package handler

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/myapp/internal/common"
	"github.com/myapp/internal/modules/attachment/models"
	"github.com/myapp/internal/modules/attachment/service"
)

type AttachmentHandler struct {
	S service.AttachmentService
}

func NewAttachmentHandler(S service.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{S}
}

func (h AttachmentHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.Fail(c, 50000, "No file uploaded")
		return
	}

	uid := uuid.New()
	randomId := uid.String()
	filetype := filepath.Ext(file.Filename)
	newFilename := randomId + filetype

	// save file to local folder ./uploads/
	savePath := filepath.Join("uploads", newFilename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		common.Fail(c, 50000, "Save file failed")
		return
	}
	// save
	att := &models.Attachment{
		Name: file.Filename,
		Type: filetype,
		Url:  newFilename,
	}
	att, err = h.S.Create(c.Request.Context(), att)

	if err != nil {
		common.Fail(c, 50000, "Create file info failed")
		return
	}
	common.Success().Append("data", att).End(c)
}

func (h AttachmentHandler) PreviewFile(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join("uploads", filename)
	c.File(filepath)
}

// clear garabage attachment
func (h AttachmentHandler) ClearGarbageAtt(c *gin.Context) {
	var r []uint
	err := c.ShouldBindBodyWithJSON(&r)
	if err != nil {
		common.Fail(c, 50000, "Parse param error!")
		return
	}
	err = h.S.BatchDeleteByIds(c.Request.Context(), r, nil)
	if err != nil {
		common.Fail(c, 50000, "delete attachment error! ")
		return
	}
	common.Success().Append("isSuccess", true).End(c)
}
