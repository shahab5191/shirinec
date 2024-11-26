package enums

type CategoryType string

const (
	Income  CategoryType = "Income"
	Expense CategoryType = "Expense"
)

type UserStatus string

const (
	StatusBanned UserStatus = "Banned"
	StatusValidated UserStatus = "Validated"
	StatusDisabled UserStatus = "Disabled"
	StatusLocked UserStatus = "Locked"
	StatusPending UserStatus = "Pending"
)
