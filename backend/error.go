package backend

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrKidAlreadyExists   = errors.New("kid already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrQuestNotFound = errors.New("quest not found")
)
