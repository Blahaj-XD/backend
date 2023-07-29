package backend

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type ParentAdminBankAccountInfoInput struct {
	ParentID    int `json:"parent_id"`
	AccessToken string
}

type ParentAdminBankAccountInfoOutput struct {
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
	User          struct {
		ID            int       `json:"id"`
		UID           int       `json:"uid"`
		AccountNumber string    `json:"account_number"`
		NIK           string    `json:"nik"`
		Username      string    `json:"username"`
		Email         string    `json:"email"`
		PhoneNumber   string    `json:"phone_number"`
		FullName      string    `json:"full_name"`
		Domisili      string    `json:"domisili"`
		TanggalLahir  string    `json:"tanggal_lahir"`
		JenisKelamin  int       `json:"jenis_kelamin"`
		Alamat        string    `json:"alamat"`
		RtRW          string    `json:"rt_rw"`
		Kelurahan     string    `json:"kelurahan"`
		Kecamatan     string    `json:"kecamatan"`
		Pekerjaan     string    `json:"pekerjaan"`
		CreatedAt     time.Time `json:"created_at"`
	} `json:"user"`
}

func (d *Dependency) ParentAdminBankAccountInfo(ctx context.Context, input ParentAdminBankAccountInfoInput) (ParentAdminBankAccountInfoOutput, error) {
	parent, err := d.repo.FindParent(ctx, "id", input.ParentID)
	if err != nil {
		return ParentAdminBankAccountInfoOutput{}, errors.Wrap(ErrUserNotFound, "parent_admin_bank_account_info: failed to find parent")
	}

	account, err := d.BankAccountInfo(ctx, BankAccountInfoInput{
		AccountNumber: parent.AccountNumber,
		AccessToken:   input.AccessToken,
	})
	if err != nil {
		return ParentAdminBankAccountInfoOutput{}, errors.Wrap(err, "parent_admin_bank_account_info: failed to get bank account info")
	}

	var output ParentAdminBankAccountInfoOutput
	output.AccountNumber = account.AccountNumber
	output.Balance = account.Balance
	output.User.ID = parent.ID
	output.User.UID = parent.UID
	output.User.AccountNumber = parent.AccountNumber
	output.User.NIK = parent.NIK
	output.User.Username = parent.Username
	output.User.Email = parent.Email
	output.User.PhoneNumber = parent.PhoneNumber
	output.User.FullName = parent.FullName
	output.User.Domisili = parent.Domisili
	output.User.TanggalLahir = parent.TanggalLahir
	output.User.JenisKelamin = parent.JenisKelamin
	output.User.Alamat = parent.Alamat
	output.User.RtRW = parent.RtRW
	output.User.Kelurahan = parent.Kelurahan
	output.User.Kecamatan = parent.Kecamatan
	output.User.Pekerjaan = parent.Pekerjaan
	output.User.CreatedAt = parent.CreatedAt

	return output, nil
}
