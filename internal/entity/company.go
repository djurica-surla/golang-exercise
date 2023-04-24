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
	Type        CompanyType
}

type CompanyType int64

// Enum for different company types.
const (
	Corporation CompanyType = iota
	NonProfit
	Cooperative
	SoleProprietorship
)

// func (c CompanyType) fromString(s string) (CompanyType, error) {
// 	switch s {
// 	case "Corporation":
// 		return Corporation, nil
// 	case "NonProfit":
// 		return NonProfit, nil
// 	case "Cooperative":
// 		return Cooperative, nil
// 	case "SoleProprietorship":
// 		return SoleProprietorship, nil
// 	default:
// 		return 0, fmt.Errorf("uknown type")
// 	}
// }
