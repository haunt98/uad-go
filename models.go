package main

// https://github.com/0x192/universal-android-debloater/wiki/FAQ
// https://github.com/Universal-Debloater-Alliance/universal-android-debloater-next-generation/wiki/FAQ
const UADRemovalRecommended = "Recommended"

type UADApps map[string]UADApp

type UADApp struct {
	List         string   `json:"list,omitempty"`
	Description  string   `json:"description,omitempty"`
	Removal      string   `json:"removal,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	NeededBy     []string `json:"neededBy,omitempty"`
	Labels       []string `json:"labels,omitempty"`
}

type UnifiedApp struct {
	ID          string
	Description string
	Safe2Remove bool
}
