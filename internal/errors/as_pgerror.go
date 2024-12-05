package server_errors

import (
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	PGForeignKeyViolation = "23503"
	PGExceptionDefault    = "P0001"
	PGCategoryNotFound    = "S0001"
)

func AsPgError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case PGExceptionDefault:
			return &InvalidInput
		case PGForeignKeyViolation:
			return &InvalidInput
		case PGCategoryNotFound:
			return &CategoryNotFound
        default:
            log.Printf("Undefined Postgresql error: %s", pgErr.Error())
		}
	}
	return nil
}
