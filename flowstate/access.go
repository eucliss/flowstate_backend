package flowstate

import (
	"log/slog"

	"gorm.io/gorm"
)

type Access struct {
	gorm.Model
	UserID uint   `json:"user_id" form:"user_id"`
	FlowID uint   `json:"flow_id" form:"flow_id"`
	Access string `json:"access" form:"access"`
}

var accessLogger *slog.Logger
var AccessDB *gorm.DB

func AccessDatabase(l *slog.Logger, dbType string) *gorm.DB {
	accessLogger = l
	accessLogger.Info("Loading access database", "dbType", dbType)
	AccessDB = Database(l, dbType, "flowstate_access", &Access{})
	return AccessDB
}

func (access *Access) Create() (int64, error) {
	result := AccessDB.Create(&access)
	if result.Error != nil {
		accessLogger.Error("Failed to create access", "error", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (access *Access) Update() (int64, error) {
	result := AccessDB.Save(&access)
	if result.Error != nil {
		accessLogger.Error("Failed to update access", "error", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (access *Access) Get() *Access {
	result := AccessDB.First(&access)
	if result.Error != nil {
		accessLogger.Error("Failed to get access", "error", result.Error)
		return nil
	}
	return access
}

func (access *Access) Delete() {
	AccessDB.Delete(&access)
}

func (access *Access) Exists() bool {
	return AccessDB.Where("id = ?", access.ID).First(&access).RowsAffected > 0
}
