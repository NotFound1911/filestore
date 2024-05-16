package mq

import (
	"github.com/NotFound1911/filestore/internal/mq/di"
	"github.com/NotFound1911/filestore/internal/mq/kafka"
)

func New() di.MessageQueue {
	// todo config
	return kafka.NewMq()
}
