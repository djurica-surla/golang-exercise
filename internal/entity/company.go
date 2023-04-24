package entity

import (
	"github.com/google/uuid"
)

// Represents a company.
type Company struct {
	ID          uuid.UUID
	Name        string
	Description string
	Employees   int
	Registered  bool
	Type        string
}
