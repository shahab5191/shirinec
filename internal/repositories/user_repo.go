package repositories

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
    GetByEmail(ctx context.Context, emaill string) (models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (id, name, email, ip, password) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := r.db.QueryRow(ctx, query, user.ID, user.Name, user.Email, user.IP, user.Password).Scan(&user.ID)
    return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (models.User, error){
    query := "SELECT id, email, ip, password FROM users WHERE users.email = $1 LIMIT 1"
    var user models.User
    err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.IP, &user.Password)
    if err != nil {
        log.Printf("ERROR: %s", err)
        return user, err
    }
    return user, nil
}
