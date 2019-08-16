package models

type ConditionalProperty struct {
	Schema   string `json:"schema" validate:"required"`
	Property string `json:"property" validate:"required"`
	Path     string `json:"path" validate:"required"`
	Required string `json:"required,omitempty"`
}

type ConditionalAPIProperties struct {
	Name       string
	properties []ConditionalProperty
}

type ConditionalPropertyCollection struct {
	APIProperties []ConditionalAPIProperties
}
