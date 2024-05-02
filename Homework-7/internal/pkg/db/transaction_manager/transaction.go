package transaction_manager

// немного копипаста, но не хочется сливать методы бд и траназкции друг с другом

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Database struct represents a database connection pool
type transaction struct {
	tx pgx.Tx
}

func newTransaction(tx pgx.Tx) *transaction {
	return &transaction{tx: tx}
}

// Get executes a SQL query and stores one result in the provided destination
func (t transaction) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, t.tx, dest, query, args...)
}

// Select executes a SQL query and stores several result in the provided destination
func (t transaction) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, t.tx, dest, query, args...)
}

// Exec executes a SQL command and returns the command tag
func (t transaction) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return t.tx.Exec(ctx, query, args...)
}

// ExecQueryRow executes a SQL command and and expects query can returning result for Scan(),
// do not use if query have no returning statement
func (t transaction) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return t.tx.QueryRow(ctx, query, args...)
}
