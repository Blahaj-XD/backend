package repo

import "github.com/jackc/pgx/v5/pgxpool"

type Dependency struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Dependency {
	return &Dependency{db: db}
}
