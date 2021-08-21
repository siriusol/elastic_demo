package main

import "encoding/json"

func GetLogString(object interface{}) string {
	bytes, _ := json.Marshal(object)
	return string(bytes)
}

func GetLogStringPretty(object interface{}) string {
	bytes, _ := json.MarshalIndent(object, "", "\t")
	return string(bytes)
}
