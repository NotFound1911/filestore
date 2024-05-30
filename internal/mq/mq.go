package mq

import (
	"github.com/NotFound1911/filestore/config"
	ldi "github.com/NotFound1911/filestore/internal/logger/di"
	mdi "github.com/NotFound1911/filestore/internal/mq/di"
	"github.com/NotFound1911/filestore/internal/mq/kafka"
	k "github.com/NotFound1911/filestore/pkg/kafka"
)

func New(conf *config.Configuration, logger ldi.Logger) mdi.MessageQueue {
	s := k.NewService(conf)
	return kafka.NewMq(s, logger)
}
