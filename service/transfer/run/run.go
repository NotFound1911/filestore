package run

import (
	"github.com/NotFound1911/filestore/config"
	"github.com/NotFound1911/filestore/internal/logger"
	"github.com/NotFound1911/filestore/internal/mq"
	"github.com/NotFound1911/filestore/internal/storage"
	"github.com/NotFound1911/filestore/service/transfer/process"
)

func Run() {
	conf := config.NewConfig("")
	log := logger.New(conf, conf.Service.Transfer.Name)
	msgQueue := mq.New(conf, conf.Service.Transfer.Name, log)
	consumerStorage := storage.New(conf, conf.Service.Transfer.Name, log)
	handler := process.NewHandler(msgQueue, consumerStorage)
	handler.Start()
}
