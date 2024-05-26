package storage

import (
	"github.com/NotFound1911/filestore/config"
	"github.com/NotFound1911/filestore/internal/storage/di"
	"github.com/NotFound1911/filestore/internal/storage/local"
)

func New(conf *config.Configuration, name string) di.CustomStorage {
	// todo config
	switch conf.Storage.Way {
	case config.LocalStorage:
		return local.NewStorage(conf)
	default:
		return local.NewStorage(conf)
	}

}
