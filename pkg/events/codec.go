package events

import "encoding/json"

func Marshal(event Event) ([]byte, error) {

	return json.Marshal(event)

}