package repo

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dropDB   = `DROP DATABASE test_devices;`
	createDB = `CREATE DATABASE test_devices;`
)

func TestMain(m *testing.M) {
	psql, err := sql.Open("postgres", "host=localhost port=5432 user=test password=test sslmode=disable")
	if err != nil {
		panic(fmt.Errorf("sql.Open() err: %v", err))
	}

	defer psql.Close()

	_, err = psql.Exec(createDB)
	if err != nil {
		panic(fmt.Errorf("sql Create DB err: %v", err))
	}
	_, err = psql.Exec(dropDB)
	if err != nil {
		panic(fmt.Errorf("sql Drop DB err: %v", err))
	}

	db, err := sql.Open("postgres", "host=localhost port=5432 user=test password=test sslmode=disable")
	if err != nil {
		panic(fmt.Errorf("sql.Open() err: %v", err))
	}

	defer db.Close()

	m.Run()
}
