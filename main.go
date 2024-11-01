package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/bytedance/sonic"
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

	// Include not safe 2 remove
	risk := strings.EqualFold(os.Getenv("RISK"), "true")

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
	for _, findFn := range []func(bool, ...string) []UnifiedApp{
		searchUAD,
	} {
		apps = findFn(risk, values...)
		if len(apps) != 0 {
			break
		}
	}

	for _, a := range apps {
		if a.Safe2Remove {
			color.Green("ID: %s", a.ID)
		} else {
			color.Red("ID: %s", a.ID)
		}
		fmt.Printf("Description: %s\n\n", a.Description)
	}
}

func searchUAD(risk bool, values ...string) []UnifiedApp {
	uadApps := UADApps{}
	if err := sonic.Unmarshal(uadBytes, &uadApps); err != nil {
		log.Fatalf("json: failed to unmarshal: %v", err)
		return nil
	}

	apps := make([]UnifiedApp, 0, len(uadApps))
	for uadAppID, uadApp := range uadApps {
		safe2Remove := uadApp.Removal == UADRemovalRecommended
		if !risk && !safe2Remove {
			continue
		}

		uadAppID = strings.TrimSpace(uadAppID)
		uadApp.Description = strings.TrimSpace(uadApp.Description)

		found := false
		for _, value := range values {
			if strings.Contains(strings.ToLower(uadAppID), value) {
				found = true
				break
			}

			if strings.Contains(strings.ToLower(uadApp.Description), value) {
				found = true
				break
			}

			for _, label := range uadApp.Labels {
				if strings.Contains(strings.ToLower(label), value) {
					found = true
					break
				}
			}

			if found {
				break
			}
		}

		if found {
			apps = append(apps, UnifiedApp{
				ID:          uadAppID,
				Description: uadApp.Description,
				Safe2Remove: safe2Remove,
			})
		}
	}

	// Sort it by safe 2 remove first, then by ID
	sort.Slice(apps, func(i, j int) bool {
		if apps[i].Safe2Remove == apps[j].Safe2Remove {
			return apps[i].ID < apps[j].ID
		}

		return apps[i].Safe2Remove
	})

	return apps
}
