package main

import "encoding/json"

func GetLogString(object interface{}) string {
	bytes, _ := json.Marshal(object)
	return string(bytes)
}
