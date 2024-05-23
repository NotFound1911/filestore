package storage

import (
	"github.com/NotFound1911/filestore/internal/storage/di"
	"github.com/NotFound1911/filestore/internal/storage/local"
)

func New() di.CustomStorage {
	// todo config
	return local.NewStorage()
}
