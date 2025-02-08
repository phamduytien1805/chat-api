package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phamduytien1805/user/domain"
)

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	q        *Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) domain.UserRepo {
	return &SQLStore{
		connPool: connPool,
		q:        New(connPool),
	}
}
