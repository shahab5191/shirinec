package handler

import (
	"shirinec.com/src/internal/repositories"
)

type Dependencies struct {
	UserRepo           repositories.UserRepository
	CategoryRepo       repositories.CategoryRepository
	ItemRepo           repositories.ItemRepository
	AccountRepo        repositories.AccountRepository
	MediaRepo          repositories.MediaRepository
	FinancialGroupRepo repositories.FinancialGroupRepository
	TransactionRepo    repositories.TransactionRepository
}
