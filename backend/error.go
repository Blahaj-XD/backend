package backend

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrGoalNotFound       = errors.New("goal not found")
	ErrKidAlreadyExists   = errors.New("kid already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrQuestNotFound = errors.New("quest not found")
)
