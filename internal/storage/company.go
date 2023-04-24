package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/djurica-surla/golang-exercise/internal/entity"
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
func (store *CompanyStore) GetCompanyByID(ctx context.Context, companyID uuid.UUID) (entity.Company, error) {
	company := entity.Company{}

	err := store.db.QueryRowContext(ctx,
		`SELECT * FROM companies WHERE id = $1`, companyID).
		Scan(&company.ID, &company.Name, &company.Description, &company.Employees, &company.Registered, &company.Type)
	if err != nil {
		return entity.Company{}, fmt.Errorf("error getting company from database %w", err)
	}

	return company, nil
}
