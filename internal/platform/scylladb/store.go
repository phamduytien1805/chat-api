package scylladb

import (
	"context"
	"log/slog"

	"github.com/phamduytien1805/internal/platform/scylladb/cql"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/migrate"
)

type Store interface {
}
type CQLStore struct {
	session gocqlx.Session
}

func NewCQLStore(session gocqlx.Session) Store {
	return &CQLStore{session: session}
}

func RunMigration(logger *slog.Logger, session gocqlx.Session) error {

	// Add callback prints
	// log := func(ctx context.Context, session gocqlx.Session, ev migrate.CallbackEvent, name string) error {
	// 	logger.Info(name, "event", ev)
	// 	return nil
	// }
	// reg := migrate.CallbackRegister{}
	// reg.Add(migrate.BeforeMigration, "m1.cql", log)
	// reg.Add(migrate.AfterMigration, "m1.cql", log)
	// reg.Add(migrate.CallComment, "1", log)
	// reg.Add(migrate.CallComment, "2", log)
	// reg.Add(migrate.CallComment, "3", log)
	// migrate.Callback = reg.Callback
	pending, err := migrate.Pending(context.Background(), session, cql.Files)
	if err != nil {
		logger.Error("Pending Error", "detail", err.Error())
		return err
	}
	logger.Info("Pending", "amount", len(pending))
	if err := migrate.FromFS(context.Background(), session, cql.Files); err != nil {
		logger.Error("Migrate Error", "detail", err.Error())
		return err
	}
	return nil
}
