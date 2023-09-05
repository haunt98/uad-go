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
	//go:embed adl_aosp.json
	adlAOSPBytes []byte

	//go:embed adl_carrier.json
	adlCarrierBytes []byte

	//go:embed adl_google.json
	adlGoogleBytes []byte

	//go:embed adl_misc.json
	adlMiscBytes []byte

	//go:embed adl_oem.json
	adlOEMBytes []byte
)

func main() {
	findStr := os.Getenv("FIND")
	if findStr == "" {
		return
	}
	findStr = strings.ToLower(findStr)

	adlApps := handleADL(findStr)

	apps := make([]UnifiedApp, 0, len(adlApps))
	apps = append(apps, adlApps...)

	for _, a := range apps {
		color.Green("ID: %s", a.ID)
		fmt.Printf("Description: %s\n", a.Description)

		fmt.Println()
	}
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
		// Only get recommend
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
