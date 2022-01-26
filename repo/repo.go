package repo

import (
	"github.com/jackc/pgx/v4"
)

type Store struct {
	DB *pgx.Conn
}
