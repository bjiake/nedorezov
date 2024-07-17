package db

import "errors"

// Обозначение ошибок
var (
	ErrMigrate        = errors.New("migration failed")
	ErrDuplicate      = errors.New("record already exists")
	ErrNotExist       = errors.New("row does not exist")
	ErrUpdateFailed   = errors.New("update failed")
	ErrNotValidAmount = errors.New("amount is not valid")
	ErrParamNotFound  = errors.New("param not found")
	ErrAuthorize      = errors.New("authorize failed")
	ErrValidate       = errors.New("validate failed")
)
