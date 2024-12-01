package enums

type CategoryType string

const (
	Income  CategoryType = "Income"
	Expense CategoryType = "Expense"
	Account CategoryType = "Account"
)

type UserStatus string

const (
	StatusBanned   UserStatus = "Banned"
	StatusVerified UserStatus = "Verified"
	StatusDisabled UserStatus = "Disabled"
	StatusLocked   UserStatus = "Locked"
	StatusPending  UserStatus = "Pending"
)

type MediaUploadBind string

const (
	BindToItem     MediaUploadBind = "item"
	BindToProfile  MediaUploadBind = "profile"
	BindToCategory MediaUploadBind = "category"
)
