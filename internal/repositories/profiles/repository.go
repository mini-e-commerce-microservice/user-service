package profiles

import (
	"github.com/Masterminds/squirrel"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
)

type repository struct {
	rdbms wsqlx.Rdbms
	sq    squirrel.StatementBuilderType
}

func NewRepository(rdbms wsqlx.Rdbms) *repository {
	return &repository{
		rdbms: rdbms,
		sq:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
