package params

import "encoding/json"

type OrderUid struct {
	OrderUid string `json:"orderUid"`
}
type PublishOrderForDocker struct {
	Text json.RawMessage `json:"text"`
}
