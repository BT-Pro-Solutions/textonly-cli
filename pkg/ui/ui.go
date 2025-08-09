package ui

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
