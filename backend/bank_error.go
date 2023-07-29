package backend

import "errors"

var (
	ErrBankRequestParameter = errors.New("bank api: request_parameter_error 1020")
	ErrHasNoDataPermission  = errors.New("bank api: has_no_data_permission 1025")
	ErrNotFoundBankAccount  = errors.New("bank api: not_found_bank_account 4002")
	ErrBankUnknownError     = errors.New("bank api: unknown_error 9999")
)
