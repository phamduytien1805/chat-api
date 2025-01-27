package scylladb

import (
	"github.com/scylladb/gocqlx/v3"
)

type Store interface {
}
type CQLStore struct {
	session gocqlx.Session
}

func NewCQLStore(session gocqlx.Session) Store {
	return &CQLStore{session: session}
}
