package service

// Error code constants — replace all hardcoded error codes with these.
const (
	ErrCodeInvalidParams      = 40001
	ErrCodeInvalidReminderID  = 40002
	ErrCodeIncorrectPassword  = 40003
	ErrCodeUsernameRequired   = 40004
	ErrCodeEmailInvalid       = 40005
	ErrCodePasswordRequired   = 40006
	ErrCodeCredentialsRequired = 40007
	ErrCodeUnauthorized       = 40100
	ErrCodeInvalidCredentials = 40101
	ErrCodeForbidden          = 40301
	ErrCodeNotFound           = 40400
	ErrCodeDuplicateUsername  = 40901
	ErrCodeDuplicateEmail     = 40902
	ErrCodeDuplicateCategory  = 40903
)
