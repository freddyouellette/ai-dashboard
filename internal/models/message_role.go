package models

type MessageRole string

const (
	MESSAGE_ROLE_USER   MessageRole = "USER"
	MESSAGE_ROLE_BOT    MessageRole = "BOT"
	MESSAGE_ROLE_SYSTEM MessageRole = "SYSTEM"
)
