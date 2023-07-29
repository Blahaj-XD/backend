package backend

import (
	"context"
	"time"

	"github.com/BlahajXD/backend/logic"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthLoginOutput struct {
	AccessToken     string `json:"access_token"`
	BankAccessToken string `json:"bank_access_token"`
	User            struct {
		ID            int    `json:"id"`
		UID           int    `json:"uid"`
		AccountNumber string `json:"account_number"`
		NIK           string `json:"nik"`
		Username      string `json:"username"`
		Email         string `json:"email"`
		PhoneNumber   string `json:"phone_number"`
		FullName      string `json:"full_name"`
		Domisili      string `json:"domisili"`
		TanggalLahir  string `json:"tanggal_lahir"`
		JenisKelamin  int    `json:"jenis_kelamin"`
		Alamat        string `json:"alamat"`
		RtRW          string `json:"rt_rw"`
		Kelurahan     string `json:"kelurahan"`
		Kecamatan     string `json:"kecamatan"`
		Pekerjaan     string `json:"pekerjaan"`
		CreatedAt     string `json:"created_at"`
	} `json:"user"`
}

func (d *Dependency) AuthLogin(ctx context.Context, input AuthLoginInput) (AuthLoginOutput, error) {
	user, err := d.repo.FindParent(ctx, "username", input.Username)
	if err != nil {
		return AuthLoginOutput{}, errors.Wrap(err, "backend.Login -> repo.FindParent")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return AuthLoginOutput{}, ErrInvalidCredentials
	}

	accessToken, err := logic.GenerateJWT(map[string]any{
		"userID": user.ID,
		"email":  user.Email,
	})
	if err != nil {
		return AuthLoginOutput{}, ErrInvalidCredentials
	}

	bankAccessToken, err := d.BankGenerateToken(BankGenerateTokenInput{
		Username:      user.Username,
		LoginPassword: input.Password,
	})
	if err != nil {
		return AuthLoginOutput{}, errors.Wrap(err, "backend.Login -> BankGenerateToken")
	}

	var output AuthLoginOutput
	output.AccessToken = accessToken
	output.BankAccessToken = bankAccessToken
	output.User.ID = user.ID
	output.User.UID = user.UID
	output.User.AccountNumber = user.AccountNumber
	output.User.NIK = user.NIK
	output.User.Username = user.Username
	output.User.Email = user.Email
	output.User.PhoneNumber = user.PhoneNumber
	output.User.FullName = user.FullName
	output.User.Domisili = user.Domisili
	output.User.TanggalLahir = user.TanggalLahir
	output.User.JenisKelamin = user.JenisKelamin
	output.User.Alamat = user.Alamat
	output.User.RtRW = user.RtRW
	output.User.Kelurahan = user.Kelurahan
	output.User.Kecamatan = user.Kecamatan
	output.User.Pekerjaan = user.Pekerjaan
	output.User.CreatedAt = user.CreatedAt.Format(time.RFC3339)

	return output, nil
}
