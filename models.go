package main

// https://github.com/0x192/universal-android-debloater/wiki/FAQ
const UADRemovalRecommended = "Recommended"

type UADApps []UADApp

type UADApp struct {
	ID           string   `json:"id,omitempty"`
	List         string   `json:"list,omitempty"`
	Description  string   `json:"description,omitempty"`
	Removal      string   `json:"removal,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	NeededBy     []string `json:"neededBy,omitempty"`
	Labels       []string `json:"labels,omitempty"`
}

// https://github.com/MuntashirAkon/android-debloat-list
const ADLRemovalDelete = "delete"

type ADLApps []ADLApp

type ADLApp struct {
	ID           string   `json:"id,omitempty"`
	Label        string   `json:"label,omitempty"`
	Description  string   `json:"description,omitempty"`
	Removal      string   `json:"removal,omitempty"`
	Warning      string   `json:"warning,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	RequiredBy   []string `json:"required_by,omitempty"`
}

type UnifiedApp struct {
	ID          string
	Description string
}
