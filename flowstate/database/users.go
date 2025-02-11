package flowstate

import (
	"log/slog"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

var userLogger *slog.Logger
var UserDB *gorm.DB

func UsersDatabase(l *slog.Logger, dbType string) *gorm.DB {
	userLogger = l
	userLogger.Info("Loading users database", "dbType", dbType)
	UserDB = Database(l, dbType, "flowstate_users", &User{})
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

func (user *User) Update() (int64, error) {
	result := UserDB.Save(&user)
	if result.Error != nil {
		userLogger.Error("Failed to update user", "error", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
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

func (user *User) LoginSuccess() bool {
	res := user.Get()
	if res == nil {
		return false
	}
	return res.Password == user.Password
}
