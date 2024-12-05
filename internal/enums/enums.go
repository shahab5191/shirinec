package enums

type CategoryType string

const (
	Income  CategoryType = "income"
	Expense CategoryType = "expense"
	Account CategoryType = "account"
)

type UserStatus string

const (
	StatusBanned   UserStatus = "banned"
	StatusVerified UserStatus = "verified"
	StatusDisabled UserStatus = "disabled"
	StatusLocked   UserStatus = "locked"
	StatusPending  UserStatus = "pending"
)

type MediaUploadBind string

const (
	BindToItem     MediaUploadBind = "item"
	BindToProfile  MediaUploadBind = "profile"
	BindToCategory MediaUploadBind = "category"
)
