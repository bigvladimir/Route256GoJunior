//go:build integration

package postgresql

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"homework/internal/pkg/db"
)

type dbOps interface {
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetPool(_ context.Context) *pgxpool.Pool
}

type TDB struct {
	DB dbOps
}

func NewTDB() *TDB {
	db, err := db.NewDb(context.Background(), "../docker-compose-test.yaml")
	if err != nil {
		panic(err)
	}
	return &TDB{DB: db}
}

func (d *TDB) Close() {
	d.DB.GetPool(context.Background()).Close()
}

func (d *TDB) SetUp(t *testing.T, tableName ...string) {
	t.Helper()
	d.truncateTable(context.Background(), tableName...)
}

func (d *TDB) TearDown(t *testing.T) {
	t.Helper()
	// шаблон для расширения функционала
}

func (d *TDB) truncateTable(ctx context.Context, tableName ...string) {
	q := fmt.Sprintf("TRUNCATE table %s RESTART IDENTITY CASCADE", strings.Join(tableName, ","))
	if _, err := d.DB.Exec(ctx, q); err != nil {
		panic(err)
	}
}
