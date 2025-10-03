// ...existing code...
package postgressql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nocson47/beaconofknowledge/config"
)

func ConnectPGX(cfg *config.Configuration) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	// Optional: ping
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
