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

type IncomeCategoryService interface {
	Create(category *models.IncomeCategory) error
	ListCategories(userID uuid.UUID, page int, size int) (*dto.ListCategoriesResponse, error)
	GetByID(userID uuid.UUID, id int) (*models.IncomeCategory, error)
    Delete(userID uuid.UUID, id int) error
}

type incomeCategoryService struct {
	incomeCategoryRepo repositories.IncomeCategoryRepository
}

func NewIncomeService(incomeCategoryRepo repositories.IncomeCategoryRepository) IncomeCategoryService {
	return &incomeCategoryService{incomeCategoryRepo: incomeCategoryRepo}
}

func (s *incomeCategoryService) Create(category *models.IncomeCategory) error {
	return s.incomeCategoryRepo.Create(context.Background(), category)
}

func (s *incomeCategoryService) ListCategories(userID uuid.UUID, page int, size int) (*dto.ListCategoriesResponse, error) {
	var response dto.ListCategoriesResponse

	limit := size
	offset := page * size
	categories, totalCount, err := s.incomeCategoryRepo.List(context.Background(), limit, offset, userID)
	totalPages := int(math.Ceil(float64(totalCount) / float64(size)))
	remainingPages := int(math.Max(float64(totalPages-page-1), 0))

	response.Pagination.PageNumber = page
	response.Pagination.TotalRecord = totalCount
	response.Pagination.PageSize = size
	response.Pagination.RemainingPages = remainingPages

	if err != nil {
		log.Printf("[Error] - IncomeCategoryService.List - Getting categories from repository: %+v\n", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Error is of type pgconn.PgError: %+v\n", pgErr)
		}

		return &response, &server_errors.InternalError
	}
	response.Categories = *categories
	return &response, err
}

func (s *incomeCategoryService) GetByID(userID uuid.UUID, id int) (*models.IncomeCategory, error) {
	category, err := s.incomeCategoryRepo.GetByID(context.Background(), id, userID)
	if err != nil {
		log.Printf("[Error] - IncomeCategoryService.GetByID - Getting category from repository: %+v\n", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Error is of type pgconn.PgError: %+v\n", pgErr)
		}
		return category, &server_errors.InternalError
	}
	return category, nil
}

func (s *incomeCategoryService) Delete(userID uuid.UUID, id int) error {
	err := s.incomeCategoryRepo.Delete(context.Background(), id, userID)
	if err != nil {
		log.Printf("[Error] - IncomeCategoryService.GetByID - Getting category from repository: %+v\n", err)
		if errors.Is(err, sql.ErrNoRows) {
			return &server_errors.ItemNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Error is of type pgconn.PgError: %+v\n", pgErr)
		}
		return &server_errors.InternalError
	}
	return nil
}
