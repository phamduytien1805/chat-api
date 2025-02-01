package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phamduytien1805/package/config"
)

const (
	defaultMaxConns          = int32(4)
	defaultMinConns          = int32(0)
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = time.Minute * 30
	defaultHealthCheckPeriod = time.Minute
	defaultConnectTimeout    = time.Second * 5
)

var PGConn *pgxpool.Pool

func NewPostgresql(config *config.DBConfig) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(config.Source)
	if err != nil {
		return nil, err
	}
	// dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	// 	pgxUUID.Register(conn.TypeMap())
	// 	return nil
	// }
	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	PGConn, err = pgxpool.NewWithConfig(context.TODO(), dbConfig)

	return PGConn, err
}
