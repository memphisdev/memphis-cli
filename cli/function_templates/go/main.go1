package main

import (
	"encoding/json"
	"github.com/memphisdev/memphis-functions.go/memphis"
)

// https://github.com/memphisdev/memphis.go#creating-a-memphis-function
func EventHandler(message []byte, headers map[string]string, inputs map[string]string) ([]byte, map[string]string,  error){
	// Here is a short example of converting the message payload to bytes and back

	var event map[string]interface{}
	json.Unmarshal(message, &event)
	event[inputs["field_to_ingest"]] = "Hello from Memphis!"
	
	// Return the payload back as []bytes
	eventBytes, _ := json.Marshal(event)
	return eventBytes, headers, nil
}


func main() {	
	memphis.CreateFunction(EventHandler)
}