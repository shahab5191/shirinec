package services

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
	"shirinec.com/internal/repositories"
)

type ItemService interface {
	Create(ctx context.Context, item *dto.ItemCreateRequest, userID uuid.UUID) (*models.Item, error)
    List(ctx context.Context, page, size int, userID uuid.UUID) (*dto.ItemsListResponse, error)
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
        log.Printf("[Error] - itemService.Create - Calling itemRepo.Create: %+v\n", err)
        return nil, &server_errors.InternalError
    }

    return &item, nil
}

func (s *itemService) List(ctx context.Context, page, size int, userID uuid.UUID) (*dto.ItemsListResponse, error){
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
        log.Printf("[Error] - itemService.List - Calling itemRepo.List: %+v\n", err)

        if errors.Is(err, sql.ErrNoRows) {
            return nil, &server_errors.ItemNotFound
        }

        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            log.Printf("Error is of type pgconn.PgError: %+v\n", pgErr)
        }

        return nil, &server_errors.InternalError
    }

    response.Items = items
    return &response, nil
}
