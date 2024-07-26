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

type GetMessagesStatResponse struct {
	TotalCount       int `json:"total_count"`
	ProcessedCount   int `json:"processed_count"`
	UnProcessedCount int `json:"unprocessed_count"`
}
