package run

import (
	"github.com/NotFound1911/filestore/internal/mq"
	"github.com/NotFound1911/filestore/internal/storage"
	"github.com/NotFound1911/filestore/service/transfer/process"
)

func Run() {
	// todo config
	msgQueue := mq.New()
	consumerStorage := storage.New()
	handler := process.NewHandler(msgQueue, consumerStorage)
	handler.Start()
}
