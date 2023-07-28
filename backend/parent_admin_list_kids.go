package backend

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type ParentAdminListKidsItem struct {
	ID            int       `json:"id"`
	ParentID      int       `json:"parent_id"`
	AccountNumber string    `json:"account_number"`
	NIK           string    `json:"nik"`
	FullName      string    `json:"full_name"`
	Domisili      string    `json:"domisili"`
	TanggalLahir  string    `json:"tanggal_lahir"`
	JenisKelamin  int       `json:"jenis_kelamin"`
	CreatedAt     time.Time `json:"created_at"`
}

type ParentAdminListKidsOutput struct {
	TotalItems int                       `json:"total_items"`
	Items      []ParentAdminListKidsItem `json:"items"`
}

func (d *Dependency) ParentAdminListKids(ctx context.Context, parentID int) (ParentAdminListKidsOutput, error) {
	var output ParentAdminListKidsOutput

	kids, err := d.repo.ListParentKids(ctx, parentID)
	if err != nil {
		return ParentAdminListKidsOutput{}, errors.Wrap(err, "backend.ParentAdminListKids -> repo.ListParentKids")
	}

	output.TotalItems = len(kids)
	output.Items = make([]ParentAdminListKidsItem, 0)
	for _, kid := range kids {
		var item ParentAdminListKidsItem
		item.ID = kid.ID
		item.ParentID = kid.ParentID
		item.AccountNumber = kid.AccountNumber
		item.NIK = kid.NIK
		item.FullName = kid.FullName
		item.Domisili = kid.Domisili
		item.TanggalLahir = kid.TanggalLahir
		item.JenisKelamin = kid.JenisKelamin
		item.CreatedAt = kid.CreatedAt

		output.Items = append(output.Items, item)
	}

	return output, nil
}
