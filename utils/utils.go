// Package utils provides basic operation and const for common use
package utils

import (
	"encoding/json"
	"fmt"
)

var (
	Version   = "unknown version"
	BuildTime = "unknown date"
)

func ShowJsonRsp(rsp []byte) error {
	var data map[string]interface{}

	err := json.Unmarshal(rsp, &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err
	}

	for key, value := range data {
		fmt.Printf("AI response| %s: %v\n", key, value)
	}
	return nil
}
