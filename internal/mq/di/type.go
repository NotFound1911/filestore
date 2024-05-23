package di

const (
	TopicName         = "filestore"
	HeaderBucket      = "bucket"
	HeaderStorageName = "storage_name"
	HeaderLocation    = "location"
)

type MessageQueue interface {
	Messages() <-chan *Message
	SendMessage(message *Message) error
	Enable() bool
}

type Message struct {
	Topic   string
	Value   []byte
	Headers []Header
}
type Header struct {
	Key   string
	Value string
}
