package utils

import "encoding/json"

func IsValidJSON(content string) bool {
	var message json.RawMessage
	err := json.Unmarshal([]byte(content), &message)
	return err == nil
}
