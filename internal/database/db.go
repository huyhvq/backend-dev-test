package database

import (
	"context"
	"database/sql"
	"github.com/huyhvq/backend-dev-test/assets"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun/migrate"
)

const defaultTimeout = 3 * time.Second

var migrations = migrate.NewMigrations()

type DB struct {
	*bun.DB
}

func New(dsn string, autoMigrate bool) (*DB, error) {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqlDB, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if autoMigrate {
		if err := migrations.Discover(assets.EmbeddedFiles); err != nil {
			return nil, err
		}
		migrator := migrate.NewMigrator(db, migrations)
		if err := migrator.Init(context.Background()); err != nil {
			return nil, err
		}
		if _, err := migrator.Migrate(context.Background()); err != nil {
			return nil, err
		}
	}
	return &DB{db}, nil
}
