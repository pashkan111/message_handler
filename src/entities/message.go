package entities

type MessageFromRequest struct {
	Content string
}

type MessageToProduce struct {
	Content   string
	MessageId int
}

type CreateMessageResponse struct {
	MessageId int `json:"message_id"`
}
