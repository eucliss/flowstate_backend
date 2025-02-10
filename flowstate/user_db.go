package flowstate

import (
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

var userLogger *slog.Logger
var UserDB *gorm.DB
var dbName string

func InitUserDB(l *slog.Logger, dbType string) {
	if dbType == "test" {
		dbName = "flowstate_users_test"
	} else if dbType == "prod" {
		dbName = "flowstate_users"
	} else {
		panic("Invalid database type: " + dbType)
	}
	userLogger = l
	UserDB = LoadUserDB(l, dbName)
}

func startGormDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func CreateDB(dbName string) bool {
	db := startGormDB()
	result := db.Exec("CREATE DATABASE " + dbName)
	return result.Error == nil
}

func DropDB(dbName string) bool {
	if dbName == "flowstate_users_test" || dbName == "dummy_db" {
		db := startGormDB()
		result := db.Exec("DROP DATABASE " + dbName)
		return result.Error == nil
	} else {
		panic("Cannot drop database " + dbName)
	}
}

func LoadUserDB(l *slog.Logger, dbName string) *gorm.DB {

	userLogger = l

	dsn := "host=localhost user=postgres password=password dbname=" + dbName + " port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to flowstate database")
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to migrate database schema")
	}

	UserDB = db
	return UserDB
}

func (user *User) Create() (int64, error) {
	result := UserDB.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		userLogger.Error("Failed to create user", "error", result.Error)
		return 0, result.Error
	} else {
		userLogger.Info("User created", "id", user.Model.ID, "rows affected", result.RowsAffected)
		return result.RowsAffected, nil
	}
}

func (user *User) Get() *User {
	result := UserDB.First(&user, "username = ?", user.Username)
	if result.Error != nil {
		userLogger.Error("Failed to get user", "error", result.Error)
		return nil
	}
	if user.ID == 0 {
		userLogger.Error("User not found", "username", user.Username)
		return nil
	}
	return user
}

func (user *User) Delete() {
	UserDB.Delete(&user)
	userLogger.Info("User deleted", "id", user.Model.ID)
}

func (user *User) Exists() bool {
	return UserDB.First(&user, "username = ?", user.Username).RowsAffected > 0
}
