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
