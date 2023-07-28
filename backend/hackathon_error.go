package backend

import "errors"

var (
	ErrHackathonRequestParameter = errors.New("hackathon api: request_parameter_error 1020")
	ErrHasNoDataPermission       = errors.New("hackathon api: has_no_data_permission 1025")
	ErrNotFoundBankAccount       = errors.New("hackathon api: not_found_bank_account 4002")
	ErrHackathonUnknownError     = errors.New("hackathon api: unknown_error 9999")
)
