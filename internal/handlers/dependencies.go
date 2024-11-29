package handler

import (
	"shirinec.com/internal/repositories"
)

type Dependencies struct {
	UserRepo     repositories.UserRepository
	CategoryRepo repositories.CategoryRepository
	ItemRepo     repositories.ItemRepository
}
