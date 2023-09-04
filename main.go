package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/samber/lo"
)

//go:embed uad_lists.json
var uadBytes []byte

func main() {
	findStr := os.Getenv("FIND")
	if findStr == "" {
		return
	}
	findStr = strings.ToLower(findStr)

	apps := Apps{}
	if err := json.Unmarshal(uadBytes, &apps); err != nil {
		slog.Error("json: failed to unmarshal apps", err)
		return
	}

	apps = lo.Filter(apps, func(app App, i int) bool {
		// Only get recommend
		if app.Removal != RemovalRecommended {
			return false
		}

		// Find it
		if strings.Contains(strings.ToLower(app.ID), findStr) {
			return true
		}

		if strings.Contains(strings.ToLower(app.Description), findStr) {
			return true
		}

		for _, label := range app.Labels {
			if strings.Contains(strings.ToLower(label), findStr) {
				return true
			}
		}

		return false
	})

	// Sort it
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].ID < apps[j].ID
	})

	for _, app := range apps {
		color.Green("ID: %s", app.ID)
		fmt.Printf("Description: %s\n", strings.TrimSpace(app.Description))

		if len(app.Dependencies) > 0 {
			color.Yellow("Dependencies: %s", strings.Join(app.Dependencies, " "))
		}

		if len(app.NeededBy) > 0 {
			color.Magenta("NeededBy: %sv", strings.Join(app.NeededBy, " "))
		}

		fmt.Println()
	}
}
