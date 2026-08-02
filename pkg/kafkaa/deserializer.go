package kafkaa

import "encoding/json"

func Deserialize(data []byte,event *Event) error {

	return json.Unmarshal(data,event)
}