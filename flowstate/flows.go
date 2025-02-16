package flowstate

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"log/slog"

	"encoding/json"

	"gorm.io/gorm"
)

type Flow struct {
	gorm.Model
	Name    string          `json:"flow_name" form:"flow_name"`
	Owner   string          `json:"owner" form:"owner"`
	Content json.RawMessage `gorm:"type:jsonb"`
}

var flowsLogger *slog.Logger
var FlowsDB *gorm.DB

func FlowsDatabase(l *slog.Logger, dbType string) *gorm.DB {
	flowsLogger = l
	flowsLogger.Info("Loading flows database", "dbType", dbType)
	FlowsDB = Database(l, dbType, "flowstate_flows", &Flow{})
	return FlowsDB
}

func (flow *Flow) Create() (int64, error) {
	result := FlowsDB.Create(&flow)
	if result.Error != nil {
		flowsLogger.Error("Failed to create flow", "error", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (flow *Flow) Update() (int64, error) {
	result := FlowsDB.Save(&flow)
	if result.Error != nil {
		flowsLogger.Error("Failed to update flow", "error", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (flow *Flow) Get() *Flow {
	result := FlowsDB.First(&flow)
	if result.Error != nil {
		flowsLogger.Error("Failed to get flow", "error", result.Error)
		return nil
	}
	return flow
}

func (flow *Flow) Delete() {
	FlowsDB.Delete(&flow)
}

func (flow *Flow) Exists() bool {
	return FlowsDB.Where("id = ?", flow.ID).First(&flow).RowsAffected > 0
}

type JSON json.RawMessage

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}
