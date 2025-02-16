package flowstate

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyStruct struct {
	ID   uint
	Name string
}

// TestDatabaseOperations groups database creation/deletion tests
func TestDatabaseOps(t *testing.T) {
	databaseLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	t.Run("Create Database", func(t *testing.T) {
		res := Create("dummy_db")
		assert.True(t, res)
	})

	t.Run("Drop Database", func(t *testing.T) {
		res := Drop("dummy_db")
		assert.True(t, res)
	})

	t.Run("Load Database", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		db := Database(logger, "test", "flowstate_users", dummyStruct{})
		assert.NotNil(t, db)
	})
}
