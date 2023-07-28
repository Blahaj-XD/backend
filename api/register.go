package api

import (
	"errors"
	"net/mail"

	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

type RegisterBody struct {
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

func (s *Server) Register(c *fiber.Ctx) error {
	var body RegisterBody

	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := body.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	output, err := s.backend.SaveParent(c.Context(), backend.SaveParentInput{
		NIK:          body.NIK,
		Username:     body.Username,
		Email:        body.Email,
		Password:     body.Password,
		PhoneNumber:  body.PhoneNumber,
		FullName:     body.FullName,
		Domisili:     body.Domisili,
		TanggalLahir: body.TanggalLahir,
		JenisKelamin: body.JenisKelamin,
		Alamat:       body.Alamat,
		RtRW:         body.RtRW,
		Kelurahan:    body.Kelurahan,
		Kecamatan:    body.Kecamatan,
		Pekerjaan:    body.Pekerjaan,
	})
	if err != nil {
		if errors.Is(err, backend.ErrUserAlreadyExists) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

func (b RegisterBody) Validate() error {
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

	if b.JenisKelamin >= 0 && b.JenisKelamin <= 1 {
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
