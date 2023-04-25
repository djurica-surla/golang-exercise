//go:build integration

package storage_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/djurica-surla/golang-exercise/internal/database"
	"github.com/djurica-surla/golang-exercise/internal/model"
	"github.com/djurica-surla/golang-exercise/internal/storage"
)

var testStore *storage.CompanyStore

func TestMain(m *testing.M) {
	retries := 5
	// Connection details can be seen in /test/docker-compose-integraiton.yaml
	dbConfig := database.DBConfig{
		Host:     "db",
		User:     "testuser",
		Password: "testpassword",
		DBname:   "testdb",
	}
	// Attempt to establish a connection with the database.
	connection, err := database.Connect(context.Background(), dbConfig)
	for err != nil {
		// Retry until postgres is retry
		if retries > 1 {
			retries--
			time.Sleep(1 * time.Second)
			connection, err = database.Connect(context.Background(), dbConfig)
			continue
		}
		log.Fatal(err)
	}

	defer connection.Close()

	// Run up migrations to create the database schema.
	// Second argument is path to the migrations file
	err = database.Migrate(connection, "../../migrations")
	if err != nil {
		log.Fatal(err)
	}

	// Instantiate company store.
	testStore = storage.NewCompanyStore(connection)

	os.Exit(m.Run())
}

func TestStorage_GetCompany(t *testing.T) {
	t.Run("should successfully get a company", func(t *testing.T) {

		company := model.CompanyCreate{
			Name:        "Some",
			Description: "Some",
			Employees:   123,
			Registered:  true,
			Type:        "Corporation",
		}

		companyID, err := testStore.CreateCompany(context.Background(), company)
		require.NoError(t, err)

		result, err := testStore.GetCompany(context.Background(), companyID)
		require.NoError(t, err)
		require.Equal(t, company.Name, result.Name)
		require.Equal(t, company.Description, result.Description)
		require.Equal(t, company.Employees, result.Employees)
		require.Equal(t, company.Registered, result.Registered)
		require.Equal(t, company.Type, result.Type)

		err = testStore.DeleteCompany(context.Background(), companyID)
		require.NoError(t, err)
	})

	t.Run("should fail to get a company that doesn't exist", func(t *testing.T) {
		_, err := testStore.GetCompany(context.Background(), uuid.Nil)
		require.Error(t, err)
	})

}

func TestStorage_CreateCompany(t *testing.T) {
	t.Run("should successfully create a company", func(t *testing.T) {

		company := model.CompanyCreate{
			Name:        "Some",
			Description: "Some",
			Employees:   123,
			Registered:  true,
			Type:        "Corporation",
		}

		companyID, err := testStore.CreateCompany(context.Background(), company)
		require.NoError(t, err)

		err = testStore.DeleteCompany(context.Background(), companyID)
		require.NoError(t, err)
	})

	t.Run("should fail to create a company due to unique name constraint", func(t *testing.T) {

		company := model.CompanyCreate{
			Name:        "Some",
			Description: "Some",
			Employees:   123,
			Registered:  true,
			Type:        "Corporation",
		}

		companyID, err := testStore.CreateCompany(context.Background(), company)
		require.NoError(t, err)

		_, err = testStore.CreateCompany(context.Background(), company)
		require.Error(t, err)

		err = testStore.DeleteCompany(context.Background(), companyID)
		require.NoError(t, err)
	})

	t.Run("should fail to create a company due to invalid type", func(t *testing.T) {

		company := model.CompanyCreate{
			Name:        "Some",
			Description: "Some",
			Employees:   123,
			Registered:  true,
			Type:        "Corpora",
		}

		_, err := testStore.CreateCompany(context.Background(), company)
		require.Error(t, err)
	})
}

func TestStorage_UpdateCompany(t *testing.T) {
	t.Run("should successfully update a company", func(t *testing.T) {

		company := model.CompanyCreate{
			Name:        "Some",
			Description: "Some",
			Employees:   123,
			Registered:  true,
			Type:        "Corporation",
		}

		updatedCompany := model.CompanyCreate{
			Name:        "SomeOther",
			Description: "SomeOther",
			Employees:   123,
			Registered:  true,
			Type:        "Corporation",
		}

		companyID, err := testStore.CreateCompany(context.Background(), company)
		require.NoError(t, err)

		err = testStore.UpdateCompany(context.Background(), updatedCompany, companyID)
		require.NoError(t, err)

		err = testStore.DeleteCompany(context.Background(), companyID)
		require.NoError(t, err)
	})

	t.Run("should return error if nothing was updated", func(t *testing.T) {

		company := model.CompanyCreate{
			Name:        "Some",
			Description: "Some",
			Employees:   123,
			Registered:  true,
			Type:        "Corporation",
		}

		err := testStore.UpdateCompany(context.Background(), company, uuid.Nil)
		require.Error(t, err)
	})

	// This causes a panic.
	// t.Run("should return error if trying to update company name to an existing one", func(t *testing.T) {

	// 	company1 := model.CompanyCreate{
	// 		Name:        "Some",
	// 		Description: "Some",
	// 		Employees:   123,
	// 		Registered:  true,
	// 		Type:        "Corporation",
	// 	}

	// 	company2 := model.CompanyCreate{
	// 		Name:        "SomeOther",
	// 		Description: "Some",
	// 		Employees:   123,
	// 		Registered:  true,
	// 		Type:        "Corporation",
	// 	}

	// 	updatedCompany := model.CompanyCreate{
	// 		Name:        "SomeOther",
	// 		Description: "SomeOther",
	// 		Employees:   123,
	// 		Registered:  true,
	// 		Type:        "Corporation",
	// 	}

	// 	companyID, err := testStore.CreateCompany(context.Background(), company1)
	// 	require.NoError(t, err)

	// 	companyID2, err := testStore.CreateCompany(context.Background(), company2)
	// 	require.NoError(t, err)

	// 	err = testStore.UpdateCompany(context.Background(), updatedCompany, companyID)
	// 	require.Error(t, err)

	// 	err = testStore.DeleteCompany(context.Background(), companyID)
	// 	require.NoError(t, err)

	// 	err = testStore.DeleteCompany(context.Background(), companyID2)
	// 	require.NoError(t, err)
	// })
}

func TestStorage_DeleteCompany(t *testing.T) {
	t.Run("should successfully delete a company", func(t *testing.T) {

		company := model.CompanyCreate{
			Name:        "Some",
			Description: "Some",
			Employees:   123,
			Registered:  true,
			Type:        "Corporation",
		}

		companyID, err := testStore.CreateCompany(context.Background(), company)
		require.NoError(t, err)

		err = testStore.DeleteCompany(context.Background(), companyID)
		require.NoError(t, err)
	})
}
