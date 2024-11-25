package repositories

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, emaill string) (models.User, error)
	UpdatePassword(ctx context.Context, newPassword string, id uuid.UUID) error
	UpdateEmail(ctx context.Context, newEmail string, id uuid.UUID) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	profileQuery := "INSERT INTO profiles DEFAULT VALUES RETURNING id"
	var profileID int
	err := r.db.QueryRow(ctx, profileQuery).Scan(&profileID)
	if err != nil {
        log.Printf("[Error] - userRepository.Create - creating new profile: %+v\n", err)
		return err
	}
    log.Printf("ProfileID: %+v\n", profileID)
	query := "INSERT INTO users (id, email, ip, password, last_login, profile_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	currentTime := time.Now().UTC()
	user.LastLogin = currentTime
	user.CreationDate = currentTime
	user.UpdateDate = currentTime
	user.ProfileID = profileID
	err = r.db.QueryRow(ctx, query, user.ID, user.Email, user.IP, user.Password, currentTime, user.ProfileID).Scan(&user.ID)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	query := "SELECT id, email, ip, password, last_login, failed_tries, status, creation_date, update_date, profile_id FROM users WHERE users.email = $1 LIMIT 1"
	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.IP, &user.Password, &user.LastLogin, &user.FailedTries, &user.Status, &user.CreationDate, &user.UpdateDate, &user.ProfileID)
	return user, err
}

func (r *userRepository) UpdatePassword(ctx context.Context, newPassword string, id uuid.UUID) error {
	query := "UPDATE users SET password = $1 WHERE id = $2"
	err := r.db.QueryRow(ctx, query, newPassword, id).Scan()
	return err
}

func (r *userRepository) UpdateEmail(ctx context.Context, newEmail string, id uuid.UUID) error {
	query := "UPDATE users SET email = $1 WHERE id = $2"
	err := r.db.QueryRow(ctx, query, newEmail, id).Scan()
	return err
}
