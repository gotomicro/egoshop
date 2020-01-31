package util

import "encoding/json"

func JsonMarshal(d interface{}) string {
	tmp, _ := json.Marshal(d)
	return string(tmp)
}
