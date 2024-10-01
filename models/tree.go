package models

import "github.com/google/uuid"

// Tree represents a tree in a plantation estate.
type Tree struct {
	ID       uuid.UUID `json:"id"`       // Unique identifier for the tree
	EstateID uuid.UUID `json:"estate_id"` // ID of the estate this tree belongs to
	X        int       `json:"x"`        // X coordinate of the tree in its plot
	Y        int       `json:"y"`        // Y coordinate of the tree in its plot
	Height   int       `json:"height"`   // Height of the tree in meters (1 to 30)
}
