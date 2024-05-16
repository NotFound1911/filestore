package process

import (
	mdi "github.com/NotFound1911/filestore/internal/mq/di"
	sdi "github.com/NotFound1911/filestore/internal/storage/di"
)

type Handler struct {
	msgQueue        mdi.MessageQueue
	customerStorage sdi.CustomStorage
}

func NewHandler(msgQueue mdi.MessageQueue, customerStorage sdi.CustomStorage) *Handler {
	return &Handler{
		msgQueue:        msgQueue,
		customerStorage: customerStorage,
	}
}

// transfer 转移服务，将文件转移到对应的存储
func (h *Handler) transfer(msg *mdi.Message) error {
	var bucket string
	var storageName string
	var location string
	for _, v := range msg.Headers {
		switch v.Key {
		case mdi.HeaderLocation:
			location = v.Value
		case mdi.HeaderBucket:
			bucket = v.Value
		case mdi.HeaderStorageName:
			storageName = v.Value
		}
	}
	return h.customerStorage.PutObject(bucket, storageName, location, "")
}
func (h *Handler) Start() {
	msgs := h.msgQueue.Messages()
	for msg := range msgs {
		if err := h.transfer(msg); err != nil {
			continue
		}
	}
}
