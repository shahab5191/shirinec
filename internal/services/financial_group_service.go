package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/enums"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

type FinancialGroupService interface {
	Create(ctx context.Context, input *dto.FinancialGroupCreateRequest, userID uuid.UUID) (*models.FinancialGroups, error)
	AddUserToGroup(ctx context.Context, financialGroupID int, newUserID uuid.UUID, userID uuid.UUID) error
	GetByID(ctx context.Context, id int, userID uuid.UUID) (*dto.FinancialGroup, error)
	List(ctx context.Context, input dto.ListRequest, userID uuid.UUID) (*dto.FinancialGroupListResponse, error)
}

type financialGroupService struct {
	financialGroupRepo repositories.FinancialGroupRepository
}

func NewFinancialGroupService(financialGroupRepo *repositories.FinancialGroupRepository) FinancialGroupService {
	return &financialGroupService{
		financialGroupRepo: *financialGroupRepo,
	}
}

func (s *financialGroupService) Create(ctx context.Context, input *dto.FinancialGroupCreateRequest, userID uuid.UUID) (*models.FinancialGroups, error) {
	var financialGroup models.FinancialGroups
	financialGroup.UserID = userID
	financialGroup.Name = input.Name
	financialGroup.ImageID = input.ImageID
	currentTime := time.Now().UTC().Truncate(time.Second)
	financialGroup.CreationDate = currentTime
	financialGroup.UpdateDate = currentTime

	if err := s.financialGroupRepo.Create(ctx, &financialGroup); err != nil {
		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return nil, pgErr
		}

		utils.Logger.Errorf("financialGroupService.Create - Calling financialGroupRepo.Create: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	return &financialGroup, nil
}

func (s *financialGroupService) AddUserToGroup(ctx context.Context, financialGroupID int, newUserID uuid.UUID, userID uuid.UUID) error {
	// First checking to see if user is owner of selected group by getting group by id
	_, err := s.financialGroupRepo.GetOwnedGroupByID(ctx, financialGroupID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &server_errors.ItemNotFound
		}
		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return pgErr
		}
		utils.Logger.Errorf("financialGroupService.AddUserToGroup - Calling financialGroupRep.GetByID(%d): %s", financialGroupID, err)
		return &server_errors.InternalError
	}

	if err = s.financialGroupRepo.AddUserToGroup(ctx, financialGroupID, newUserID); err != nil {
		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return pgErr
		}

		utils.Logger.Errorf("financialGroupService.AddUserToGroup - Calling financialGroupRepo.AddUserToGroup: (id: %d) %s", financialGroupID, err.Error())
		return &server_errors.InternalError
	}

	return nil
}

func (s *financialGroupService) GetByID(ctx context.Context, id int, userID uuid.UUID) (*dto.FinancialGroup, error) {
	financialGroup, err := s.financialGroupRepo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}

		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return nil, pgErr
		}

		utils.Logger.Errorf("financialGroupService.GetByID - Calling GetByID: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	if userID == financialGroup.OwnerID {
		financialGroup.UserRole = enums.FinancialGroupOwner
	} else {
		financialGroup.UserRole = enums.FinancialGroupMember
	}

	return financialGroup, nil
}

func (s *financialGroupService) List(ctx context.Context, input dto.ListRequest, userID uuid.UUID) (*dto.FinancialGroupListResponse, error) {
	return nil, nil
}
