package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/src/internal/dto"
	server_errors "shirinec.com/src/internal/errors"
	"shirinec.com/src/internal/models"
)

type FinancialGroupRepository interface {
	Create(ctx context.Context, financialGroup *models.FinancialGroups) error
	AddUserToGroup(ctx context.Context, financialGroupID int, userID uuid.UUID) error
	GetRelatedGroupByID(ctx context.Context, finacialGroupID int, userID uuid.UUID) (*dto.FinancialGroup, error)
	GetOwnedGroupByID(ctx context.Context, financialGroupID int, userID uuid.UUID) (*models.FinancialGroups, error)
	ListOwnedGroups(ctx context.Context, page, size int, userID uuid.UUID) ([]dto.FinancialGroupListItem, int, error)
	ListMemberedGroups(ctx context.Context, page, size int, userID uuid.UUID) ([]dto.FinancialGroupListItem, int, error)
	RemoveGroupMember(ctx context.Context, financialGroupID int, memberID, userID uuid.UUID) error
	Delete(ctx context.Context, financialGroupID int) error
	GetByID(ctx context.Context, financialGroupID int) (*models.FinancialGroups, error)
}

type financialGroupRepository struct {
	db *pgxpool.Pool
}

func NewFinancialGroupRepository(db *pgxpool.Pool) FinancialGroupRepository {
	return &financialGroupRepository{
		db: db,
	}
}

func (r *financialGroupRepository) Create(ctx context.Context, financialGroup *models.FinancialGroups) error {
	query := `
        INSERT INTO financial_groups (name, image_id, user_id, creation_date, update_date)
        VALUES ($1, $2, $3, $4, $4)
        RETURNING id
    `
	err := r.db.QueryRow(ctx, query, financialGroup.Name, financialGroup.ImageID, financialGroup.UserID, financialGroup.CreationDate).Scan(&financialGroup.ID)
	return err
}

func (r *financialGroupRepository) AddUserToGroup(ctx context.Context, financialGroupID int, userID uuid.UUID) error {
	query := `
        INSERT INTO user_financial_groups 
        (financial_group_id, user_id)
        VALUES ($1, $2)
        RETURNING id
    `
	var relID int
	err := r.db.QueryRow(ctx, query, financialGroupID, userID).Scan(&relID)
	return err
}

func (r *financialGroupRepository) GetRelatedGroupByID(ctx context.Context, financialGroupID int, userID uuid.UUID) (*dto.FinancialGroup, error) {
	var financialGroup dto.FinancialGroup
	var usersList []uuid.UUID
	query := `
        SELECT
            fg.name, m.url, fg.user_id, fg.creation_date, fg.update_date, ARRAY_AGG(ufg.user_id) AS members
        FROM financial_groups fg
        LEFT JOIN media m ON m.id = fg.image_id
        LEFT JOIN user_financial_groups ufg ON ufg.financial_group_id = fg.id
        WHERE
            fg.id = $1
            AND (
                fg.user_id = $2
                OR
                EXISTS(
                    SELECT 1 FROM user_financial_groups
                    WHERE financial_group_id = $1
                    AND user_id = $2
                )
            )
        GROUP BY
            fg.name, m.url, fg.user_id, fg.creation_date, fg.update_date;
    `
	usersQuery := `
        SELECT u.id, pm.url profie_picture, p.name, p.family_name
        FROM users u
        JOIN profiles p ON p.id = u.profile_id
        LEFT JOIN media pm ON pm.id = p.picture_id
        WHERE u.id = ANY($1)
    `

	if err := r.db.QueryRow(ctx, query, financialGroupID, userID).Scan(&financialGroup.Name, &financialGroup.ImageURL, &financialGroup.OwnerID, &financialGroup.CreationDate, &financialGroup.UpdateDate, &usersList); err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, usersQuery, usersList)
	if err != nil {
		return nil, err
	}

	var users []*dto.UserGetResponse
	for rows.Next() {
		var user dto.UserGetResponse
		if err := rows.Scan(&user.ID, &user.ProfilePictureURL, &user.Name, &user.FamilyName); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	financialGroup.Users = users
	return &financialGroup, err
}

