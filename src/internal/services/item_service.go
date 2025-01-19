package services

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/google/uuid"
	"shirinec.com/src/internal/dto"
	"shirinec.com/src/internal/errors"
	"shirinec.com/src/internal/models"
	"shirinec.com/src/internal/repositories"
	"shirinec.com/src/internal/utils"
)

type ItemService interface {
	Create(ctx context.Context, item *dto.ItemCreateRequest, userID uuid.UUID) (*models.Item, error)
	List(ctx context.Context, page, size int, userID uuid.UUID) (*dto.ItemsListResponse, error)
	GetByID(ctx context.Context, id int, userID uuid.UUID) (*dto.ItemJoinedResponse, error)
	Update(ctx context.Context, input *dto.ItemUpdateRequest, id int, userID uuid.UUID) (*dto.ItemJoinedResponse, error)
	Delete(ctx context.Context, id int, userID uuid.UUID) error
}

type itemService struct {
	itemRepo repositories.ItemRepository
}

func NewItemService(itemRepo *repositories.ItemRepository) ItemService {
	return &itemService{itemRepo: *itemRepo}
}

func (s *itemService) Create(ctx context.Context, input *dto.ItemCreateRequest, userID uuid.UUID) (*models.Item, error) {
	var item models.Item
	item.UserID = userID
	item.CategoryID = &input.CategoryID
	item.Name = &input.Name
	item.ImageID = input.ImageID

	err := s.itemRepo.Create(ctx, &item)
	if err != nil {
		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return nil, pgErr
		}
		utils.Logger.Errorf("itemService.Create - Calling itemRepo.Create: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	return &item, nil
}

func (s *itemService) List(ctx context.Context, page, size int, userID uuid.UUID) (*dto.ItemsListResponse, error) {
	var response dto.ItemsListResponse

	limit := size
	offset := page * size
	items, totalCount, err := s.itemRepo.List(context.Background(), limit, offset, userID)
	totalPages := int(math.Ceil(float64(totalCount) / float64(size)))
	remainingPages := int(math.Max(float64(totalPages-page-1), 0))

	response.Pagination.PageNumber = page
	response.Pagination.TotalRecord = totalCount
	response.Pagination.PageSize = size
	response.Pagination.RemainingPages = remainingPages

	if err != nil {
		utils.Logger.Errorf("itemService.List - Calling itemRepo.List: %s", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}

		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			utils.Logger.Errorf("Error is of type pgconn.PgError: %+v\n", pgErr.Error())
		}

		return nil, &server_errors.InternalError
	}

	response.Items = items
	return &response, nil
}

func (s *itemService) GetByID(ctx context.Context, id int, userID uuid.UUID) (*dto.ItemJoinedResponse, error) {
	item, err := s.itemRepo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}
		utils.Logger.Errorf("itemService.GetByID - Calling itemRepository.GetByID: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	return item, nil
}

func (s *itemService) Delete(ctx context.Context, id int, userID uuid.UUID) error {
	err := s.itemRepo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &server_errors.ItemNotFound
		}
		utils.Logger.Errorf("itemService.Delete - Calling itemRepo.Delete: %s", err.Error())
		return &server_errors.InternalError
	}

	return nil
}

func (s *itemService) Update(ctx context.Context, input *dto.ItemUpdateRequest, id int, userID uuid.UUID) (*dto.ItemJoinedResponse, error) {
	var item models.Item
	item.ID = id
	item.UserID = userID
	item.Name = input.Name
	item.ImageID = input.ImageID
	item.CategoryID = input.CategoryID

	itemJoined, err := s.itemRepo.Update(ctx, &item)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}
		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return nil, pgErr
		}
		utils.Logger.Errorf("itemService.Update - Calling itemRepo.Update: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	return itemJoined, nil
}
