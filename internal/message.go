package internal

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateMessageID() string {
	// Generate a UUID (Universally Unique Identifier)
	uuid := uuid.New().String()

	// Get the current Unix timestamp in milliseconds
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Concatenate the timestamp and UUID to create a unique messageId
	messageID := fmt.Sprintf("%d-%s", timestamp, uuid)

	return messageID
}
