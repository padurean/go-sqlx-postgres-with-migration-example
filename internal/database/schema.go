package database

import (
	"fmt"
	"log"

	"github.com/padurean/go-sqlx-postgres-with-migration-example/internal/env"
)

var migrations []*Migration

// InitDDL ...
func InitDDL() {
	migrations = append(migrations, getLedgerStatsMigrations()...)
}

// InitDML ...
func InitDML() {
	initMigrationDML()

	initLedgerStatsDML()
}

// Migrate ...
func Migrate(db *DB) error {
	log.Println("migrating database ...")

	// 1. create schema and migrations table
	createSchemaSQL :=
		"CREATE SCHEMA IF NOT EXISTS " + env.Global.DB.Schema + " AUTHORIZATION " + env.Global.DB.User + ";"
	if _, err := db.Exec(createSchemaSQL); err != nil {
		return err
	}
	if _, err := db.Exec(getMigrationDDL()); err != nil {
		return err
	}

	// 2. load any previously run migrations and group them by entity
	oldMigrations, err := (&Migration{}).List(db)
	if err != nil {
		return fmt.Errorf("error loading past migrations from db: %v", err)
	}
	oldMigrationsPerEntity := make(map[string][]*Migration)
	for _, om := range oldMigrations {
		oldMigrationsPerEntity[om.Entity] =
			append(oldMigrationsPerEntity[om.Entity], om)
	}

	// 3. loop through all migrations, executing any new ones and
	//    recording them in the migrations table
outer:
	for _, m := range migrations {
		oms := oldMigrationsPerEntity[m.Entity]
		for _, om := range oms {
			if om.Version == m.Version {
				continue outer
			}
		}
		log.Printf("executing %s migration %d ...\n", m.Entity, m.Version)
		if _, err := db.Exec(m.SQL); err != nil {
			return fmt.Errorf("error executing database migration %+v: %v", m, err)
		}
		_, err := m.Create(db)
		if err != nil {
			return fmt.Errorf(
				"migration %+v executed successfully, but there was an error "+
					"recording it in the migrations table: %v",
				m, err)
		}
	}

	return nil
}
