//go:build integration

package postgresql

import (
	"context"
	"fmt"
	"homework/internal/pkg/db"
	"strings"
	"testing"
)

const bdCfg = "../docker-compose-test.yaml"

type TDB struct {
	DB db.DBops
}

func NewTDB() *TDB {
	db, err := db.NewDb(context.Background(), bdCfg)
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
