package logger

import (
	"github.com/NotFound1911/filestore/config"
	"github.com/NotFound1911/filestore/internal/logger/di"
	z "github.com/NotFound1911/filestore/internal/logger/zap"
	"github.com/NotFound1911/filestore/pkg/zap"
)

func New(conf *config.Configuration, name string) di.Logger {
	s := zap.NewService(conf, name)
	return z.NewLogger(s)
}
