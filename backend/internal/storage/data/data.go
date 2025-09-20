package data

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type UserCache interface {
	Exists(ctx context.Context, dataType, data string) (bool, error)
}

type dataStorage struct {
	db        *sqlx.DB
	userCache UserCache
}

func NewDataStorage(db *sqlx.DB, usercCache UserCache) *dataStorage {
	return &dataStorage{
		db:        db,
		userCache: usercCache,
	}
}

func (s *dataStorage) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, nil)
}
