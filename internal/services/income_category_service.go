package services

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	server_errors "shirinec.com/internal/errors"
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
        log.Printf("[Error] - IncomeCategoryService.List - Getting categories from repository: %+v\n", err)

        if errors.Is(err, sql.ErrNoRows){
            return nil, &server_errors.ItemNotFound
        }

        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr){
            log.Printf("Error is of type pgconn.PgError: %+v\n", pgErr)
        }

        return categories, &server_errors.InternalError
    }
	return categories, err
}

func (s *incomeCategoryService) GetByID(userID uuid.UUID, id int) (*models.IncomeCategory, error) {
    category, err := s.incomeCategoryRepo.GetByID(context.Background(), id, userID)
    if err != nil {
        log.Printf("[Error] - IncomeCategoryService.GetByID - Getting category from repository: %+v\n", err)
        if errors.Is(err, sql.ErrNoRows){
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
