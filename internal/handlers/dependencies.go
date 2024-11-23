package handler

import (
	"shirinec.com/internal/repositories"
)

type Dependencies struct {
    UserRepo                repositories.UserRepository
    IncomeCategoryRepo      repositories.IncomeCategoryRepository
}
