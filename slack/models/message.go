package models

type Message struct {
	Sender    string
	Receiver  string
	Content   string
	Timestamp int64
	Type      MessageType `json:"type"`
}

type MessageType string

const (
	BroadCast     = "broadCast"
	DirectMessage = "directMessage"
)

type WebsocketRequest struct {
	Client  string
	Type    MessageType `json:"type"`
	Message string
}
