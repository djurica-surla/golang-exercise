package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/djurica-surla/golang-exercise/internal/model"
)

// CompanyStorer represents necessary company storage implementation for Company service.
type CompanyStorer interface {
	GetCompanyByID(ctx context.Context, companyID uuid.UUID) (model.Company, error)
	CreateCompany(ctx context.Context, company model.CompanyCreate) error
	UpdateCompany(ctx context.Context, company model.CompanyCreate, companyID uuid.UUID) error
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

// GetCompanyByID handles the logic for getting the company by the id.
func (s *CompanyService) GetCompanyByID(ctx context.Context, companyID uuid.UUID) (model.Company, error) {
	company, err := s.companyStore.GetCompanyByID(ctx, companyID)
	if err != nil {
		return company, err
	}

	return company, err
}

// CreateCompany handles the logic for creating a company.
func (s *CompanyService) CreateCompany(ctx context.Context, company model.CompanyCreate) error {
	err := s.companyStore.CreateCompany(ctx, company)
	if err != nil {
		return err
	}

	return nil
}

// UpdateCompany handles the logic for updating a company.
func (s *CompanyService) UpdateCompany(ctx context.Context, company model.CompanyCreate, companyID uuid.UUID) error {
	err := s.companyStore.UpdateCompany(ctx, company, companyID)
	if err != nil {
		return err
	}

	return nil
}
