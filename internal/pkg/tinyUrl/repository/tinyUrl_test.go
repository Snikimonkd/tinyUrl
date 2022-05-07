package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCreate_InMemory(t *testing.T) {
	repo := TinyUrlInMemoryRepository{DB: make(map[string]string)}

	err := repo.Create("full", "tiny")
	assert.Equal(t, nil, err)

	res, err := repo.Get("tiny")
	assert.Equal(t, nil, err)

	assert.Equal(t, res, "full")
}

func TestGet_NotFound(t *testing.T) {
	repo := TinyUrlInMemoryRepository{DB: make(map[string]string)}

	res, err := repo.Get("tiny")
	assert.Equal(t, nil, err)

	assert.Equal(t, "", res)
}

func TestCheckIfTinyUrlExist(t *testing.T) {
	repo := TinyUrlInMemoryRepository{DB: make(map[string]string)}

	err := repo.Create("full", "tiny")
	assert.Equal(t, nil, err)

	exist, err := repo.CheckIfTinyUrlExist("tiny")
	assert.Equal(t, nil, err)

	assert.Equal(t, exist, true)
}

func TestCheckIfTinyUrlExist_False(t *testing.T) {
	repo := TinyUrlInMemoryRepository{DB: make(map[string]string)}

	err := repo.Create("full", "tiny")
	assert.Equal(t, nil, err)

	exist, err := repo.CheckIfTinyUrlExist("not tiny")
	assert.Equal(t, nil, err)

	assert.Equal(t, exist, false)
}

func TestCheckIfFullUrlExist(t *testing.T) {
	repo := TinyUrlInMemoryRepository{DB: make(map[string]string)}

	err := repo.Create("full", "tiny")
	assert.Equal(t, nil, err)

	res, err := repo.CheckIfFullUrlExist("full")
	assert.Equal(t, nil, err)

	assert.Equal(t, "tiny", res)
}

func TestCheckIfFullUrlExist_False(t *testing.T) {
	repo := TinyUrlInMemoryRepository{DB: make(map[string]string)}

	err := repo.Create("full", "tiny")
	assert.Equal(t, nil, err)

	res, err := repo.CheckIfFullUrlExist("not full")
	assert.Equal(t, nil, err)

	assert.Equal(t, "", res)
}
