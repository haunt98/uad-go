package main

// https://github.com/0x192/universal-android-debloater/wiki/FAQ

const (
	RemovalRecommended = "Recommended"
	RemovalAdvanced    = "Advanced"
	RemovalExpert      = "Expert"
	RemovalUnsafe      = "Unsafe"
)

// Copy def from uad
type Apps []App

type App struct {
	ID           string   `json:"id,omitempty"`
	List         string   `json:"list,omitempty"`
	Description  string   `json:"description,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	NeededBy     []string `json:"neededBy,omitempty"`
	Labels       []string `json:"labels,omitempty"`
	Removal      string   `json:"removal,omitempty"`
}
