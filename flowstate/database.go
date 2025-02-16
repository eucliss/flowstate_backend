package flowstate

import (
	"log/slog"
	"slices"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var databaseLogger *slog.Logger

var validDBTypes = []string{"test", "prod"}
var validDBNames = []string{
	"flowstate_flows",
	"flowstate_users",
	"flowstate_access",
}

func Database(l *slog.Logger, dbType string, dbName string, migrateStruct ...interface{}) (db *gorm.DB) {
	databaseLogger = l
	if !slices.Contains(validDBTypes, dbType) {
		panic("Invalid database type: " + dbType)
	}
	if !slices.Contains(validDBNames, dbName) {
		panic("Invalid database name: " + dbName)
	}
	switch dbType {
	case "test":
		databaseLogger.Info("Loading test database", "dbName", dbName)
		db = load(dbName+"_test", migrateStruct...)
	case "prod":
		databaseLogger.Info("Loading production database", "dbName", dbName)
		db = load(dbName, migrateStruct...)
	}
	return db
}

func connect() *gorm.DB {
	databaseLogger.Info("Connecting to database")
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func Create(dbName string) bool {
	databaseLogger.Info("Creating database", "dbName", dbName)
	db := connect()
	result := db.Exec("CREATE DATABASE " + dbName)
	return result.Error == nil
}

func Drop(dbName string) bool {
	databaseLogger.Info("Dropping database", "dbName", dbName)
	if strings.Contains(dbName, "test") || strings.Contains(dbName, "dummy") {
		db := connect()
		result := db.Exec("DROP DATABASE " + dbName)
		return result.Error == nil
	} else {
		panic("Cannot drop database " + dbName)
	}
}

func load(dbName string, migrateStruct ...interface{}) *gorm.DB {
	databaseLogger.Info("Loading database", "dbName", dbName)
	dsn := "host=localhost user=postgres password=password dbname=" + dbName + " port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to flowstate database")
	}

	// Auto migrate the schema
	err = db.AutoMigrate(migrateStruct...)
	if err != nil {
		panic("failed to migrate database schema")
	}
	return db
}
