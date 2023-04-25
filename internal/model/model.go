package model

import (
	"github.com/google/uuid"
)

var companyTypes = map[string]struct{}{
	"Corporation":        {},
	"Cooperative":        {},
	"NonProfit":          {},
	"SoleProprietorship": {},
}

// Represents a company.
type Company struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Employees   int       `json:"employees"`
	Registered  bool      `json:"registered"`
	Type        string    `json:"type"`
}

// CompanyCreate is used for company creation and input validation.
type CompanyCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Employees   int    `json:"employees" validate:"required"`
	Registered  bool   `json:"registered" validate:"required"`
	Type        string `json:"type" validate:"required"`
}

// IsValidType checks if company type is valid and returns all valid types.
func (c CompanyCreate) IsValidType() (bool, []string) {
	var validTypes []string
	_, ok := companyTypes[c.Type]
	for key := range companyTypes {
		validTypes = append(validTypes, key)
	}
	return ok, validTypes
}
