package kafka

import (
	"github.com/IBM/sarama"
	"github.com/NotFound1911/filestore/config"
	"time"
)

type Service struct {
	Addr   []string
	Cfg    *sarama.Config
	Enable bool
}

func NewService(conf *config.Configuration) *Service {
	addr := conf.Mq.Addr
	cfg := sarama.NewConfig()
	cfg.Net.DialTimeout = time.Duration(conf.Mq.NetDiaTimeout) * time.Second
	cfg.Net.ReadTimeout = time.Duration(conf.Mq.NetReadTimeout) * time.Second
	cfg.Producer.RequiredAcks = sarama.RequiredAcks(conf.Mq.PReqiredAcks)
	cfg.Producer.Timeout = time.Duration(conf.Mq.PTimeout) * time.Second
	cfg.Producer.Return.Successes = conf.Mq.PReturnSuccess
	cfg.Producer.Return.Errors = conf.Mq.PReturnErr
	cfg.Consumer.MaxWaitTime = time.Duration(conf.Mq.CMaxWaitTime) * time.Millisecond
	cfg.Consumer.Fetch.Min = int32(conf.Mq.CFetchMin)
	cfg.Consumer.Fetch.Default = int32(conf.Mq.CFetchDefault)
	enable := conf.Mq.Enable
	return &Service{
		Addr:   addr,
		Cfg:    cfg,
		Enable: enable,
	}
}
