package flowstate

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var dummyUser = User{Username: "test", Password: "test"}

// TestUserOperations groups all user-related operations
func TestUserOperations(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db := UsersDatabase(logger, "test")
	t.Run("Create User", func(t *testing.T) {
		res := db.Create(&dummyUser)
		assert.Nil(t, res.Error)
		assert.Equal(t, int64(1), res.RowsAffected)
	})

	t.Run("Get User", func(t *testing.T) {
		t.Cleanup(func() {
			Drop("flowstate_users_test")
		})

		db.Create(&dummyUser)
		user := dummyUser.Get()
		assert.Equal(t, dummyUser.Username, user.Username)
	})

	t.Run("Check User Exists", func(t *testing.T) {
		assert.True(t, dummyUser.Exists())
	})

	t.Run("Login Success", func(t *testing.T) {
		assert.True(t, dummyUser.LoginSuccess())
	})

	t.Run("Delete User", func(t *testing.T) {
		dummyUser.Delete()
		user := dummyUser.Get()
		assert.Nil(t, user)
	})
}
