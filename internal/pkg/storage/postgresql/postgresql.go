package postgresql

import (
	"database/sql"
	"ozon-test-project/internal/pkg/storage"
)

type postgresqlRepository struct {
	db *sql.DB
}

func NewPostgresqlRepository(db *sql.DB) storage.Storage {
	return &postgresqlRepository{
		db: db,
	}
}
