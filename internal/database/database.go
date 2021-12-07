package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	// initialize database (PosgreSQL) driver
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

// DB ...
type DB struct {
	*sqlx.DB
}

// CRUDFields ...
type CRUDFields struct {
	ID      uint64       `json:"-"`
	Created time.Time    `json:"-"`
	Updated time.Time    `json:"-"`
	Deleted sql.NullTime `json:"-"`
}

// LimitAndOffset ...
type LimitAndOffset struct {
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// MustConnect ...
func MustConnect(driver string, url string) *DB {
	db := DB{DB: sqlx.MustConnect(driver, url)}
	db.SetConnMaxLifetime(60 * time.Second)
	db.SetConnMaxIdleTime(15 * time.Second)
	return &db
}

// Upsert ...
func Upsert(db *DB, sqlUpsert string, sqlSelectByID string, argUpsert interface{}, dest interface{}) error {
	stmtUpsert, err := db.PrepareNamed(sqlUpsert)
	if err != nil {
		return fmt.Errorf("error preparing named db upsert: %v", err)
	}
	defer stmtUpsert.Close()
	var id int64
	if err := stmtUpsert.Get(&id, argUpsert); err != nil {
		return fmt.Errorf("error executing db upsert: %v", err)
	}
	stmtSelectByID, err := db.Preparex(sqlSelectByID)
	if err != nil {
		return fmt.Errorf("error preparing db select by ID: %v", err)
	}
	defer stmtSelectByID.Close()
	if err := stmtSelectByID.Get(dest, id); err != nil {
		return fmt.Errorf("error executing db select by ID: %v", err)
	}
	return nil
}

// UpsertBatch ...
func UpsertBatch(db *DB, sqlUpsert string, argUpsert []interface{}, dest []int64) error {
	if len(argUpsert) == 0 || len(argUpsert) != len(dest) {
		return errors.New("invalid arguments to upsert")
	}
	stmtUpsert, err := db.PrepareNamed(sqlUpsert)
	if err != nil {
		return fmt.Errorf("error preparing named db upsert: %v", err)
	}
	defer stmtUpsert.Close()
	for i := range argUpsert {
		var id int64
		if err := stmtUpsert.Get(&id, argUpsert[i]); err != nil {
			return fmt.Errorf("error executing db upsert: %v", err)
		}
		dest = append(dest, id)
	}
	return nil
}

// UpsertTx ...
func UpsertTx(tx *sqlx.Tx, sqlUpsert string, sqlSelectByID string, argUpsert interface{}, dest interface{}) error {
	stmtUpsert, err := tx.PrepareNamed(sqlUpsert)
	if err != nil {
		return fmt.Errorf("error preparing named db upsert: %v", err)
	}
	defer stmtUpsert.Close()
	var id int64
	if err := stmtUpsert.Get(&id, argUpsert); err != nil {
		return fmt.Errorf("error executing db upsert: %v", err)
	}
	stmtSelectByID, err := tx.Preparex(sqlSelectByID)
	if err != nil {
		return fmt.Errorf("error preparing db select by ID: %v", err)
	}
	defer stmtSelectByID.Close()
	if err := stmtSelectByID.Get(dest, id); err != nil {
		return fmt.Errorf("error executing db select by ID: %v", err)
	}
	return nil
}

// SelectOne ...
func SelectOne(db *DB, sqlSelect string, argSelect, dest interface{}) error {
	stmtSelect, err := db.Preparex(sqlSelect)
	if err != nil {
		return fmt.Errorf("error preparing db select one: %v", err)
	}
	err = stmtSelect.Get(dest, argSelect)
	stmtSelect.Close()
	return err
}

// UpsertMany ...
func UpsertMany(db *DB, sqlUpsert string) (int64, error) {
	res, err := db.Exec(sqlUpsert)
	if err != nil {
		return 0, fmt.Errorf("error executing db multi insert: %v", err)
	}
	return res.RowsAffected()
}

// DeleteMany ...
func DeleteMany(db *DB, sqlDelete string, args ...interface{}) (int64, error) {
	res, err := db.Exec(sqlDelete, args...)
	if err != nil {
		return 0, fmt.Errorf("error deleting many: %w", err)
	}
	return res.RowsAffected()
}

// UpdateOne will update a single ts column
func UpdateOne(db *DB, sqlStmt, opType string, id uint64) error {
	stmtUpdate, err := db.Preparex(sqlStmt)
	if err != nil {
		return fmt.Errorf("error preparing db stmt for %s ID %d: %v", opType, id, err)
	}
	defer stmtUpdate.Close()
	result, err := stmtUpdate.Exec(id)
	if err != nil {
		return fmt.Errorf("error executing db stmt for %s ID %d: %v", opType, id, err)
	}
	nbUpdated, err := result.RowsAffected()
	if nbUpdated == 0 {
		return fmt.Errorf(
			"error %s: ID %d - %w", opType, id, sql.ErrNoRows)
	} else if nbUpdated > 1 {
		return fmt.Errorf(
			"error %s: ID %d: %d entities have been affected instead of exactly 1",
			opType, id, nbUpdated)
	}
	return nil
}
