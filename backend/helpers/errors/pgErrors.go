package customErrors

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Gorm doesn't cover each single error the driver it uses returns. E.g the postgres driver used by gorm is pgx.
// All database errors follow the convention SQLSTATE which assigns error codes to certain error conditions. To test
// for a certain error condition, the code needs to be checked. The codes are found here:
// https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	invalidInputSyntaxErrorCode = "22P02"
	uniqueViolationErrorCode    = "23505"
)

func isPostgresErrorCode(err error, errorCode string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == errorCode
	}
	return false
}

func IsInvalidInputSyntaxError(err error) bool {
	return isPostgresErrorCode(err, invalidInputSyntaxErrorCode)
}

func IsUniqueConstraintViolationError(err error) bool {
	return isPostgresErrorCode(err, uniqueViolationErrorCode)
}
