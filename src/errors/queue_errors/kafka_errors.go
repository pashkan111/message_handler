package queue_errors

import "fmt"

type PublishError struct {
	Message string
}

func (e PublishError) Error() string {
	return fmt.Sprintf(
		"Error publishing message to kafka. Message: %s", e.Message,
	)
}
