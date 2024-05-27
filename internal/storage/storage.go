package storage

import (
	"github.com/NotFound1911/filestore/config"
	ldi "github.com/NotFound1911/filestore/internal/logger/di"
	"github.com/NotFound1911/filestore/internal/storage/di"
	"github.com/NotFound1911/filestore/internal/storage/local"
)

func New(conf *config.Configuration, name string, logger ldi.Logger) di.CustomStorage {
	// todo config
	switch conf.Storage.Way {
	case config.LocalStorage:
		return local.NewStorage(conf, local.WithLogger(logger))
	default:
		return local.NewStorage(conf, local.WithLogger(logger))
	}

}
