package gormdb

import (
	"luis/app/models"

	"github.com/revel/revel"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	DB      *gorm.DB
	gormLog = revel.AppLog
)

func init() {
	revel.RegisterModuleInit(func(m *revel.Module) {
		gormLog = m.Log
	})
}

func InitDB() {
	revel.AppLog.Warn("InitDB")
	// dbUser := revel.Config.StringDefault("db.user", "default")
	// dbPassword := revel.Config.StringDefault("db.password", "")
	dbName := revel.Config.StringDefault("db.name", "default")

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		gormLog.Fatalf("failed to open database: %s", err.Error())
	}

	DB = db

	// Auto-migrate all models.
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		gormLog.Fatalf("migration of %q failed: %s", "User", err.Error())
	}
}
