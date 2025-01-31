package services

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"shirinec.com/src/internal/dto"
	"shirinec.com/src/internal/enums"
	"shirinec.com/src/internal/errors"
	"shirinec.com/src/internal/models"
	"shirinec.com/src/internal/repositories"
	"shirinec.com/src/internal/utils"
)

type FinancialGroupService interface {
	Create(ctx context.Context, input *dto.FinancialGroupCreateRequest, userID uuid.UUID) (*models.FinancialGroups, error)
	AddUserToGroup(ctx context.Context, financialGroupID int, newUserID uuid.UUID, userID uuid.UUID) error
	GetByID(ctx context.Context, id int, userID uuid.UUID) (*dto.FinancialGroup, error)
	List(ctx context.Context, input dto.FinancialGroupListRequest, userID uuid.UUID) (*dto.FinancialGroupListResponse, error)
	RemoveGroupMember(ctx context.Context, financialGroupID int, memberID, userID uuid.UUID) error
	Delete(ctx context.Context, financialGroupID int, userID uuid.UUID) error
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
	financialGroup, err := s.financialGroupRepo.GetRelatedGroupByID(ctx, id, userID)
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

func (s *financialGroupService) List(ctx context.Context, input dto.FinancialGroupListRequest, userID uuid.UUID) (*dto.FinancialGroupListResponse, error) {
	limit := input.Size
	offset := input.Page * input.Size

	var financialGroups []dto.FinancialGroupListItem
	var totalCount int
	var err error
	switch input.Role {
	case enums.FinancialGroupOwner:
		financialGroups, totalCount, err = s.financialGroupRepo.ListOwnedGroups(ctx, limit, offset, userID)
	case enums.FinancialGroupMember:
		financialGroups, totalCount, err = s.financialGroupRepo.ListMemberedGroups(ctx, limit, offset, userID)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &server_errors.ItemNotFound
		}

		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return nil, pgErr
		}

		utils.Logger.Errorf("financialGroupService.List - Calling financialGroupRepo.List: %s", err.Error())
		return nil, &server_errors.InternalError
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(input.Size)))
	remainingPages := int(math.Max(float64(totalPages-input.Page-1), 0))

	var response dto.FinancialGroupListResponse

	response.FinancialGroup = financialGroups
	response.Pagination.PageNumber = input.Page
	response.Pagination.PageSize = input.Size
	response.Pagination.RemainingPages = remainingPages
	response.Pagination.TotalRecord = totalCount
	return &response, nil
}

func (s *financialGroupService) RemoveGroupMember(ctx context.Context, financialGroupID int, memberID, userID uuid.UUID) error {
	if err := s.financialGroupRepo.RemoveGroupMember(ctx, financialGroupID, memberID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &server_errors.ItemNotFound
		}

		var sErr *server_errors.SError
		if errors.As(err, &sErr) {
			return sErr
		}

		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return pgErr
		}

		utils.Logger.Errorf("financialGroupService.RemoveUserFromGroup - Calling financialGroupRepo.RemoveUserFromGroup: %s", err.Error())
		return &server_errors.InternalError
	}

	return nil
}

func (s *financialGroupService) Delete(ctx context.Context, financialGroupID int, userID uuid.UUID) error {
	financialGroup, err := s.financialGroupRepo.GetByID(ctx, financialGroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &server_errors.ItemNotFound
		}

		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return pgErr
		}

		utils.Logger.Errorf("financialGroupService.Delete - Calling financialGroupRepo.GetOwnedGroupByID: %s", err.Error())
		return &server_errors.InternalError
	}

	if financialGroup.UserID != userID {
		return &server_errors.Unauthorized
	}

	if err := s.financialGroupRepo.Delete(ctx, financialGroupID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &server_errors.ItemNotFound
		}

		if pgErr := server_errors.AsPgError(err); pgErr != nil {
			return pgErr
		}

		utils.Logger.Errorf("financialGroupService.Delete - Calling financialGroupRepo.Delete: %s", err.Error())
		return &server_errors.InternalError
	}

	return nil
}
