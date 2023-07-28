package repo

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type Kid struct {
	ID            int
	ParentID      int
	AccountNumber string
	NIK           string
	FullName      string
	Domisili      string
	TanggalLahir  string
	JenisKelamin  int
	CreatedAt     time.Time
}

func (d *Dependency) SaveKid(ctx context.Context, params Kid) (Kid, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cols := []string{
		"parent_id",
		"account_number",
		"nik",
		"full_name",
		"domisili",
		"tanggal_lahir",
		"jenis_kelamin",
		"created_at"}

	query := qb.Insert("kids").
		Columns(cols...).
		Values(
			params.ParentID, params.AccountNumber, params.NIK,
			params.FullName, params.Domisili, params.TanggalLahir,
			params.JenisKelamin, params.CreatedAt).
		Suffix("RETURNING \"id\"")

	sql, args, err := query.ToSql()
	if err != nil {
		return Kid{}, errors.Wrap(err, "repo.SaveKid")
	}

	var id int
	err = d.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return Kid{}, errors.Wrap(err, "repo.SaveKid")
	}

	var output Kid
	output.ID = id
	output.ParentID = params.ParentID
	output.AccountNumber = params.AccountNumber
	output.NIK = params.NIK
	output.FullName = params.FullName
	output.Domisili = params.Domisili
	output.TanggalLahir = params.TanggalLahir
	output.JenisKelamin = params.JenisKelamin
	output.CreatedAt = params.CreatedAt

	return output, nil
}

func (d *Dependency) FindKid(ctx context.Context, col string, value any) (Kid, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := qb.Select("*").
		From("kids").
		Where(sq.Eq{col: value})

	sql, args, err := query.ToSql()
	if err != nil {
		return Kid{}, errors.Wrap(err, "repo.FindKid")
	}

	var output Kid
	err = d.db.QueryRow(ctx, sql, args...).
		Scan(
			&output.ID,
			&output.ParentID,
			&output.AccountNumber,
			&output.NIK,
			&output.FullName,
			&output.Domisili,
			&output.TanggalLahir,
			&output.JenisKelamin,
			&output.CreatedAt,
		)
	if err != nil {
		return Kid{}, errors.Wrap(err, "repo.FindKid")
	}

	return output, nil
}

func (d *Dependency) ListParentKids(ctx context.Context, parentID int) ([]Kid, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := qb.Select("*").
		From("kids").
		Where(sq.Eq{"parent_id": parentID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := d.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var output []Kid
	for rows.Next() {
		var kid Kid
		err = rows.Scan(
			&kid.ID,
			&kid.ParentID,
			&kid.AccountNumber,
			&kid.NIK,
			&kid.FullName,
			&kid.Domisili,
			&kid.TanggalLahir,
			&kid.JenisKelamin,
			&kid.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		output = append(output, kid)
	}

	return output, nil
}
