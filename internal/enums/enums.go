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

type MediaStatus string

const (
	MediaStatusTemp     MediaStatus = "temp"
	MediaStatusAttached MediaStatus = "attached"
	MediaStatusRemoved  MediaStatus = "removed"
)

type MediaBindType string

const (
	MediaBindItem        MediaBindType = "item"
	MediaBindTransaction MediaBindType = "transaction"
)

type MediaAccess string

const (
	Owner  MediaAccess = "owner"
	Group  MediaAccess = "group"
	Public MediaAccess = "public"
)
