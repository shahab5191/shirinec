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

type ExpenseCategoryService interface {
	Create(category *models.ExpenseCategory) error
	ListCategories(userID uuid.UUID, page int, size int) (*dto.ListExpenseCategoriesResponse, error)
	GetByID(userID uuid.UUID, id int) (*models.ExpenseCategory, error)
    Delete(userID uuid.UUID, id int) error
    Update(userID *uuid.UUID, id int, category *dto.UpdateIncomeCategoryRequest) (*models.ExpenseCategory ,error)
}

type expenseCategoryService struct {
	expenseCategoryRepo repositories.ExpenseCategoryRepository
}

func NewExpenseService(expenseCategoryRepo repositories.ExpenseCategoryRepository) ExpenseCategoryService {
	return &expenseCategoryService{expenseCategoryRepo: expenseCategoryRepo}
}

func (s *expenseCategoryService) Create(category *models.ExpenseCategory) error {
	return s.expenseCategoryRepo.Create(context.Background(), category)
}

func (s *expenseCategoryService) ListCategories(userID uuid.UUID, page int, size int) (*dto.ListExpenseCategoriesResponse, error) {
	var response dto.ListExpenseCategoriesResponse

	limit := size
	offset := page * size
	categories, totalCount, err := s.expenseCategoryRepo.List(context.Background(), limit, offset, userID)
	totalPages := int(math.Ceil(float64(totalCount) / float64(size)))
	remainingPages := int(math.Max(float64(totalPages-page-1), 0))

	response.Pagination.PageNumber = page
	response.Pagination.TotalRecord = totalCount
	response.Pagination.PageSize = size
	response.Pagination.RemainingPages = remainingPages

	if err != nil {
		log.Printf("[Error] - ExpenseCategoryService.List - Getting categories from repository: %+v\n", err)

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

func (s *expenseCategoryService) GetByID(userID uuid.UUID, id int) (*models.ExpenseCategory, error) {
	category, err := s.expenseCategoryRepo.GetByID(context.Background(), id, userID)
	if err != nil {
		log.Printf("[Error] - ExpenseCategoryService.GetByID - Getting category from repository: %+v\n", err)
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

func (s *expenseCategoryService) Delete(userID uuid.UUID, id int) error {
	err := s.expenseCategoryRepo.Delete(context.Background(), id, userID)
	if err != nil {
		log.Printf("[Error] - ExpenseCategoryService.GetByID - Getting category from repository: %+v\n", err)
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

func (s *expenseCategoryService) Update(userID *uuid.UUID, id int, categoryDTO *dto.UpdateIncomeCategoryRequest) (*models.ExpenseCategory ,error) {
    var category models.ExpenseCategory
    category.ID = id
    category.UserID = *userID
    category.Color = categoryDTO.Color
    category.Name = categoryDTO.Name

	err := s.expenseCategoryRepo.Update(context.Background(), &category)
	if err != nil {
		log.Printf("[Error] - ExpenseCategoryService.GetByID - Getting category from repository: %+v\n", err)
		if errors.Is(err, sql.ErrNoRows) {
			return &category, &server_errors.ItemNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Error is of type pgconn.PgError: %+v\n", pgErr)
		}
		return &category, &server_errors.InternalError
	}
	return &category, nil
}
