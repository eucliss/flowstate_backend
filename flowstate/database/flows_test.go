package flowstate

import (
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlowsDatabase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db := FlowsDatabase(logger, "test")
	assert.NotNil(t, db)
}

// TestUserOperations groups all user-related operations
func TestFlowOperations(t *testing.T) {
	dummyFlow := Flow{
		Name:    "Test Flow",
		Owner:   "testuser",
		Content: json.RawMessage([]byte(`{"nodes": [{"id": "1", "type": "queryNode", "position": {"x": 100, "y": 100}, "data": {"label": "Query Node", "content": "Sample query content"}}, {"id": "2", "type": "textNode", "position": {"x": 200, "y": 200}, "data": {"label": "Text Node", "content": "Sample text content"}}]}`)),
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db := FlowsDatabase(logger, "test")

	t.Run("Create Flow", func(t *testing.T) {
		res := db.Create(&dummyFlow)
		assert.Nil(t, res.Error)
		assert.Equal(t, int64(1), res.RowsAffected)
	})

	t.Run("Get Flow", func(t *testing.T) {
		t.Cleanup(func() {
			Drop("flowstate_flows_test")
		})

		db.Create(&dummyFlow)
		flow := dummyFlow.Get()
		logger.Info("Flow", "flow", flow)
		logger.Info("Flow Content", "content", flow.Content)
		assert.Equal(t, dummyFlow.Name, flow.Name)
	})

	t.Run("Check User Exists", func(t *testing.T) {
		assert.True(t, dummyFlow.Exists())
	})

	t.Run("Delete User", func(t *testing.T) {
		dummyFlow.Delete()
		flow := dummyFlow.Get()
		assert.Nil(t, flow)
	})
}
