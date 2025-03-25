package postgresql

import (
	"database/sql"
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage"
)

type postgresqlRepository struct {
	db *sql.DB
}

func NewPostgresqlRepository(db *sql.DB) storage.Storage {
	return &postgresqlRepository{
		db: db,
	}
}
