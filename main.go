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

//go:embed adl_aosp.json
var adlAOSP []byte

//go:embed adl_carrier.json
var adlCarrier []byte

//go:embed adl_google.json
var adlGoogle []byte

//go:embed adl_misc.json
var adlMisc []byte

//go:embed adl_oem.json
var adlOEM []byte

func main() {
	findStr := os.Getenv("FIND")
	if findStr == "" {
		return
	}
	findStr = strings.ToLower(findStr)

	apps := handleUAD(findStr)

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
		// Only get recommend
		if a.Removal != UADRemovalRecommended {
			return UnifiedApp{}, false
		}

		a.ID = strings.TrimSpace(a.ID)
		a.Description = strings.TrimSpace(a.Description)

		// Find it
		if strings.Contains(strings.ToLower(a.ID), findStr) {
			return UnifiedApp{
				ID:          a.ID,
				Description: a.Description,
			}, true
		}

		if strings.Contains(strings.ToLower(a.Description), findStr) {
			return UnifiedApp{
				ID:          a.ID,
				Description: a.Description,
			}, true
		}

		for _, label := range a.Labels {
			if strings.Contains(strings.ToLower(label), findStr) {
				return UnifiedApp{
					ID:          a.ID,
					Description: a.Description,
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
