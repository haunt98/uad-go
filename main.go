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
	findStr = strings.ToLower(strings.TrimSpace(findStr))

	// For example xiaomi+miui
	// Use + instead of | because | is unix pipe
	values := strings.Split(findStr, "+")
	values = lo.Filter(values, func(v string, i int) bool {
		return strings.TrimSpace(v) != ""
	})

	var apps []UnifiedApp

	// Loop all data sources
	// Stop if found
	// Otherwise keep going
	for _, findFn := range []func(...string) []UnifiedApp{
		searchUAD,
	} {
		apps = findFn(values...)
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

func searchUAD(values ...string) []UnifiedApp {
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
		for _, value := range values {
			if strings.Contains(strings.ToLower(a.ID), value) {
				return UnifiedApp{
					ID:          a.ID,
					Description: description,
				}, true
			}

			if strings.Contains(strings.ToLower(a.Description), value) {
				return UnifiedApp{
					ID:          a.ID,
					Description: description,
				}, true
			}

			for _, label := range a.Labels {
				if strings.Contains(strings.ToLower(label), value) {
					return UnifiedApp{
						ID:          a.ID,
						Description: description,
					}, true
				}
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
