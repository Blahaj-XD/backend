package repo

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Parent struct {
	ID           int
	NIK          string
	Username     string
	Email        string
	Password     string `json:"-"`
	PhoneNumber  string
	FullName     string
	Domisili     string
	TanggalLahir string
	JenisKelamin int
	Alamat       string
	RtRW         string
	Kelurahan    string
	Kecamatan    string
	Pekerjaan    string
	CreatedAt    time.Time
}

func (d *Dependency) SaveParent(ctx context.Context, params Parent) (int, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cols := []string{
		"nik",
		"username",
		"email",
		"password",
		"phone_number",
		"full_name",
		"domisili",
		"tanggal_lahir",
		"jenis_kelamin",
		"alamat",
		"rt_rw",
		"kelurahan",
		"kecamatan",
		"pekerjaan",
		"created_at"}

	query := qb.Insert("parents").
		Columns(cols...).
		Values(
			params.NIK, params.Username, params.Email, params.Password, params.PhoneNumber, params.FullName, params.Domisili, params.TanggalLahir,
			params.JenisKelamin, params.Alamat, params.RtRW, params.Kelurahan,
			params.Kecamatan, params.Pekerjaan, params.CreatedAt).
		Suffix("RETURNING \"id\"")

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	var id int

	if err := d.db.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (d *Dependency) FindParent(ctx context.Context, col string, value any) (Parent, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	cols := []string{
		"id",
		"nik",
		"username",
		"email",
		"password",
		"phone_number",
		"full_name",
		"domisili",
		"tanggal_lahir",
		"jenis_kelamin",
		"alamat",
		"rt_rw",
		"kelurahan",
		"kecamatan",
		"pekerjaan",
		"created_at"}

	query := qb.Select(cols...).
		From("parents").
		Where(sq.Eq{col: value})

	sql, args, err := query.ToSql()
	if err != nil {
		return Parent{}, err
	}

	var parent Parent

	if err := d.db.QueryRow(ctx, sql, args...).Scan(
		&parent.ID, &parent.NIK, &parent.Username, &parent.Email, &parent.Password, &parent.PhoneNumber,
		&parent.FullName, &parent.Domisili, &parent.TanggalLahir, &parent.JenisKelamin,
		&parent.Alamat, &parent.RtRW, &parent.Kelurahan, &parent.Kecamatan,
		&parent.Pekerjaan, &parent.CreatedAt); err != nil {
		return Parent{}, err
	}

	return parent, nil
}