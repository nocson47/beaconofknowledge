package postgressql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/nocson47/invoker_board/config"
)

func ConnectPGX(cfg *config.Configuration) (*pgx.Conn, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	return pgx.Connect(context.Background(), connStr)
}
