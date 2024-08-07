package models

type Message struct {
	BaseEntity
	UserScopedEntity
	ChatID     uint        `json:"chat_id"`
	Text       string      `json:"text"`
	Correction string      `json:"correction"`
	Role       MessageRole `json:"role"`
	Active     bool        `json:"active"`
	BreakAfter bool        `json:"break_after"`
}

func (m *Message) SetUserId(userId uint) {
	m.UserId = userId
}

type GetMessagesOptions struct {
	ChatID  uint `json:"chat_id"`
	Page    int  `json:"page"`
	PerPage int  `json:"per_page"`
}

type MessagesDTO struct {
	Messages []*Message `json:"messages"`
	Page     int        `json:"page"`
	PerPage  int        `json:"per_page"`
	Total    int        `json:"total"`
}
