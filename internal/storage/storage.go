package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/djurica-surla/golang-exercise/internal/model"
)

// Represents implementation of company storage.
type CompanyStore struct {
	db *sql.DB
}

// NewCompanyStore creates a new instance of the CompanyStore.
func NewCompanyStore(connection *sql.DB) *CompanyStore {
	return &CompanyStore{db: connection}
}

// Retrieves a company from the database by the id.
func (store *CompanyStore) GetCompany(ctx context.Context, companyID uuid.UUID) (model.Company, error) {
	company := model.Company{}

	err := store.db.QueryRowContext(ctx,
		`SELECT * FROM companies WHERE id = $1`, companyID).
		Scan(&company.ID, &company.Name, &company.Description, &company.Employees, &company.Registered, &company.Type)
	if err != nil {
		return model.Company{}, fmt.Errorf("error getting company from database %w", err)
	}

	return company, nil
}

// Creates a new company in the database and returns the id of the created company.
func (store *CompanyStore) CreateCompany(ctx context.Context, company model.CompanyCreate) (uuid.UUID, error) {
	var companyID uuid.UUID

	err := store.db.QueryRowContext(ctx,
		`INSERT INTO companies (name, description, employees, registered, type)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		company.Name, company.Description, company.Employees, company.Registered, company.Type).Scan(&companyID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return uuid.Nil, fmt.Errorf("error creating company in the database - company with that name already exists")
		}
		return uuid.Nil, fmt.Errorf("error creating company in the database %w", err)
	}

	return companyID, nil
}

// Update company updates the company in the database.
func (store *CompanyStore) UpdateCompany(ctx context.Context, company model.CompanyCreate, companyID uuid.UUID) error {
	res, err := store.db.ExecContext(ctx,
		`UPDATE companies SET 
		name = $1,
		description = $2,
		employees = $3,
		registered = $4,
		type = $5
		WHERE id = $6`, company.Name, company.Description, company.Employees, company.Registered, company.Type, companyID)
	n, _ := res.RowsAffected()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return fmt.Errorf("error creating company in the database - company with that name already exists")
		}
		return fmt.Errorf("error updating company in the database %w", err)
	}
	if n == 0 {
		return fmt.Errorf("error updating company in the database - record doesnt exit found")
	}

	return nil
}

// Delete company deletes the company from the database.
func (store *CompanyStore) DeleteCompany(ctx context.Context, companyID uuid.UUID) error {
	err := store.db.QueryRowContext(ctx,
		`DELETE FROM companies
		WHERE id = $1`, companyID).Err()
	if err != nil {
		return fmt.Errorf("error deleting company from the database %w", err)
	}

	return nil
}
