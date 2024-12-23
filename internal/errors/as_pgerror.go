package server_errors

import (
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	PGForeignKeyViolation     = "23503"
	PGExceptionDefault        = "P0001"
	PGCategoryNotFound        = "S0001"
	PGInvalidMediaRefrence    = "S0002"
	PGUserAlreadyInGroup      = "S0003"
	PGTransactionUnAuthorized = "S0004"
)

func AsPgError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case PGExceptionDefault:
			return &InvalidInput
		case PGForeignKeyViolation:
			return &InvalidRefrencedEntity
		case PGCategoryNotFound:
			return &CategoryNotFound
		case PGInvalidMediaRefrence:
			return &InvalidMediaRefrence
		case PGUserAlreadyInGroup:
			return &UserAlreadyInFinancialGroup
        case PGTransactionUnAuthorized:
            return &Unauthorized
		default:
			log.Printf("Undefined Postgresql error: %s", pgErr.Error())
		}
	}
	return nil
}
