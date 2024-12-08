package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/google/uuid"

	"shirinec.com/config"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

type MediaService interface {
	Create(ctx context.Context, savePath string, userID uuid.UUID, input *dto.MediaUploadQuery) (*dto.MediaUploadResponse, error)
	GetMedia(ctx context.Context, mediaName string, userID uuid.UUID) (string, error)
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

func (s *mediaService) Create(ctx context.Context, fileName string, userID uuid.UUID, input *dto.MediaUploadQuery) (*dto.MediaUploadResponse, error) {
	var media models.Media
	media.UserID = userID
	media.FilePath = fileName
	currentTime := time.Now().UTC().Truncate(time.Second)
	media.CreationDate = currentTime
	media.UpdateDate = currentTime
	media.Access = &input.Access
	media.FinancialGroupID = input.FinancialGroupID
	url := fmt.Sprintf("/file/media-%s", uuid.New().String())
	media.Url = url

	if err := s.mediaRepo.Create(ctx, &media); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}

		utils.Logger.Errorf("mediaService.Create - Calling mediaRepo.CreateFor: %s", err.Error())
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

func (s *mediaService) GetMedia(ctx context.Context, mediaName string, userID uuid.UUID) (string, error) {
	url := fmt.Sprintf("/file/%s", mediaName)
	media, err := s.mediaRepo.GetByMediaName(ctx, url, userID)
	if err != nil {
        utils.Logger.Infof("mediaService.GetMedia - calling mediaRepo.GetByMediaName: %s", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return "", &server_errors.ItemNotFound
		}
		if sErr := server_errors.AsPgError(err); sErr != nil {
			return "", sErr
		}

		utils.Logger.Errorf("mediaService.GetMedia - Calling mediaRepo.GetByMediaName: %s", err.Error())
		return "", &server_errors.InternalError
	}

	mediaPath := path.Join(config.AppConfig.UploadFolder, media.FilePath)
	return mediaPath, nil
}
