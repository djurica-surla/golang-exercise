package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	httpError "github.com/djurica-surla/golang-exercise/internal/errors"
	"github.com/djurica-surla/golang-exercise/internal/model"
)

// RegisterRoutes links routes with the handler.
func (h *CompanyHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/companies/{id}", handleError(h.GetCompany)).Methods(http.MethodGet)
	router.HandleFunc("/companies", handleError(h.CreateCompany)).Methods(http.MethodPost)
	router.HandleFunc("/companies/{id}", handleError(h.UpdateCompany)).Methods(http.MethodPatch)
	router.HandleFunc("/companies/{id}", handleError(h.DeleteCompany)).Methods(http.MethodDelete)
}

// CompanyServicer represents necessary company service implementation for company handler.
type CompanyServicer interface {
	GetCompanyByID(ctx context.Context, companyID uuid.UUID) (model.Company, error)
	CreateCompany(ctx context.Context, company model.CompanyCreate) (uuid.UUID, error)
	UpdateCompany(ctx context.Context, company model.CompanyCreate, companyID uuid.UUID) error
	DeleteCompany(ctx context.Context, companyID uuid.UUID) error
}

// Custom handler func which allows returning error so that they can be handled in once place.
type handlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

// CompanyHandler handles http requests for companies.
type CompanyHandler struct {
	companyService CompanyServicer
}

// NewCompanyHandler creates a new instance of a company handler.
func NewCompanyHandler(companyService CompanyServicer) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

// Message will be returned as a response.
type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

// GetCompanyById handles retrieveing company by the id.
func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uuidString := mux.Vars(r)["id"]
	companyID, err := uuid.Parse(uuidString)
	if err != nil {
		return nil, httpError.NewBadRequestError(err)
	}

	res, err := h.companyService.GetCompanyByID(r.Context(), companyID)
	if err != nil {
		return nil, httpError.NewInternalServerError(err)
	}

	return res, nil
}

// CreateCompany handles creating a company.
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var company model.CompanyCreate

	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		return nil, httpError.NewBadRequestError(err)
	}

	err = h.ValidateModel(company)
	if err != nil {
		return nil, httpError.NewBadRequestError(err)
	}

	ok, validTypes := company.IsValidType()
	if !ok {
		return nil, httpError.NewBadRequestError(fmt.Errorf("invalid company type, must be one of these: %s", validTypes))
	}

	companyID, err := h.companyService.CreateCompany(r.Context(), company)
	if err != nil {
		return nil, httpError.NewInternalServerError(err)
	}

	return companyID, nil
}

// UpdateCompany handles updating a company.
func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var company model.CompanyCreate

	uuidString := mux.Vars(r)["id"]
	companyID, err := uuid.Parse(uuidString)
	if err != nil {
		return nil, httpError.NewBadRequestError(err)
	}

	err = json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		return nil, httpError.NewBadRequestError(err)
	}

	err = h.ValidateModel(company)
	if err != nil {
		return nil, httpError.NewBadRequestError(err)
	}

	ok, validTypes := company.IsValidType()
	if !ok {
		return nil, httpError.NewBadRequestError(fmt.Errorf("invalid company type, must be one of these: %s", validTypes))
	}

	err = h.companyService.UpdateCompany(r.Context(), company, companyID)
	if err != nil {
		return nil, httpError.NewInternalServerError(err)
	}

	return "Successfully updated company", nil
}

// DeleteCompany handles deleting a company.
func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uuidString := mux.Vars(r)["id"]
	companyID, err := uuid.Parse(uuidString)
	if err != nil {
		return nil, httpError.NewBadRequestError(err)
	}

	err = h.companyService.DeleteCompany(r.Context(), companyID)
	if err != nil {
		return nil, httpError.NewInternalServerError(err)
	}

	return "Successfully deleted a company", nil
}

// Handle error wraps around the custom handler func and resolves error type to a proper response.
func handleError(f handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := f(w, r)
		if err != nil {
			if httpError.IsUnauthorizedError(err) {
				w.WriteHeader(http.StatusUnauthorized)
			} else if httpError.IsBadRequestError(err) {
				w.WriteHeader(http.StatusBadRequest)
			} else if httpError.IsNotFoundError(err) {
				w.WriteHeader(http.StatusNotFound)
			} else if httpError.IsForbiddenError(err) {
				w.WriteHeader(http.StatusForbidden)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			errorResponse := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// Validate validates a model struct using validator package.
func (h *CompanyHandler) ValidateModel(data interface{}) error {
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		switch err := err.(type) {
		case validator.ValidationErrors:
			var fields []string
			for _, field := range err {
				fields = append(fields, strings.ToLower(field.Field()))
			}
			return fmt.Errorf("field validation error, missing fields: %v", fields)
		default:
			return fmt.Errorf("field validation error: %w", err)
		}
	}
	return nil
}
