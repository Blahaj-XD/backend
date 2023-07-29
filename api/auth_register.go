package api

import (
	"errors"
	"net/mail"

	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type AuthRegisterBody struct {
	NIK          string `json:"nik"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Pin          string `json:"pin"`
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

func (s *Server) AuthRegister(c *fiber.Ctx) error {
	var body AuthRegisterBody
	if err := c.BodyParser(&body); err != nil {
		log.Debug().Msg("Error parsing body")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := body.Validate(); err != nil {
		log.Debug().Msg("Error validating body")
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	log.Debug().Msg("Registering user in bank 3rd party API")
	uid, err := s.backend.BankCreateUser(backend.BankCreateUserInput{
		KTPID:         body.NIK,
		Username:      body.Username,
		Email:         body.Email,
		LoginPassword: body.Password,
		PhoneNumber:   body.PhoneNumber,
		BirthDate:     body.TanggalLahir,
		Gender:        body.JenisKelamin,
	})
	if err != nil {
		if errors.Is(err, backend.ErrBankRequestParameter) {
			log.Debug().Msg("Error registering user in bank 3rd party API: ErrBankRequestParameter")
			errMsg := errors.Unwrap(err)
			return fiber.NewError(fiber.StatusUnprocessableEntity, errMsg.Error())
		}

		if errors.Is(err, backend.ErrUserAlreadyExists) {
			log.Debug().Msg("Error registering user in bank 3rd party API: ErrUserAlreadyExists")
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}

		log.Debug().Msg("Error registering user in bank 3rd party API: ErrServiceUnavailable")
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	accessToken, err := s.backend.BankGenerateToken(backend.BankGenerateTokenInput{
		Username:      body.Username,
		LoginPassword: body.Password,
	})
	if err != nil {
		if errors.Is(err, backend.ErrHasNoDataPermission) {
			log.Debug().Msg("Error generating token in bank 3rd party API: ErrHasNoDataPermission")
			errMsg := errors.Unwrap(err)
			return fiber.NewError(fiber.StatusUnprocessableEntity, errMsg.Error())
		}

		log.Debug().Msg("Error generating token in bank 3rd party API: ErrServiceUnavailable")
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	accountNumber, err := s.backend.BankCreateBankAccount(backend.BankCreateBankAccountInput{
		Balance:     0,
		AccessToken: accessToken,
	})
	if err != nil {
		if errors.Is(err, backend.ErrHasNoDataPermission) {
			log.Debug().Msg("Error creating bank account in bank 3rd party API: ErrHasNoDataPermission")
			errMsg := errors.Unwrap(err)
			return fiber.NewError(fiber.StatusUnprocessableEntity, errMsg.Error())
		}

		log.Debug().Msg("Error creating bank account in bank 3rd party API: ErrServiceUnavailable")
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	log.Debug().Msg("Registering user in database")

	output, err := s.backend.AuthRegister(c.Context(), backend.AuthRegisterInput{
		UID:           uid,
		AccountNumber: accountNumber,
		NIK:           body.NIK,
		Username:      body.Username,
		Email:         body.Email,
		Pin:           body.Pin,
		Password:      body.Password,
		PhoneNumber:   body.PhoneNumber,
		FullName:      body.FullName,
		Domisili:      body.Domisili,
		TanggalLahir:  body.TanggalLahir,
		JenisKelamin:  body.JenisKelamin,
		Alamat:        body.Alamat,
		RtRW:          body.RtRW,
		Kelurahan:     body.Kelurahan,
		Kecamatan:     body.Kecamatan,
		Pekerjaan:     body.Pekerjaan,
	})
	if err != nil {
		if errors.Is(err, backend.ErrUserAlreadyExists) {
			log.Debug().Msg("Error registering user in database: ErrUserAlreadyExists")
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}

		log.Debug().Msg("Error registering user in database: ErrInternalServerError")
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

func (b AuthRegisterBody) Validate() error {
	if b.NIK == "" {
		return errors.New("nik is required")
	}

	if b.Username == "" {
		return errors.New("username is required")
	}

	if b.Email == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(b.Email); err != nil {
		return errors.New("invalid email")
	}

	if b.Pin == "" {
		return errors.New("pin is required")
	}

	if b.Pin != "" && len(b.Pin) != 6 {
		return errors.New("pin must be 6 digits")
	}

	if b.Password == "" {
		return errors.New("password is required")
	}

	if b.PhoneNumber == "" {
		return errors.New("phone_number is required")
	}

	if b.FullName == "" {
		return errors.New("full_name is required")
	}

	if b.Domisili == "" {
		return errors.New("domisili is required")
	}

	if b.TanggalLahir == "" {
		return errors.New("tanggal_lahir is required")
	}

	if b.JenisKelamin < 0 && b.JenisKelamin > 1 {
		return errors.New("jenis_kelamin is required (0: Laki-laki, 1: Perempuan)")
	}

	if b.Alamat == "" {
		return errors.New("alamat is required")
	}

	if b.RtRW == "" {
		return errors.New("rt_rw is required")
	}

	if b.Kelurahan == "" {
		return errors.New("kelurahan is required")
	}

	if b.Kecamatan == "" {
		return errors.New("kecamatan is required")
	}

	if b.Pekerjaan == "" {
		return errors.New("pekerjaan is required")
	}

	return nil
}
