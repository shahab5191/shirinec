package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/models"
)

type FinancialGroupRepository interface {
	Create(ctx context.Context, financialGroup *models.FinancialGroups) error
	AddUserToGroup(ctx context.Context, financialGroupID int, userID uuid.UUID) error
	GetByID(ctx context.Context, finacialGroupID int, userID uuid.UUID) (*dto.FinancialGroup, error)
	GetOwnedGroupByID(ctx context.Context, financialGroupID int, userID uuid.UUID) (*models.FinancialGroups, error)
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

func (r *financialGroupRepository) GetByID(ctx context.Context, financialGroupID int, userID uuid.UUID) (*dto.FinancialGroup, error) {
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

	if err := r.db.QueryRow(ctx, query, financialGroupID, userID).Scan(&financialGroup.Name, &financialGroup.ImageURL, &financialGroup.UserID, &financialGroup.CreationDate, &financialGroup.UpdateDate, &usersList); err != nil {
        return nil, err
    }
    
    rows, err := r.db.Query(ctx, usersQuery, usersList)
    if err != nil{
        return nil, err
    }

    var users []*dto.UserGetResponse
    for rows.Next(){
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
