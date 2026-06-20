package ws

import "encoding/json"

const MessageTypeSnapshot = "snapshot"

// Envelope is the wire format for WebSocket messages.
type Envelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// MarshalEnvelope serializes a typed message payload.
func MarshalEnvelope(messageType string, payload any) ([]byte, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return json.Marshal(Envelope{
		Type:    messageType,
		Payload: raw,
	})
}
