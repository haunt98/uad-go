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

//go:embed data/uad_lists.json
var uadBytes []byte

func main() {
	findStr := os.Getenv("FIND")
	if findStr == "" {
		return
	}
	findStr = strings.ToLower(findStr)

	var apps []UnifiedApp

	// Loop all data sources
	// Stop if found
	// Otherwise keep going
	for _, findFn := range []func(string) []UnifiedApp{
		handleUAD,
	} {
		apps = findFn(findStr)
		if len(apps) != 0 {
			break
		}
	}

	for _, a := range apps {
		color.Green("ID: %s", a.ID)
		fmt.Printf("Description: %s\n", a.Description)

		fmt.Println()
	}
}

func handleUAD(findStr string) []UnifiedApp {
	uadApps := UADApps{}
	if err := json.Unmarshal(uadBytes, &uadApps); err != nil {
		slog.Error("json: failed to unmarshal apps", err)
		return nil
	}

	apps := lo.FilterMap(uadApps, func(a UADApp, i int) (UnifiedApp, bool) {
		if a.Removal != UADRemovalRecommended {
			return UnifiedApp{}, false
		}

		a.ID = strings.TrimSpace(a.ID)
		description := strings.TrimSpace(a.Description)

		// Find it
		if strings.Contains(strings.ToLower(a.ID), findStr) {
			return UnifiedApp{
				ID:          a.ID,
				Description: description,
			}, true
		}

		if strings.Contains(strings.ToLower(a.Description), findStr) {
			return UnifiedApp{
				ID:          a.ID,
				Description: description,
			}, true
		}

		for _, label := range a.Labels {
			if strings.Contains(strings.ToLower(label), findStr) {
				return UnifiedApp{
					ID:          a.ID,
					Description: description,
				}, true
			}
		}

		return UnifiedApp{}, false
	})

	// Sort it
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].ID < apps[j].ID
	})

	return apps
}
