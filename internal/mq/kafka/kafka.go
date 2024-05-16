package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/NotFound1911/filestore/internal/mq/di"
	"github.com/NotFound1911/filestore/pkg/kafka"
	"log"
)

const (
	mqSize int = 1000
)

type Mq struct {
	ctrl     *kafka.Controller
	producer sarama.AsyncProducer
	consumer sarama.ConsumerGroup
	handler  *consumerHandler
}

func (m *Mq) Messages() <-chan *di.Message {
	return m.handler.mc
}

func (m *Mq) SendMessage(message *di.Message) (err error) {
	msgs := m.producer.Input()
	msgs <- m.handler.toProducerMessage(message)
	select {
	case msg := <-m.producer.Successes():
		log.Printf("发送成功:%v\n", string(msg.Value.(sarama.StringEncoder)))
	case pErr := <-m.producer.Errors():
		log.Printf("发送失败:%v,%v\n", pErr.Err, pErr.Msg)
		err = fmt.Errorf("%v:%v", pErr.Err, pErr.Msg)
	}
	return err
}

type consumerHandler struct {
	mc chan *di.Message
}

func (c *consumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
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
func (c *consumerHandler) toMessage(msg *sarama.ConsumerMessage) *di.Message {
	m := &di.Message{
		Topic:   msg.Topic,
		Value:   msg.Value,
		Headers: make([]di.Header, 0, len(msg.Headers)),
	}
	for _, v := range msg.Headers {
		tmp := di.Header{
			Key:   string(v.Key),
			Value: string(v.Value),
		}
		m.Headers = append(m.Headers, tmp)
	}
	return m
}
func (c *consumerHandler) toProducerMessage(msg *di.Message) *sarama.ProducerMessage {
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
func NewMq() di.MessageQueue {
	// todo cfg
	addr := []string{"localhost:9094"}
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	ctrl := kafka.NewController(addr, cfg)
	q := &Mq{
		ctrl: ctrl,
	}
	p, err := sarama.NewAsyncProducer(ctrl.Addr, ctrl.Cfg)
	if err != nil {
		panic(err)
	}
	q.producer = p
	c, err := sarama.NewConsumerGroup(ctrl.Addr, di.TopicName, ctrl.Cfg)
	if err != nil {
		panic(err)
	}
	q.consumer = c
	q.handler = &consumerHandler{
		mc: make(chan *di.Message, mqSize),
	}
	go func() {
		if err := q.consumer.Consume(context.Background(), []string{di.TopicName}, q.handler); err != nil {
			fmt.Println("q.consumer.Consume err:", err)
		}
	}()
	return q
}
