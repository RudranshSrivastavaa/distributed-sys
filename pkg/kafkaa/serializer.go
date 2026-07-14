package kafkaa

import "encoding/json"

func Serialize[T any](event Event) ([]byte, error) {

	return json.Marshal(event)

}