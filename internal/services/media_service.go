package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"shirinec.com/internal/dto"
	"shirinec.com/internal/enums"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

type MediaService interface {
	Create(ctx context.Context, input *dto.MediaUploadRequest, savePath string, userID uuid.UUID) (*dto.MediaUploadResponse, error)
}

type mediaService struct {
	mediaRepo    repositories.MediaRepository
	itemRepo     repositories.ItemRepository
	categoryRepo repositories.CategoryRepository
}

func NewMediaService(mediaRepo repositories.MediaRepository, itemRepo repositories.ItemRepository, categoryRepo repositories.CategoryRepository) MediaService {
	return &mediaService{
		mediaRepo:    mediaRepo,
		categoryRepo: categoryRepo,
		itemRepo:     itemRepo,
	}
}

func (s *mediaService) Create(ctx context.Context, input *dto.MediaUploadRequest, savePath string, userID uuid.UUID) (*dto.MediaUploadResponse, error) {
	var media models.Media
	media.UserID = userID
	media.FilePath = savePath
	currentTime := time.Now().UTC().Truncate(time.Second)
	media.CreationDate = currentTime
	media.UpdateDate = currentTime
	url := fmt.Sprintf("/file/%s-%s", input.BindsTo, uuid.New().String())
	media.Url = url

	var err error
	switch input.BindsTo {
	case enums.BindToItem:
		err = s.mediaRepo.CreateForEntity(ctx, "items", "image_id", &media, input.BindID)
	case enums.BindToCategory:
		err = s.mediaRepo.CreateForEntity(ctx, "categories", "icon_id", &media, input.BindID)
	case enums.BindToProfile:
		err = s.mediaRepo.CreateForProfile(ctx, &media)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}

        utils.Logger.Errorf("mediaService.Create - Calling mediaRepo.CreateFor%s: %s", input.BindsTo, err.Error())
		return nil, &server_errors.InternalError
	}

	var mediaResponse *dto.MediaUploadResponse = &dto.MediaUploadResponse{
		ID:           media.ID,
		URL:          media.Url,
		Metadata:     media.Metadata,
		UpdateDate:   media.UpdateDate,
		CreationDate: media.CreationDate,
	}

	return mediaResponse, nil
}
