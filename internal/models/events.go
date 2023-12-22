package models

type EventType string

const (
	EVENT_TYPE_MESSAGE_CREATED EventType = "message_created"
)

type MessageCreated struct {
	Message *Message
}
