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

type CategoryService interface {
	Create(category *models.Category) error
	ListCategories(userID uuid.UUID, page int, size int) (*dto.CategoriesListResponse, error)
	GetByID(userID uuid.UUID, id int) (*models.Category, error)
	Delete(userID uuid.UUID, id int) error
	Update(userID *uuid.UUID, id int, category *dto.CategoryUpdateRequest) (*models.Category, error)
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
		utils.Logger.Errorf("CategoryService.List - Getting categories from repository: %s", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}

		if pgErr := server_errors.AsPgError(err); pgErr != nil {
            return nil, pgErr
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

func (s *categoryService) Update(userID *uuid.UUID, id int, categoryDTO *dto.CategoryUpdateRequest) (*models.Category, error) {
	var category models.Category
	category.ID = id
	category.UserID = *userID
	category.Color = categoryDTO.Color
	category.Name = categoryDTO.Name
	category.IconID = categoryDTO.IconID

	err := s.categoryRepo.Update(context.Background(), &category)
	if err != nil {
		var sError *server_errors.SError
		if errors.As(err, &sError) {
			return nil, sError
		}
		if errors.Is(err, sql.ErrNoRows) {
			return &category, &server_errors.ItemNotFound
		}

        if pgErr := server_errors.AsPgError(err); pgErr != nil {
            return nil, pgErr
		}

        utils.Logger.Errorf("CategoryService.Update - Getting category from repository: %s", err.Error())
		return &category, &server_errors.InternalError
	}
	return &category, nil
}
