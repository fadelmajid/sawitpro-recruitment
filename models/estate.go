package models

import "github.com/google/uuid"

// Estate represents a plantation estate with dimensions.
type Estate struct {
	ID     uuid.UUID `json:"id"`     // Unique identifier for the estate
	Width  int       `json:"width"`  // Width of the estate in 10m plots
	Length int       `json:"length"` // Length of the estate in 10m plots
}
