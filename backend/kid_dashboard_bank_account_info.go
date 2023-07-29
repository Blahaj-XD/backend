package backend

import (
	"context"
	"time"
)

type KidDashboardBankAccountInfoInput struct {
	KidID       int `json:"kid_id"`
	ParentID    int `json:"parent_id"`
	AccessToken string
}

type KidDashboardBankAccountInfoOutput struct {
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
	User          struct {
		ID            int       `json:"id"`
		ParentID      int       `json:"parent_id"`
		AccountNumber string    `json:"account_number"`
		NIK           string    `json:"nik"`
		FullName      string    `json:"full_name"`
		Domisili      string    `json:"domisili"`
		TanggalLahir  string    `json:"tanggal_lahir"`
		JenisKelamin  int       `json:"jenis_kelamin"`
		CreatedAt     time.Time `json:"created_at"`
	} `json:"user"`
}

func (d *Dependency) KidDashboardBankAccountInfo(ctx context.Context, input KidDashboardBankAccountInfoInput) (KidDashboardBankAccountInfoOutput, error) {
	kid, err := d.repo.FindKid(ctx, "id", input.KidID)
	if err != nil {
		return KidDashboardBankAccountInfoOutput{}, ErrUserNotFound
	}

	if kid.ParentID != input.ParentID {
		return KidDashboardBankAccountInfoOutput{}, ErrUserNotFound
	}

	info, err := d.BankAccountInfo(ctx, BankAccountInfoInput{
		AccountNumber: kid.AccountNumber,
		AccessToken:   input.AccessToken,
	})
	if err != nil {
		return KidDashboardBankAccountInfoOutput{}, err
	}

	var output KidDashboardBankAccountInfoOutput
	output.AccountNumber = info.AccountNumber
	output.Balance = info.Balance
	output.User.ID = kid.ID
	output.User.ParentID = kid.ParentID
	output.User.AccountNumber = kid.AccountNumber
	output.User.NIK = kid.NIK
	output.User.FullName = kid.FullName
	output.User.Domisili = kid.Domisili
	output.User.TanggalLahir = kid.TanggalLahir
	output.User.JenisKelamin = kid.JenisKelamin
	output.User.CreatedAt = kid.CreatedAt
	return output, nil
}
