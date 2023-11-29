package util

import (
	"encoding/json"
	"log/slog"
)

func StructToString(s interface{}) string {

	empJSON, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		slog.Error("Could not marshal struct to json string")
		slog.Error(err.Error())
		return ""
	}

	return string(empJSON)

}
