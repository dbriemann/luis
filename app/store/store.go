package store

import (
	"io/ioutil"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/revel/revel"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dblog = revel.AppLog
	DB    *sqlx.DB
)

func init() {
	revel.RegisterModuleInit(func(m *revel.Module) {
		dblog = m.Log
	})
}

func InitDB() {
	dbName := revel.Config.StringDefault("db.name", "db.sqlite")
	db, err := sqlx.Open("sqlite3", dbName)
	if err != nil {
		dblog.Fatalf("failed to open database: %q", err.Error())
	}

	DB = db

	// Enable foreign keys.
	if _, err := DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		dblog.Fatalf("could not enable foreign keys: %q", err.Error())
	}

	// Create schemas ( TODO: migrate ).
	schemaFile := filepath.Join(revel.BasePath, "/sql/", "v0.sql")
	b, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		dblog.Fatalf("could not open sql schemas v0: %q", err.Error())
	}
	_, err = DB.Exec(string(b))
	if err != nil {
		dblog.Fatalf("could not exec schemas v0: %q", err.Error())
	}
}
