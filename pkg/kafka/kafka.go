package kafka

import "github.com/IBM/sarama"

type Controller struct {
	Addr []string
	Cfg  *sarama.Config
}

func NewController(addr []string, cfg *sarama.Config) *Controller {
	return &Controller{
		Addr: addr,
		Cfg:  cfg,
	}
}
