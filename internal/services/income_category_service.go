package services

import (
	"context"
	"log"

	"github.com/google/uuid"
	"shirinec.com/internal/models"
	"shirinec.com/internal/repositories"
)

type IncomeCategoryService interface {
	ListCategories(userID uuid.UUID, limit int, offset int) (*[]models.IncomeCategory, error)
    GetByID(userID uuid.UUID, ID int) (*models.IncomeCategory, error)
}

type incomeCategoryService struct {
	incomeCategoryRepo repositories.IncomeCategoryRepository
}

func NewIncomeService(incomeCategoryRepo repositories.IncomeCategoryRepository) IncomeCategoryService {
	return &incomeCategoryService{incomeCategoryRepo: incomeCategoryRepo}
}

func (s *incomeCategoryService) ListCategories(userID uuid.UUID, limit int, offset int) (*[]models.IncomeCategory, error) {
    categories, err := s.incomeCategoryRepo.List(context.Background(), limit, offset, userID)
    if err != nil {
        log.Printf("%+v", err)
    }
	return categories, err
}

func (s *incomeCategoryService) GetByID(userID uuid.UUID, id int) (*models.IncomeCategory, error) {
    category, err := s.incomeCategoryRepo.GetByID(context.Background(), id, userID)
    if err != nil {
        log.Printf("[Error] - IncomeCategoryService.GetByID - while getting category from repository: %+v\n", err)
    }
    return category, err
}
