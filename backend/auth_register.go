package backend

import (
	"context"
	"time"

	"github.com/BlahajXD/backend/repo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthRegisterInput struct {
	UID           int    `json:"uid"`
	AccountNumber string `json:"account_number"`
	NIK           string `json:"nik"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Pin           string `json:"pin"`
	Password      string `json:"password"`
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
}

type AuthRegisterOutput struct {
	User struct {
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

func (d *Dependency) AuthRegister(ctx context.Context, input AuthRegisterInput) (AuthRegisterOutput, error) {
	// TODO: Better validation
	// Check if nik already exists
	_, err := d.repo.FindParent(ctx, "nik", input.NIK)
	if err == nil {
		return AuthRegisterOutput{}, ErrUserAlreadyExists
	}

	log.Debug().Err(err).Msgf("input.Email: %s", input.Email)

	// Check if email already exists
	_, err = d.repo.FindParent(ctx, "email", input.Email)
	if err == nil {
		return AuthRegisterOutput{}, ErrUserAlreadyExists
	}

	// Check if username already exists
	_, err = d.repo.FindParent(ctx, "username", input.Username)
	if err == nil {
		return AuthRegisterOutput{}, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthRegisterOutput{}, errors.Wrap(err, "backend.SaveParent: bcrypt.GenerateFromPassword #1")
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(input.Pin), bcrypt.DefaultCost)
	if err != nil {
		return AuthRegisterOutput{}, errors.Wrap(err, "backend.SaveParent: bcrypt.GenerateFromPassword #2")
	}

	input.Password = string(hashedPassword)
	input.Pin = string(hashedPin)
	now := time.Now()
	params := repo.Parent{
		UID:           input.UID,
		AccountNumber: input.AccountNumber,
		NIK:           input.NIK,
		Username:      input.Username,
		Email:         input.Email,
		Pin:           input.Pin,
		Password:      input.Password,
		PhoneNumber:   input.PhoneNumber,
		FullName:      input.FullName,
		Domisili:      input.Domisili,
		TanggalLahir:  input.TanggalLahir,
		JenisKelamin:  input.JenisKelamin,
		Alamat:        input.Alamat,
		RtRW:          input.RtRW,
		Kelurahan:     input.Kelurahan,
		Kecamatan:     input.Kecamatan,
		Pekerjaan:     input.Pekerjaan,
		CreatedAt:     now,
	}

	userID, err := d.repo.SaveParent(ctx, params)
	if err != nil {
		return AuthRegisterOutput{}, errors.Wrap(err, "backend.SaveParent -> d.repo.SaveParent")
	}

	var output AuthRegisterOutput
	output.User.ID = userID
	output.User.UID = input.UID
	output.User.AccountNumber = input.AccountNumber
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
