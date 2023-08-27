package main

// Copy def from uad
type Apps []App

type App struct {
	ID           string   `json:"id:omitempty"`
	List         string   `json:"list:omitempty"`
	Description  string   `json:"description:omitempty"`
	Dependencies []string `json:"dependencies:omitempty"`
	NeededBy     []string `json:"neededBy:omitempty"`
	Labels       []string `json:"labels:omitempty"`
	Removal      string   `json:"removal:omitempty"`
}
