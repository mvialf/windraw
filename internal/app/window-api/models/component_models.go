package models

type AuxProfile struct {
	ID    string  `json:"id"`
	Angle float64 `json:"angle,omitempty"`
	Notes string  `json:"notes,omitempty"`
}

type Module struct {
	ID         string      `json:"id"`
	Elements   []Element
	AuxProfile *AuxProfile `json:"aux_profile,omitempty"`
}

type Component struct {
	ID      string   `json:"id"`
	Modules []Module `json:"modules"`
}