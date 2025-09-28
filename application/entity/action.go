package entity

import "encoding/json"

type ActionData struct {
	Method  string            `json:"method" binding:"required"`
	URL     string            `json:"url" binding:"required"`
	Headers map[string]string `json:"headers,omitempty"`
	Payload json.RawMessage   `json:"payload,omitempty"`
}
