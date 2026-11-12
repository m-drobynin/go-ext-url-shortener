package storage

import (
	"m-drobynin/go-ext-url-shortener/internal/model"
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
	key4 := "key_4"
	key5 := "key_5"

	value1 := "value_1"
	value2 := "value_2"
	value3 := "value_3"
	value4 := "value_4"
	value5 := "value_5"

	fs.Put(t.Context(), key1, value1)
	fs.Put(t.Context(), key2, value2)
	fs.Put(t.Context(), key3, value3)
	fs.PutBatch(t.Context(), []model.ValuesPair{
		{
			Key:   key4,
			Value: value4,
		},
		{
			Key:   key5,
			Value: value5,
		},
	})

	err = fs.read()
	assert.NoError(t, err, "should not return an reading error")

	res1, err := fs.Get(t.Context(), key1)
	assert.NoError(t, err, "error while getting value 1")
	assert.Equal(t, res1, &value1)

	res2, err := fs.Get(t.Context(), key2)
	assert.NoError(t, err, "error while getting value 2")
	assert.Equal(t, res2, &value2)

	res3, err := fs.Get(t.Context(), key3)
	assert.NoError(t, err, "error while getting value 3")
	assert.Equal(t, res3, &value3)

	res4, err := fs.Get(t.Context(), key4)
	assert.NoError(t, err, "error while getting value 4")
	assert.Equal(t, res4, &value4)

	res5, err := fs.Get(t.Context(), key5)
	assert.NoError(t, err, "error while getting value 5")
	assert.Equal(t, res5, &value5)

	noval, err := fs.Get(t.Context(), "key_12345")

	assert.Error(t, err, "should return an error if key does not exist")
	assert.Nil(t, noval, "should return nil value if key does not exist")

	err = fs.Put(t.Context(), key1, "smth")
	assert.Error(t, err, "should return an error for putting an existing key")
}
