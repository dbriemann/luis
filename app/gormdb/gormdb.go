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
	dbName := revel.Config.StringDefault("db.name", "luis.sqlite")

	revel.AppLog.Info("init database %q", dbName)

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		gormLog.Fatalf("failed to open database: %s", err.Error())
	}

	DB = db

	// Auto-migrate all models.
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		gormLog.Fatalf("migration of %q failed: %s", "User", err.Error())
	}
	if err := DB.AutoMigrate(&models.File{}); err != nil {
		gormLog.Fatalf("migration of %q failed: %s", "File", err.Error())
	}
	if err := DB.AutoMigrate(&models.Collection{}); err != nil {
		gormLog.Fatalf("migration of %q failed: %s", "Collection", err.Error())
	}
	if err := DB.AutoMigrate(&models.Star{}); err != nil {
		gormLog.Fatalf("migration of %q failed: %s", "Star", err.Error())
	}
	if err := DB.AutoMigrate(&models.Comment{}); err != nil {
		gormLog.Fatalf("migration of %q failed: %s", "Comment", err.Error())
	}
	// if err := DB.AutoMigrate(&models.Tag{}); err != nil {
	// 	gormLog.Fatalf("migration of %q failed: %s", "Tag", err.Error())
	// }
}
