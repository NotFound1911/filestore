package storage

import (
	"github.com/NotFound1911/filestore/config"
	ldi "github.com/NotFound1911/filestore/internal/logger/di"
	"github.com/NotFound1911/filestore/internal/storage/di"
	"github.com/NotFound1911/filestore/internal/storage/local"
	"github.com/NotFound1911/filestore/internal/storage/minio"
	m "github.com/NotFound1911/filestore/pkg/minio"
)

func New(conf *config.Configuration, logger ldi.Logger) di.CustomStorage {
	// todo config
	switch conf.Storage.Way {
	case config.LocalStorage:
		return local.NewStorage(conf, local.WithLogger(logger))
	case config.MinioStorage:
		service := m.NewService(conf)
		return minio.NewStorage(service, logger)
	default:
		return local.NewStorage(conf, local.WithLogger(logger))
	}

}
