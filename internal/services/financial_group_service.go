package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

type FinancialGroupService interface {
	Create(ctx context.Context, input *dto.FinancialGroupCreateRequest, userID uuid.UUID) (*models.FinancialGroups, error)
	AddUserToGroup(ctx context.Context, financialGroupID int, newUserID uuid.UUID, userID uuid.UUID) error
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
	financialGroup, err := s.financialGroupRepo.GetByID(ctx, financialGroupID)
	if err != nil {
		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return pgErr
		}
		utils.Logger.Errorf("financialGroupService.AddUserToGroup - Calling financialGroupRep.GetByID(%d): %s", financialGroupID, err)
		return &server_errors.InternalError
	}

	if financialGroup.UserID != userID {
		return &server_errors.Unauthorized
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
