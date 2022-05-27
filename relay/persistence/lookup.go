package persistence

import (
	"context"
	"database/sql"
	"errors"
	"github.com/enclaive/relay/models"
	"github.com/jmoiron/sqlx"
)

type PgLookupRepository struct {
	db *sqlx.DB
}

func (p PgLookupRepository) GetByUsername(ctx context.Context, email string) (l models.Lookup, err error) {
	const query = "SELECT * FROM lookup WHERE username = $1 LIMIT 1"

	err = p.db.GetContext(ctx, &l, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		err = errEmptyResultSet
	}

	return
}

func (p PgLookupRepository) Set(ctx context.Context, l models.Lookup) (err error) {
	const query = "INSERT INTO lookup(username, enclave_address) VALUES($1,$2)"

	_, err = p.db.ExecContext(ctx, query, l.Username, l.EnclaveAddress)

	return
}
