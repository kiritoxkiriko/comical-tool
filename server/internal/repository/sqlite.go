package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	db     *sqlx.DB
	driver string
}

func Open(driver, dsn string) (*Store, error) {
	name, err := normalizeDriver(driver)
	if err != nil {
		return nil, err
	}
	sqlx.BindDriver("pgx", sqlx.DOLLAR)
	db, err := sqlx.Open(name, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db, driver: name}, nil
}

func OpenSQLite(dsn string) (*Store, error) {
	return Open("sqlite", dsn)
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Migrate(ctx context.Context) error {
	switch s.driver {
	case "sqlite":
		return s.execStatements(ctx, sqliteSchema)
	case "pgx":
		return s.execStatements(ctx, postgresSchema)
	case "mysql":
		return s.execStatements(ctx, mysqlSchema)
	default:
		return fmt.Errorf("unsupported database driver %q", s.driver)
	}
}

func (s *Store) exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, s.db.Rebind(query), args...)
}

func (s *Store) get(ctx context.Context, dest any, query string, args ...any) error {
	return s.db.GetContext(ctx, dest, s.db.Rebind(query), args...)
}

func (s *Store) selectRows(ctx context.Context, dest any, query string, args ...any) error {
	return s.db.SelectContext(ctx, dest, s.db.Rebind(query), args...)
}

func (s *Store) timeArg(t dbTime) any {
	if !t.Valid {
		return nil
	}
	if s.driver == "sqlite" {
		return t.Time.UTC().Format(timeLayout)
	}
	return t.Time.UTC()
}

func (s *Store) nowArg() any {
	return s.timeArg(newDBTime(now()))
}

func (s *Store) execStatements(ctx context.Context, schema string) error {
	for _, statement := range strings.Split(schema, ";") {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}
		if _, err := s.db.ExecContext(ctx, statement); err != nil && !isIgnorableMigrationError(err) {
			return err
		}
	}
	return nil
}

func normalizeDriver(driver string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(driver)) {
	case "", "sqlite", "sqlite3":
		return "sqlite", nil
	case "postgres", "postgresql", "pgx":
		return "pgx", nil
	case "mysql":
		return "mysql", nil
	default:
		return "", fmt.Errorf("unsupported database driver %q", driver)
	}
}

func isIgnorableMigrationError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1061 {
		return true
	}
	return false
}
