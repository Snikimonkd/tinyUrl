package tinyUrl

import (
	"database/sql"
	"embed"
	"errors"
	"os"
	"time"

	_ "embed"

	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var f embed.FS

func getPostgres() *sql.DB {
	// подождать пока база поднимется
	time.Sleep(5 * time.Second)

	dsn := os.Getenv("DB")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic("cant parse config" + err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic("can`t ping db" + err.Error())
	}

	db.SetMaxOpenConns(10)

	source, err := iofs.New(f, "migrations")
	if err != nil {
		panic("can`t open migrations source:" + err.Error())
	}

	m, err := migrate.NewWithSourceInstance(
		"embed",
		source,
		dsn)
	if err != nil {
		panic("can`t init migrations:" + err.Error())
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}

	return db
}

func Init() *sqlx.DB {
	return sqlx.NewDb(getPostgres(), "psx")
}
