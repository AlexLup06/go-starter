package customErrors

import "errors"

var (
	ErrRecordNotFound           = errors.New("record not found")
	ErrDuplicatedKey            = errors.New("duplicated key not allowed")
	ErrEmailAndProviderExist    = errors.New("email with the provider already exists")
	ErrEmailExists              = errors.New("email exists")
	ErrAuthProviderExists       = errors.New("Auth provider already exists for the email")
	ErrAuthProviderDoesNotExist = errors.New("The auth Provider does not exist")
	ErrUserDoesNotExist         = errors.New("user does not exist")
)

type baseErr struct {
	message string
}

func (b baseErr) Error() string {
	return b.message
}

type ValidationError struct {
	baseErr
	// originalErr contains the originalErr error causing this validation error.
	originalErr error
}

func (v ValidationError) Original() error {
	return v.originalErr
}

func NewValidationError(message string) error {
	return ValidationError{baseErr: baseErr{message: message}}
}

func NewValidationErrorWithOriginal(message string, originalErr error) error {
	return ValidationError{baseErr: baseErr{message: message}, originalErr: originalErr}
}

func IsValidationError(err error) bool {
	return errors.As(err, &ValidationError{})
}

type UnauthorizedError struct {
	baseErr
}

func NewUnauthorizedError(message string) error {
	return UnauthorizedError{baseErr: baseErr{message: message}}
}

func IsUnauthorizedError(err error) bool {
	return errors.As(err, &UnauthorizedError{})
}

type NotFoundError struct {
	baseErr
}

func NewNotFoundError(message string) error {
	return NotFoundError{baseErr: baseErr{message: message}}
}

func IsNotFoundError(err error) bool {
	return errors.As(err, &NotFoundError{})
}
