package entities

type CreateMessageRequest struct {
	Content string `json:"content"`
}

type CreateMessageResponse struct {
	MessageId int `json:"message_id"`
}
