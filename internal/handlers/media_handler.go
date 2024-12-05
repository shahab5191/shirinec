package handler

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shirinec.com/config"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/services"
	"shirinec.com/internal/utils"
)

type MediaHandler interface {
	Upload(c *gin.Context)
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
		utils.Logger.Errorf("Calling c.FormFile: %s", err.Error())
		c.JSON(server_errors.FileRequired.Unwrap())
		return
	}

	var input dto.MediaUploadRequest
	if err = c.ShouldBindQuery(&input); err != nil {
        if errList := server_errors.AsValidatorError(err); errList != nil {
			c.JSON(server_errors.ValidationErrorBuilder(errList).Unwrap())
			return
		}
        utils.Logger.Warnf("Undefined error while binding input to dto.MediaUploadRequest: %s", err.Error())
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

	media, err := h.mediaService.Create(context.Background(), &input, savePath, userID)
	if err != nil {
        c.JSON(err.(*server_errors.SError).Unwrap())
		return
	}

	c.JSON(http.StatusOK, media)
}
