package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
)

//go:embed uad_lists.json
var uadBytes []byte

func main() {
	apps := Apps{}
	if err := json.Unmarshal(uadBytes, &apps); err != nil {
		slog.Error("json: failed to unmarshal apps", err)
		return
	}

	fmt.Println(len(apps))
}
