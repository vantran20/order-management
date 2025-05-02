package testutil

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// ToJSONString helper function to convert interface{} to JSON string
func ToJSONString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return `"` + v + `"`
	case gin.H:
		json, _ := json.Marshal(v)
		return string(json)
	default:
		json, _ := json.Marshal(v)
		return string(json)
	}
}
