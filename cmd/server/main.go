package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/padurean/go-sqlx-postgres-with-migration-example/internal/database"
	"github.com/padurean/go-sqlx-postgres-with-migration-example/internal/env"
)

func main() {
	log.Printf("env: %+v\n", env.Global)

	log.Println("connecting to database ...")
	db := database.MustConnect(env.Global.DB.Driver, env.Global.DB.URL)
	database.InitDDL()
	database.InitDML()
	if err := database.Migrate(db); err != nil {
		panic(err)
	}

	log.Println("inserting some ledger stats ...")
	ledgerStats := &database.LedgerStats{
		LedgerID:         11,
		LastInsertionAt:  sql.NullTime{Time: time.Now(), Valid: true},
		ActiveSources:    111,
		NumberOfEntries:  1111,
		TotalSizeInBytes: 1024,
	}
	insertedLedgerStats, err := ledgerStats.Create(db)
	if err != nil {
		panic(fmt.Sprintf("failed to create ledger stats: %+v: %v\n", ledgerStats, err))
	}
	log.Printf("inserted ledger stats: %+v\n", insertedLedgerStats)

	log.Println("updating ledger stats ...")
	insertedLedgerStats.Tampered = sql.NullTime{Time: time.Now(), Valid: true}
	updatedLedgerStats, err := insertedLedgerStats.Update(db)
	if err != nil {
		panic(fmt.Sprintf("failed to update ledger stats: %+v: %v\n", insertedLedgerStats, err))
	}
	log.Printf("updated ledger stats: %+v\n", updatedLedgerStats)

	// print some party emojis :)
	fmt.Println(unquoteCodePoint("\\U0001f389"), unquoteCodePoint("\\U0001f973"))
}

func unquoteCodePoint(s string) string {
	r, _ := strconv.ParseInt(strings.TrimPrefix(s, "\\U"), 16, 32)
	return string(rune(r))
}
