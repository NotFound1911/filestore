package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	ldi "github.com/NotFound1911/filestore/internal/logger/di"
	mdi "github.com/NotFound1911/filestore/internal/mq/di"
	"github.com/NotFound1911/filestore/pkg/kafka"
)

const (
	mqSize int = 1000
)

var _ mdi.MessageQueue = &Mq{}

type Mq struct {
	ctrl     *kafka.Service
	producer sarama.AsyncProducer
	consumer sarama.ConsumerGroup
	handler  *consumerHandler
	logger   ldi.Logger
}

func (m *Mq) Enable() bool {
	return m.ctrl.Enable
}

func (m *Mq) Messages() <-chan *mdi.Message {
	return m.handler.mc
}

func (m *Mq) SendMessage(message *mdi.Message) (err error) {
	msgs := m.producer.Input()
	msgs <- m.handler.toProducerMessage(message)
	select {
	case msg := <-m.producer.Successes():
		m.logger.Info(fmt.Sprintf("发送成功:%v\n", string(msg.Value.(sarama.StringEncoder))))
	case pErr := <-m.producer.Errors():
		m.logger.Error(fmt.Sprintf("发送失败:%v,%v\n", pErr.Err, pErr.Msg))
		err = fmt.Errorf("%v:%v", pErr.Err, pErr.Msg)
	}
	return err
}

type consumerHandler struct {
	mc chan *mdi.Message
}

func (c *consumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		c.mc <- c.toMessage(msg)
		session.MarkMessage(msg, "")
	}
	return nil
}
func (c *consumerHandler) toMessage(msg *sarama.ConsumerMessage) *mdi.Message {
	m := &mdi.Message{
		Topic:   msg.Topic,
		Value:   msg.Value,
		Headers: make([]mdi.Header, 0, len(msg.Headers)),
	}
	for _, v := range msg.Headers {
		tmp := mdi.Header{
			Key:   string(v.Key),
			Value: string(v.Value),
		}
		m.Headers = append(m.Headers, tmp)
	}
	return m
}
func (c *consumerHandler) toProducerMessage(msg *mdi.Message) *sarama.ProducerMessage {
	m := &sarama.ProducerMessage{
		Topic:   msg.Topic,
		Value:   sarama.StringEncoder(msg.Value),
		Headers: make([]sarama.RecordHeader, 0, len(msg.Headers)),
	}
	for _, v := range msg.Headers {
		tmp := sarama.RecordHeader{
			Key:   []byte(v.Key),
			Value: []byte(v.Value),
		}
		m.Headers = append(m.Headers, tmp)
	}
	return m
}
func NewMq(ctrl *kafka.Service, logger ldi.Logger) mdi.MessageQueue {
	q := &Mq{
		ctrl:   ctrl,
		logger: logger,
	}
	p, err := sarama.NewAsyncProducer(ctrl.Addr, ctrl.Cfg)
	if err != nil {
		panic(err)
	}
	q.producer = p
	c, err := sarama.NewConsumerGroup(ctrl.Addr, mdi.TopicName, ctrl.Cfg)
	if err != nil {
		panic(err)
	}
	q.consumer = c
	q.handler = &consumerHandler{
		mc: make(chan *mdi.Message, mqSize),
	}
	go func() {
		if err := q.consumer.Consume(context.Background(), []string{mdi.TopicName}, q.handler); err != nil {
			q.logger.Error(fmt.Sprintf("q.consumer.Consume err:%v", err))
		}
	}()
	return q
}
