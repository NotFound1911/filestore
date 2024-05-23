package mq

import (
	"github.com/NotFound1911/filestore/config"
	"github.com/NotFound1911/filestore/internal/mq/di"
	"github.com/NotFound1911/filestore/internal/mq/kafka"
	k "github.com/NotFound1911/filestore/pkg/kafka"
)

func New(conf *config.Configuration, name string) di.MessageQueue {
	s := k.NewService(conf, name)
	return kafka.NewMq(s)
}
