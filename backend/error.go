package backend

import "errors"

var (
	ErrUserAlreadyExists               = errors.New("user already exists")
	ErrUserNotFound                    = errors.New("user not found")
	ErrGoalNotFound                    = errors.New("goal not found")
	ErrQuestNotFound                   = errors.New("quest not found")
	ErrQuestNotAvailable               = errors.New("quest not available")
	ErrCantDepositToGoalWhenNotOngoing = errors.New("can't deposit to goal when not ongoing")
	ErrKidAlreadyExists                = errors.New("kid already exists")
	ErrInvalidCredentials              = errors.New("invalid credentials")
)
