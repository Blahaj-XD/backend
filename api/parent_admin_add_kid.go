package api

import (
	"errors"

	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

type ParentAdminAddKidBody struct {
	NIK          string `json:"nik"`
	FullName     string `json:"full_name"`
	Domisili     string `json:"domisili"`
	TanggalLahir string `json:"tanggal_lahir"`
	JenisKelamin int    `json:"jenis_kelamin"`
}

func (s *Server) ParentAdminAddKid(c *fiber.Ctx) error {
	var body ParentAdminAddKidBody
	parentID := int(c.Locals("userID").(float64))
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := body.Validate(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	accountNumber, err := s.backend.HackathonCreateBankAccount(backend.HackathonCreateBankAccountInput{
		Balance:     0,
		AccessToken: c.Locals("hackathonAccessToken").(string),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	kid, err := s.backend.AddKid(c.Context(), backend.ParentAdminAddKidInput{
		ParentID:      parentID,
		AccountNumber: accountNumber,
		NIK:           body.NIK,
		FullName:      body.FullName,
		Domisili:      body.Domisili,
		TanggalLahir:  body.TanggalLahir,
		JenisKelamin:  body.JenisKelamin,
	})
	if err != nil {
		if errors.Is(err, backend.ErrKidAlreadyExists) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(kid)
}

func (b ParentAdminAddKidBody) Validate() error {
	if b.NIK == "" {
		return errors.New("nik is required")
	}

	if b.FullName == "" {
		return errors.New("full name is required")
	}

	if b.Domisili == "" {
		return errors.New("domisili is required")
	}

	if b.TanggalLahir == "" {
		return errors.New("tanggal lahir is required")
	}

	if b.JenisKelamin < 0 || b.JenisKelamin > 1 {
		return errors.New("jenis kelamin must be 0 or 1")
	}

	return nil
}
