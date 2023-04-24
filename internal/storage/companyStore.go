package storage

import (
	"database/sql"
)

// Represents implementation of company storage.
type CompanyStore struct {
	db *sql.DB
}

// NewCompanyStore creates a new instance of the CompanyStore.
func NewCompanyStore(connection *sql.DB) *CompanyStore {
	return &CompanyStore{db: connection}
}
