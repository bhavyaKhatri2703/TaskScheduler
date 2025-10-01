package entity

type ActionData struct {
	Method  string            `json:"method" binding:"required"`
	URL     string            `json:"url" binding:"required"`
	Headers map[string]string `json:"headers,omitempty"`
	Payload interface{}       `json:"payload,omitempty"`
}
