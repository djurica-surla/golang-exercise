package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/djurica-surla/golang-exercise/internal/entity"
)

// CompanyStorer represents necessary company storage implementation for Company service.
type CompanyStorer interface {
	GetCompanyByID(ctx context.Context, companyID uuid.UUID) (entity.Company, error)
}

// CompanyService contains business logic for working with company object.
type CompanyService struct {
	companyStore CompanyStorer
}

// Instantiates a new company service struct with a company repo.
func NewCompanyService(companyStore CompanyStorer) *CompanyService {
	return &CompanyService{
		companyStore: companyStore,
	}
}

// GetCompanyByID handles the logic for getting a company by the id.
func (s *CompanyService) GetCompanyByID(ctx context.Context, CompanyID uuid.UUID) (entity.Company, error) {
	company, err := s.companyStore.GetCompanyByID(ctx, CompanyID)
	if err != nil {
		return company, err
	}

	return company, err
}
