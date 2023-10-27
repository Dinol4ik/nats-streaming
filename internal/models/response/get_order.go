package response

import "encoding/json"

type GetOrder struct {
	OrderUid  string          `json:"order_uid"`
	OrderInfo json.RawMessage `json:"order_info"`
}
