package psql

import (
	"log/slog"

	"github.com/SHSanderland/EffMobTest/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func InitDB(log *slog.Logger, dsn string) *Storage {

	return &Storage{log: log}
}

func (s Storage) CreateSubscription(sub *model.Subscription) error {

	return nil
}
