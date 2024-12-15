package services

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

type AccountService interface {
	Create(ctx context.Context, account *dto.AccountCreateRequest, userID uuid.UUID) (*models.Account, error)
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
	var account models.Account
	account.UserID = userID
	account.CategoryID = &input.CategoryID
	account.Name = &input.Name
	account.Balance = &input.Balance
	account.Type = &input.Type

	err := s.accountRepo.Create(ctx, &account)
	if err != nil {
		pgErr := server_errors.AsPgError(err)
		if pgErr != nil {
			return nil, pgErr
		}
		utils.Logger.Errorf("accountService.Create - Calling accountRepo.Create: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	return &account, nil
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}

		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return nil, pgErr
		}

		utils.Logger.Errorf("accountService.List - Calling accountRepo.List: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	response.Accounts = accounts
	return &response, nil
}

func (s *accountService) GetByID(ctx context.Context, id int, userID uuid.UUID) (*dto.AccountJoinedResponse, error) {
	account, err := s.accountRepo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}
		utils.Logger.Errorf("accountService.GetByID - Calling accountRepository.GetByID: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	return account, nil
}

func (s *accountService) Delete(ctx context.Context, id int, userID uuid.UUID) error {
	err := s.accountRepo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &server_errors.ItemNotFound
		}
		utils.Logger.Errorf("accountService.Delete - Calling accountRepo.Delete: %s", err.Error())
		return &server_errors.InternalError
	}

	return nil
}

func (s *accountService) Update(ctx context.Context, input *dto.AccountUpdateRequest, id int, userID uuid.UUID) (*dto.AccountJoinedResponse, error) {
	var account models.Account
	account.ID = id
	account.UserID = userID
	account.Name = input.Name
	account.Balance = input.Balance
	account.CategoryID = input.CategoryID

	accountJoined, err := s.accountRepo.Update(ctx, &account)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}

		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return nil, pgErr
		}

		utils.Logger.Errorf("accountService.Update - Calling accountRepo.Update: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	return accountJoined, nil
}
