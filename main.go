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

var (
	//go:embed data/uad_lists.json
	uadBytes []byte

	//go:embed data/adl_aosp.json
	adlAOSPBytes []byte

	//go:embed data/adl_carrier.json
	adlCarrierBytes []byte

	//go:embed data/adl_google.json
	adlGoogleBytes []byte

	//go:embed data/adl_misc.json
	adlMiscBytes []byte

	//go:embed data/adl_oem.json
	adlOEMBytes []byte
)

func main() {
	findStr := os.Getenv("FIND")
	if findStr == "" {
		return
	}
	findStr = strings.ToLower(findStr)

	var apps []UnifiedApp

	if strings.EqualFold(strings.TrimSpace(os.Getenv("MODE")), "uad") {
		apps = handleUAD(findStr)
	} else {
		// Default mode is adl
		apps = handleADL(findStr)
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

func handleADL(findStr string) []UnifiedApp {
	adlApps := ADLApps{}
	for _, bytes := range [][]byte{
		adlAOSPBytes,
		adlCarrierBytes,
		adlGoogleBytes,
		adlMiscBytes,
		adlOEMBytes,
	} {
		adlSubApps := ADLApps{}
		if err := json.Unmarshal(bytes, &adlSubApps); err != nil {
			slog.Error("json: failed to unmarshal apps", err)
			return nil
		}

		adlApps = append(adlApps, adlSubApps...)
	}

	apps := lo.FilterMap(adlApps, func(a ADLApp, i int) (UnifiedApp, bool) {
		if a.Removal != ADLRemovalDelete {
			return UnifiedApp{}, false
		}

		a.ID = strings.TrimSpace(a.ID)
		description := strings.TrimSpace(a.Description) + " " + strings.TrimSpace(a.Warning)

		// Find it
		if strings.Contains(strings.ToLower(a.ID), findStr) {
			return UnifiedApp{
				ID:          a.ID,
				Description: description,
			}, true
		}

		if strings.Contains(strings.ToLower(a.Label), findStr) {
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

		if strings.Contains(strings.ToLower(a.Warning), findStr) {
			return UnifiedApp{
				ID:          a.ID,
				Description: description,
			}, true
		}

		return UnifiedApp{}, false
	})

	// Sort it
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].ID < apps[j].ID
	})

	return apps
}
