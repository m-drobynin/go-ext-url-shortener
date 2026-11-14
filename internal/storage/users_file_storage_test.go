package storage

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersFileStorage(t *testing.T) {
	fs, err := NewUsersFileStorage(filepath.Join(t.TempDir(), "fs_test.json"))

	assert.NoError(t, err, "error while creating a db")

	user1, err := fs.Create(t.Context())
	assert.NoError(t, err, "error while creating a user1")

	user2, err := fs.Create(t.Context())
	assert.NoError(t, err, "error while creating a user2")

	user3, err := fs.Create(t.Context())
	assert.NoError(t, err, "error while creating a user3")

	err = fs.read()
	assert.NoError(t, err, "should not return an reading error")

	res1, err := fs.GetByID(t.Context(), user1.UserID)
	assert.NoError(t, err, "error while getting user 1")
	assert.Equal(t, res1.UserID, user1.UserID)

	res2, err := fs.GetByID(t.Context(), user2.UserID)
	assert.NoError(t, err, "error while getting user 2")
	assert.Equal(t, res2.UserID, user2.UserID)

	res3, err := fs.GetByID(t.Context(), user3.UserID)
	assert.NoError(t, err, "error while getting user 3")
	assert.Equal(t, res3.UserID, user3.UserID)
}
