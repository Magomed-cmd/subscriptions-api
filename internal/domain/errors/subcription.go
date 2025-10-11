package errors2

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Доменные ошибки
var (
	ErrNotFound            = errors.New("record not found")
	ErrDuplicate           = errors.New("duplicate record")
	ErrInvalidInput        = errors.New("invalid input data")
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrNullViolation       = errors.New("null value not allowed")
	ErrCheckViolation      = errors.New("check constraint violated")
	ErrInvalidUUID         = errors.New("invalid UUID format")
	ErrInvalidDateFormat   = errors.New("invalid date format")
	ErrForeignKeyViolation = errors.New("foreign key constraint violated")
	ErrInternal            = errors.New("internal server error")
	ErrNothingToUpdate     = errors.New("nothing to update")
)

var pgErrorMap = map[string]error{
	"23502": ErrNullViolation,       // not_null_violation
	"23514": ErrCheckViolation,      // check_violation
	"22P02": ErrInvalidUUID,         // invalid_text_representation
	"22008": ErrInvalidDateFormat,   // invalid_datetime_format
	"23505": ErrDuplicate,           // unique_violation
	"23503": ErrForeignKeyViolation, // foreign_key_violation
}

func MapPgError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if mapped, ok := pgErrorMap[pgErr.Code]; ok {
			return mapped
		}
		return ErrInternal
	}
	return err
}
