package util

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintStruct(s interface{}) {

	empJSON, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s\n", string(empJSON))

}