func (r *financialGroupRepository) GetOwnedGroupByID(ctx context.Context, financialGroupID int, userID uuid.UUID) (*models.FinancialGroups, error) {
	var financialGroup models.FinancialGroups
	query := `
        SELECT name, image_id, user_id, creation_date, update_date
        FROM financial_groups
        WHERE id = $1 AND user_id = $2
    `

	err := r.db.QueryRow(ctx, query, financialGroupID, userID).Scan(&financialGroup.Name, &financialGroup.ImageID, &financialGroup.UserID, &financialGroup.CreationDate, &financialGroup.UpdateDate)

	return &financialGroup, err
}

func (r *financialGroupRepository) ListOwnedGroups(ctx context.Context, limit, offset int, userID uuid.UUID) ([]dto.FinancialGroupListItem, int, error) {
	query := `
        SELECT id, name
        FROM financial_groups
        WHERE user_id = $1
        LIMIT $2
        OFFSET $3
    `

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var financialGroups []dto.FinancialGroupListItem
	for rows.Next() {
		var financialGroup dto.FinancialGroupListItem
		if err := rows.Scan(&financialGroup.ID, &financialGroup.Name); err != nil {
			return nil, 0, err
		}
		financialGroups = append(financialGroups, financialGroup)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := `
        SELECT COUNT(*) FROM financial_groups WHERE user_id = $1
    `

	var totalCount int

	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	return financialGroups, totalCount, nil
}

func (r *financialGroupRepository) ListMemberedGroups(ctx context.Context, limit, offset int, userID uuid.UUID) ([]dto.FinancialGroupListItem, int, error) {
	query := `
        SELECT fg.id, fg.name as role
        FROM user_financial_groups ufg
        JOIN financial_groups fg ON fg.id = ufg.financial_group_id
        WHERE ufg.user_id = $1
        LIMIT $2
        OFFSET $3
    `

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var financialGroups []dto.FinancialGroupListItem
	for rows.Next() {
		var financialGroup dto.FinancialGroupListItem
		if err := rows.Scan(&financialGroup.ID, &financialGroup.Name); err != nil {
			return nil, 0, err
		}
		financialGroups = append(financialGroups, financialGroup)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := `
        SELECT COUNT(*) FROM user_financial_groups WHERE user_id = $1
    `

	var totalCount int

	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	return financialGroups, totalCount, nil
}

func (r *financialGroupRepository) RemoveGroupMember(ctx context.Context, financialGroupID int, memberID, userID uuid.UUID) error {
	if userID != memberID {
		getOwnerQuery := `
            SELECT user_id
            FROM financial_groups
            WHERE id = $1
        `
		var ownerID uuid.UUID
		if err := r.db.QueryRow(ctx, getOwnerQuery, financialGroupID).Scan(&ownerID); err != nil {
			return err
		}
		if ownerID != userID {
			return &server_errors.Unauthorized
		}
	}

	query := `
        DELETE FROM user_financial_groups
        WHERE
            financial_group_id = $1
            AND
            user_id = $2
        RETURNING id
    `

	var id int
	err := r.db.QueryRow(ctx, query, financialGroupID, memberID).Scan(&id)
	return err
}

func (r *financialGroupRepository) Delete(ctx context.Context, financialGroupID int) error {
	query := `
        DELETE FROM financial_groups
        WHERE id = $1
        RETURNING id
    `
	var id int
	err := r.db.QueryRow(ctx, query, financialGroupID).Scan(&id)
	return err
}

func (r *financialGroupRepository) GetByID(ctx context.Context, financialGroupID int) (*models.FinancialGroups, error) {
	query := `
        SELECT id, name, user_id, image_id, creation_date, update_date
        FROM financial_groups
        WHERE id = $1
    `

	var financialGroup models.FinancialGroups
	err := r.db.QueryRow(ctx, query, financialGroupID).Scan(&financialGroup.ID, &financialGroup.Name, &financialGroup.UserID, &financialGroup.ImageID, &financialGroup.CreationDate, &financialGroup.UpdateDate)
	return &financialGroup, err
}
