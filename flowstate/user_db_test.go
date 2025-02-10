package flowstate

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var dummyUser = User{Username: "test", Password: "test"}

func loadTestDB() *gorm.DB {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db := LoadUserDB(logger, "flowstate_users_test")
	if db == nil {
		panic("Failed to load test database")
	}
	return db
}

// TestDatabaseOperations groups database creation/deletion tests
func TestDatabaseOperations(t *testing.T) {
	t.Run("Create Database", func(t *testing.T) {
		res := CreateDB("dummy_db")
		assert.True(t, res)
	})

	t.Run("Drop Database", func(t *testing.T) {
		res := DropDB("dummy_db")
		assert.True(t, res)
	})

	t.Run("Load Database", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		db := LoadUserDB(logger, "flowstate_users_test")
		assert.NotNil(t, db)
	})
}

// TestUserOperations groups all user-related operations
func TestUserOperations(t *testing.T) {
	t.Run("Create User", func(t *testing.T) {
		db := loadTestDB()
		res := db.Create(&dummyUser)
		assert.Nil(t, res.Error)
		assert.Equal(t, int64(1), res.RowsAffected)
	})

	t.Run("Get User", func(t *testing.T) {
		db := loadTestDB()
		t.Cleanup(func() {
			DropDB("flowstate_users_test")
		})

		db.Create(&dummyUser)
		user := dummyUser.Get()
		assert.Equal(t, dummyUser.Username, user.Username)
	})

	t.Run("Check User Exists", func(t *testing.T) {
		loadTestDB()
		assert.True(t, dummyUser.Exists())
	})

	t.Run("Delete User", func(t *testing.T) {
		loadTestDB()
		dummyUser.Delete()
		user := dummyUser.Get()
		assert.Nil(t, user)
	})
}
