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
	"shirinec.com/internal/utils"
)

type CategoryService interface {
	Create(category *models.Category) error
	ListCategories(userID uuid.UUID, page int, size int) (*dto.CategoriesListResponse, error)
	GetByID(userID uuid.UUID, id int) (*models.Category, error)
    Delete(userID uuid.UUID, id int) error
    Update(userID *uuid.UUID, id int, category *dto.CategoryUpdateRequest) (*models.Category ,error)
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) Create(category *models.Category) error {
	return s.categoryRepo.Create(context.Background(), category)
}

func (s *categoryService) ListCategories(userID uuid.UUID, page int, size int) (*dto.CategoriesListResponse, error) {
	var response dto.CategoriesListResponse

	limit := size
	offset := page * size
	categories, totalCount, err := s.categoryRepo.List(context.Background(), limit, offset, userID)
	totalPages := int(math.Ceil(float64(totalCount) / float64(size)))
	remainingPages := int(math.Max(float64(totalPages-page-1), 0))

	response.Pagination.PageNumber = page
	response.Pagination.TotalRecord = totalCount
	response.Pagination.PageSize = size
	response.Pagination.RemainingPages = remainingPages

	if err != nil {
		log.Printf("[Error] - CategoryService.List - Getting categories from repository: %+v\n", err)

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

func (s *categoryService) GetByID(userID uuid.UUID, id int) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(context.Background(), id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}

        utils.Logger.Errorf("Calling categoryService.GetByID: %s", err.Error())
		return category, &server_errors.InternalError
	}
	return category, nil
}

func (s *categoryService) Delete(userID uuid.UUID, id int) error {
	err := s.categoryRepo.Delete(context.Background(), id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &server_errors.ItemNotFound
		}
        utils.Logger.Errorf("Calling categoryService.Delete: %s", err.Error())
		return &server_errors.InternalError
	}
	return nil
}

func (s *categoryService) Update(userID *uuid.UUID, id int, categoryDTO *dto.CategoryUpdateRequest) (*models.Category ,error) {
    var category models.Category
    category.ID = id
    category.UserID = *userID
    category.Color = categoryDTO.Color
    category.Name = categoryDTO.Name

	err := s.categoryRepo.Update(context.Background(), &category)
	if err != nil {
		log.Printf("[Error] - CategoryService.Update - Getting category from repository: %+v\n", err)
        var sError *server_errors.SError
        if errors.As(err, &sError) {
            return nil, err
        }
		if errors.Is(err, sql.ErrNoRows) {
            log.Printf("Item was not found!")
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
