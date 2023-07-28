package backend

import (
	"context"
	"time"

	"github.com/BlahajXD/backend/repo"
	"github.com/pkg/errors"
)

type ParentAdminAddKidInput struct {
	ParentID     int    `json:"parent_id"`
	NIK          string `json:"nik"`
	FullName     string `json:"full_name"`
	Domisili     string `json:"domisili"`
	TanggalLahir string `json:"tanggal_lahir"`
	JenisKelamin int    `json:"jenis_kelamin"`
}

type ParentAdminAddKidOutput struct {
	ID           int       `json:"id"`
	ParentID     int       `json:"parent_id"`
	NIK          string    `json:"nik"`
	FullName     string    `json:"full_name"`
	Domisili     string    `json:"domisili"`
	TanggalLahir string    `json:"tanggal_lahir"`
	JenisKelamin int       `json:"jenis_kelamin"`
	CreatedAt    time.Time `json:"created_at"`
}

func (d *Dependency) AddKid(ctx context.Context, input ParentAdminAddKidInput) (ParentAdminAddKidOutput, error) {
	var output ParentAdminAddKidOutput

	_, err := d.repo.FindKid(ctx, "nik", input.NIK)
	if err == nil {
		return output, ErrKidAlreadyExists
	}

	params := repo.Kid{
		ParentID:     input.ParentID,
		NIK:          input.NIK,
		FullName:     input.FullName,
		Domisili:     input.Domisili,
		TanggalLahir: input.TanggalLahir,
		JenisKelamin: input.JenisKelamin,
		CreatedAt:    time.Now(),
	}

	kid, err := d.repo.SaveKid(ctx, params)
	if err != nil {
		return output, errors.Wrap(err, "backend.ParentAdminAddKid -> repo.SaveKid")
	}

	output.ID = kid.ID
	output.ParentID = kid.ParentID
	output.NIK = kid.NIK
	output.FullName = kid.FullName
	output.Domisili = kid.Domisili
	output.TanggalLahir = kid.TanggalLahir
	output.JenisKelamin = kid.JenisKelamin
	output.CreatedAt = kid.CreatedAt

	return output, nil
}
