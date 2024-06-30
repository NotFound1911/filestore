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
	var storageHandler di.CustomStorage
	switch conf.Storage.Way {
	case config.LocalStorage:
		return local.NewStorage(conf, local.WithLogger(logger))
	case config.MinioStorage:
		service := m.NewService(conf)
		storageHandler = minio.NewStorage(service, logger)
	default:
		storageHandler = local.NewStorage(conf, local.WithLogger(logger))
	}
	for _, bucket := range []string{"image", "video", "audio", "archive", "unknown", "doc"} {
		if err := storageHandler.MakeBucket(bucket); err != nil {
			panic(err)
		}
	}
	return storageHandler
}
