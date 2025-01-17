//go:build integration

package database

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	gomigrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// required for gomigrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"

	"github.com/gojek/mlp/api/log"
)

func connectionString(db string) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password='%s' sslmode=disable",
		host, 5432, user, db, password)
}

func create(conn *sql.DB, dbName string) (*sql.DB, error) {
	if _, err := conn.Exec("CREATE DATABASE " + dbName); err != nil {
		return nil, err
	} else if testDb, err := sql.Open("postgres", connectionString(dbName)); err != nil {
		if _, err := conn.Exec("DROP DATABASE " + dbName); err != nil {
			log.Fatalf("Failed to cleanup integration test database: \n%s", err)
		}
		return nil, err
	} else {
		return testDb, nil
	}
}

func migrate(db *sql.DB, dbName string) (*sql.DB, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	defer driver.Close()

	if migrations, err := gomigrate.NewWithDatabaseInstance("file://../../db-migrations", dbName, driver); err != nil {
		return nil, err
	} else if err = migrations.Up(); err != nil {
		return nil, err
	}
	return sql.Open("postgres", connectionString(dbName))
}

// Connects to test postgreSQL instance (either local or the one at CI environment)
// and creates a new database with an up-to-date schema
func CreateTestDatabase() (*gorm.DB, func(), error) {
	testDbName := fmt.Sprintf("mlp_id_%d", time.Now().UnixNano())

	connStr := connectionString(database)
	log.Infof("connecting to test db: %s", connStr)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, err
	}

	testDb, err := create(conn, testDbName)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := testDb.Close(); err != nil {
			log.Fatalf("Failed to close connection to integration test database: \n%s", err)
		} else if _, err := conn.Exec("DROP DATABASE " + testDbName); err != nil {
			log.Fatalf("Failed to cleanup integration test database: \n%s", err)
		} else if err = conn.Close(); err != nil {
			log.Fatalf("Failed to close database: \n%s", err)
		}
	}

	if testDb, err = migrate(testDb, testDbName); err != nil {
		log.Errorf("Error %v", err)

		cleanup()
		return nil, nil, err
	} else if gormDb, err := gorm.Open("postgres", testDb); err != nil {
		cleanup()
		return nil, nil, err
	} else {
		return gormDb, cleanup, nil
	}
}

func WithTestDatabase(t *testing.T, test func(t *testing.T, db *gorm.DB)) {
	if testDb, cleanupFn, err := CreateTestDatabase(); err != nil {
		t.Fatalf("Fail to create an integration test database: \n%s", err)
	} else {
		test(t, testDb)
		cleanupFn()
	}
}
