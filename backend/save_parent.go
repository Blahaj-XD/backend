package backend

import (
	"context"
	"time"

	"github.com/BlahajXD/backend/logic"
	"github.com/BlahajXD/backend/repo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type SaveParentInput struct {
	NIK          string `json:"nik"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PhoneNumber  string `json:"phone_number"`
	FullName     string `json:"full_name"`
	Domisili     string `json:"domisili"`
	TanggalLahir string `json:"tanggal_lahir"`
	JenisKelamin int    `json:"jenis_kelamin"`
	Alamat       string `json:"alamat"`
	RtRW         string `json:"rt_rw"`
	Kelurahan    string `json:"kelurahan"`
	Kecamatan    string `json:"kecamatan"`
	Pekerjaan    string `json:"pekerjaan"`
}

type SaveParentOutput struct {
	AccessToken string `json:"access_token"`
	User        struct {
		ID           int    `json:"id"`
		NIK          string `json:"nik"`
		Username     string `json:"username"`
		Email        string `json:"email"`
		PhoneNumber  string `json:"phone_number"`
		FullName     string `json:"full_name"`
		Domisili     string `json:"domisili"`
		TanggalLahir string `json:"tanggal_lahir"`
		JenisKelamin int    `json:"jenis_kelamin"`
		Alamat       string `json:"alamat"`
		RtRW         string `json:"rt_rw"`
		Kelurahan    string `json:"kelurahan"`
		Kecamatan    string `json:"kecamatan"`
		Pekerjaan    string `json:"pekerjaan"`
		CreatedAt    string `json:"created_at"`
	} `json:"user"`
}

func (d *Dependency) SaveParent(ctx context.Context, input SaveParentInput) (SaveParentOutput, error) {
	// TODO: Better validation
	// Check if nik already exists
	_, err := d.repo.FindParent(ctx, "nik", input.NIK)
	if err == nil {
		return SaveParentOutput{}, ErrUserAlreadyExists
	}

	log.Debug().Err(err).Msgf("input.Email: %s", input.Email)

	// Check if email already exists
	_, err = d.repo.FindParent(ctx, "email", input.Email)
	if err == nil {
		return SaveParentOutput{}, ErrUserAlreadyExists
	}

	// Check if username already exists
	_, err = d.repo.FindParent(ctx, "username", input.Username)
	if err == nil {
		return SaveParentOutput{}, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return SaveParentOutput{}, errors.Wrap(err, "backend.SaveParent: bcrypt.GenerateFromPassword")
	}

	input.Password = string(hashedPassword)
	now := time.Now()
	params := repo.Parent{
		NIK:          input.NIK,
		Username:     input.Username,
		Email:        input.Email,
		Password:     input.Password,
		PhoneNumber:  input.PhoneNumber,
		FullName:     input.FullName,
		Domisili:     input.Domisili,
		TanggalLahir: input.TanggalLahir,
		JenisKelamin: input.JenisKelamin,
		Alamat:       input.Alamat,
		RtRW:         input.RtRW,
		Kelurahan:    input.Kelurahan,
		Kecamatan:    input.Kecamatan,
		Pekerjaan:    input.Pekerjaan,
		CreatedAt:    now,
	}

	userID, err := d.repo.SaveParent(ctx, params)
	if err != nil {
		return SaveParentOutput{}, errors.Wrap(err, "backend.SaveParent -> d.repo.SaveParent")
	}

	accessToken, err := logic.GenerateJWT(map[string]any{
		"userID": userID,
		"email":  input.Email,
	})
	if err != nil {
		return SaveParentOutput{}, errors.Wrap(err, "backend.SaveParent -> logic.GenerateJWT")
	}

	var output SaveParentOutput
	output.AccessToken = accessToken
	output.User.ID = userID
	output.User.NIK = input.NIK
	output.User.Username = input.Username
	output.User.Email = input.Email
	output.User.PhoneNumber = input.PhoneNumber
	output.User.FullName = input.FullName
	output.User.Domisili = input.Domisili
	output.User.TanggalLahir = input.TanggalLahir
	output.User.JenisKelamin = input.JenisKelamin
	output.User.Alamat = input.Alamat
	output.User.RtRW = input.RtRW
	output.User.Kelurahan = input.Kelurahan
	output.User.Kecamatan = input.Kecamatan
	output.User.Pekerjaan = input.Pekerjaan
	output.User.CreatedAt = now.Format(time.RFC3339)

	return output, nil
}
