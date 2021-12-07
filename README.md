# go-sqlx-postgres-with-migration-example

Go example showing how to use [sqlx](https://github.com/jmoiron/sqlx) with PostgreSQL and a very simple way to do database migrations

## Prerequisites

- a running PostgreSQL server
- values in the .env file must match the ones PostgreSQL ones:
```console
DB_DRIVER=pgx
DB_HOST_NAME=localhost
DB_USER=some_user
DB_PASSWORD=some_password
DB_NAME=some_database
DB_SCHEMA=some_schema
```

## Run

```console
go run ./cmd/server
```

The output should look like this:
```console
➤ go run ./cmd/server
2021/12/07 18:08:59 env: {DB:{Driver:pgx HostName:localhost User:some_user Password:some_password Name:some_database Schema:some_schema URL:postgres://some_user:some_password@localhost:5432/some_database}}
2021/12/07 18:08:59 connecting to database ...
2021/12/07 18:08:59 migrating database ...
2021/12/07 18:08:59 executing ledger_stats migration 1 ...
2021/12/07 18:08:59 executing ledger_stats migration 2 ...
2021/12/07 18:08:59 executing ledger_stats migration 3 ...
2021/12/07 18:08:59 inserting some ledger stats ...
2021/12/07 18:08:59 inserted ledger stats: &{ID:1 Created:2021-12-07 18:08:59.693316 +0200 EET Updated:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false} LedgerID:11 LastInsertionAt:{Time:2021-12-07 18:08:59.690109 +0200 EET Valid:true} ActiveSources:111 NumberOfEntries:1111 TotalSizeInBytes:1024 Tampered:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}}
2021/12/07 18:08:59 updating ledger stats ...
2021/12/07 18:08:59 updated ledger stats: &{ID:1 Created:2021-12-07 18:08:59.693316 +0200 EET Updated:{Time:2021-12-07 18:08:59.703779 +0200 EET Valid:true} LedgerID:11 LastInsertionAt:{Time:2021-12-07 18:08:59.690109 +0200 EET Valid:true} ActiveSources:111 NumberOfEntries:1111 TotalSizeInBytes:1024 Tampered:{Time:2021-12-07 18:08:59.702133 +0200 EET Valid:true}}
🎉 🥳
```
