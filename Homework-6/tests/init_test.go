//go:build integration

package tests

import (
	"io"
	"log"
	"os"
	"testing"

	"homework/tests/postgresql"
	"homework/tests/testserver"
)

var (
	db   *postgresql.TDB
	serv *testserver.TServer
)

func TestMain(m *testing.M) {
	db = postgresql.NewTDB()
	defer db.Close()
	serv = testserver.NewTServer(db)

	log.SetOutput(io.Discard)

	exitCode := m.Run()
	os.Exit(exitCode)
}
