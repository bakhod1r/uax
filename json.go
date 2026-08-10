package uax

import "encoding/json"

// ToJSON returns the JSON representation of the Client.
func (c *Client) ToJSON() []byte {
	b, _ := json.Marshal(c)
	return b
}
