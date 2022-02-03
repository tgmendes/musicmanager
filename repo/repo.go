package repo

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type Store struct {
	DB *pgx.Conn
}

func SubQuery(sb sq.SelectBuilder) sq.Sqlizer {
	sql, args, _ := sb.ToSql()
	return sq.Expr("("+sql+")", args...)
}
