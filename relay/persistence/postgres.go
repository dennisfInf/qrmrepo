package persistence

import (
	"errors"
	"fmt"
	"github.com/enclaive/relay/config"
	"github.com/enclaive/relay/server"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var errEmptyResultSet = errors.New("result set was empty")

type PgRepository struct {
	db *sqlx.DB
}

func New(conf config.PostgresConfig) *PgRepository {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		conf.Host, conf.User, conf.Password, conf.DBName, conf.Port)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("failed to connect to database")
	}

	return &PgRepository{db: db}
}

func (r *PgRepository) ApplySchema() error {
	for _, schema := range schemas {
		_, err := r.db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PgRepository) Lookup() server.LookupRepository {
	return &PgLookupRepository{
		db: r.db,
	}
}

func (r *PgRepository) IsEmptyResultSetError(err error) bool {
	return errors.Is(err, errEmptyResultSet)
}
