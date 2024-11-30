package services

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"shirinec.com/internal/db"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
	"shirinec.com/internal/repositories"
)

type AccountService interface {
	Create(ctx context.Context, item *dto.AccountCreateRequest, userID uuid.UUID) (*models.Account, error)
	List(ctx context.Context, page, size int, userID uuid.UUID) (*dto.AccountListResponse, error)
	GetByID(ctx context.Context, id int, userID uuid.UUID) (*dto.AccountJoinedResponse, error)
	Update(ctx context.Context, input *dto.AccountUpdateRequest, id int, userID uuid.UUID) (*dto.AccountJoinedResponse, error)
	Delete(ctx context.Context, id int, userID uuid.UUID) error
}

type accountService struct {
	accountRepo repositories.AccountRepository
}

func NewAccountService(accountRepo *repositories.AccountRepository) AccountService {
	return &accountService{accountRepo: *accountRepo}
}

func (s *accountService) Create(ctx context.Context, input *dto.AccountCreateRequest, userID uuid.UUID) (*models.Account, error) {
	var item models.Account
	item.UserID = userID
	item.CategoryID = &input.CategoryID
	item.Name = &input.Name
	item.Balance = &input.Balance

	err := s.accountRepo.Create(ctx, &item)
	if err != nil {
		log.Printf("[Error] - accountService.Create - Calling accountRepo.Create: %+v\n", err)
		return nil, &server_errors.InternalError
	}

	return &item, nil
}

func (s *accountService) List(ctx context.Context, page, size int, userID uuid.UUID) (*dto.AccountListResponse, error) {
	var response dto.AccountListResponse

	limit := size
	offset := page * size
	accounts, totalCount, err := s.accountRepo.List(context.Background(), limit, offset, userID)
	totalPages := int(math.Ceil(float64(totalCount) / float64(size)))
	remainingPages := int(math.Max(float64(totalPages-page-1), 0))

	response.Pagination.PageNumber = page
	response.Pagination.TotalRecord = totalCount
	response.Pagination.PageSize = size
	response.Pagination.RemainingPages = remainingPages

	if err != nil {
		log.Printf("[Error] - accountService.List - Calling accountRepo.List: %+v\n", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Error is of type pgconn.PgError: %+v\n", pgErr)
		}

		return nil, &server_errors.InternalError
	}

	response.Accounts = accounts
	return &response, nil
}

func (s *accountService) GetByID(ctx context.Context, id int, userID uuid.UUID) (*dto.AccountJoinedResponse, error) {
	item, err := s.accountRepo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}
		log.Printf("[Error] - accountService.GetByID - Calling accountRepository.GetByID: %+v\n", err)
		return nil, &server_errors.InternalError
	}

	return item, nil
}

func (s *accountService) Delete(ctx context.Context, id int, userID uuid.UUID) error {
	err := s.accountRepo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &server_errors.ItemNotFound
		}
		log.Printf("[Error] - accountService.Delete - Calling accountRepo.Delete: %+v\n", err)
		return &server_errors.InternalError
	}

	return nil
}

func (s *accountService) Update(ctx context.Context, input *dto.AccountUpdateRequest, id int, userID uuid.UUID) (*dto.AccountJoinedResponse, error) {
	var item models.Account
	item.ID = id
	item.UserID = userID
	item.Name = input.Name
	item.Balance = input.Balance
	item.CategoryID = input.CategoryID

	itemJoined, err := s.accountRepo.Update(ctx, &item)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == db.ForeignKeyViolation {
				return nil, &server_errors.InvalidInput
			}
		}
		log.Printf("[Error] - accountService.Update - Calling accountRepo.Update: %+v\n", err)
		return nil, &server_errors.InternalError
	}

	return itemJoined, nil
}
