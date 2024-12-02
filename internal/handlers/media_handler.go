package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shirinec.com/config"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/services"
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
		log.Printf("[Error] - mediaHandler.Upload - Calling c.FormFile: %+v\n", err)
		c.JSON(server_errors.FileRequired.Unwrap())
		return
	}

	var input dto.MediaUploadRequest
	if err = c.ShouldBindQuery(&input); err != nil {
        log.Printf("[Error] - mediaHandler.Upload - Binding input to dto.MediaUploadRequest: %+v\n", err)
		c.JSON(server_errors.InvalidInput.Unwrap())
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		log.Printf("[Error] - mediaHandler.Upload - Parsing uuid from user_id string: %+v\n", err)
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
		log.Printf("[Error] - mediaHandler.Upload - Saving file: %+v\n", err)
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	media, err := h.mediaService.Create(context.Background(), &input, savePath, userID)
	if err != nil {
		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			c.JSON(sErr.Unwrap())
			return
		}
		c.JSON(server_errors.InternalError.Unwrap())
		return
	}

	c.JSON(http.StatusOK, media)
}
