package flowstate

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessDatabase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db := AccessDatabase(logger, "test")
	assert.NotNil(t, db)
}

// TestUserOperations groups all user-related operations
func TestAccessOperations(t *testing.T) {
	dummyAccess := Access{
		UserID: 1,
		FlowID: 1,
		Access: "read",
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db := AccessDatabase(logger, "test")

	t.Run("Create Access", func(t *testing.T) {
		res := db.Create(&dummyAccess)
		assert.Nil(t, res.Error)
		assert.Equal(t, int64(1), res.RowsAffected)
	})

	t.Run("Get Flow", func(t *testing.T) {
		t.Cleanup(func() {
			Drop("flowstate_access_test")
		})

		db.Create(&dummyAccess)
		access := dummyAccess.Get()
		logger.Info("Access", "access", access)
		assert.Equal(t, dummyAccess.UserID, access.UserID)
		assert.Equal(t, dummyAccess.FlowID, access.FlowID)
		assert.Equal(t, dummyAccess.Access, access.Access)
	})

	t.Run("Check User Exists", func(t *testing.T) {
		assert.True(t, dummyAccess.Exists())
	})

	t.Run("Delete User", func(t *testing.T) {
		dummyAccess.Delete()
		access := dummyAccess.Get()
		assert.Nil(t, access)
	})
}
