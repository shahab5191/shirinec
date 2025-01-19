package handler

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shirinec.com/config"
	"shirinec.com/src/internal/dto"
	"shirinec.com/src/internal/enums"
	"shirinec.com/src/internal/errors"
	"shirinec.com/src/internal/services"
	"shirinec.com/src/internal/utils"
)

type MediaHandler interface {
	Upload(c *gin.Context)
	GetMedia(c *gin.Context)
	UpdateMedia(c *gin.Context)
}

type mediaHandler struct {
	mediaService services.MediaService
}

func NewMediaHandler(mediaService services.MediaService) MediaHandler {
	return &mediaHandler{mediaService: mediaService}
}

func (h *mediaHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(server_errors.FileRequired.Unwrap())
		return
	}

	var input dto.MediaUploadQuery
	input.Access = enums.Owner
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("mediaHandler.Upload - Parsing uuid from user_id string: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		c.JSON(server_errors.InvalidFileFormat.Unwrap())
		return
	}

	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := fmt.Sprintf("%s/%s", config.AppConfig.UploadFolder, fileName)
	if err = c.SaveUploadedFile(file, savePath); err != nil {
		utils.Logger.Errorf("mediaHandler.Upload - Saving file: %s", err.Error())
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	media, err := h.mediaService.Create(context.Background(), fileName, userID, &input)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, media)
}

func (h *mediaHandler) GetMedia(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		utils.Logger.Errorf("mediaHandler.GetMedia - Parsing uuid from user_id string: %s", err.Error())
		return
	}

	mediaName := c.Param("fileName")
	mediaPath, err := h.mediaService.GetMedia(context.Background(), mediaName, userID)
	if err != nil {
		c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.File(mediaPath)
}

func (h *mediaHandler) UpdateMedia(c *gin.Context) {

}
