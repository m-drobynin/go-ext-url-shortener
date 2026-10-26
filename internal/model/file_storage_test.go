package model

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileStorage(t *testing.T) {
	fs, err := NewFileStorage(filepath.Join(t.TempDir(), "fs_test.json"))

	assert.NoError(t, err, "error while creating a db")

	key1 := "key_1"
	key2 := "key_2"
	key3 := "key_3"
	value1 := "value_1"
	value2 := "value_2"
	value3 := "value_3"

	fs.TryPut(key1, value1)
	fs.TryPut(key2, value2)
	fs.TryPut(key3, value3)

	err = fs.read()
	assert.NoError(t, err, "should not return an reading error")

	res1, err := fs.Get(key1)
	assert.NoError(t, err, "error while getting value 1")
	assert.Equal(t, res1, &value1)

	res2, err := fs.Get(key2)
	assert.NoError(t, err, "error while getting value 2")
	assert.Equal(t, res2, &value2)

	res3, err := fs.Get(key3)
	assert.NoError(t, err, "error while getting value 3")
	assert.Equal(t, res3, &value3)

	noval, err := fs.Get("key_12345")

	assert.Error(t, err, "should return an error if key does not exist")
	assert.Nil(t, noval, "should return nil value if key does not exist")

	err = fs.TryPut(key1, "smth")
	assert.Error(t, err, "should return an error for putting an existing key")
}
