package utils

import "encoding/json"

func IsValidJSON(content []byte) bool {
	var message json.RawMessage
	err := json.Unmarshal(content, &message)
	return err == nil
}
