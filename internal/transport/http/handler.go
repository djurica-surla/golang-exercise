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

	"github.com/djurica-surla/golang-exercise/internal/model"
)

// RegisterRoutes links routes with the handler.
func (h *CompanyHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/companies/{id}", h.GetCompany()).Methods(http.MethodGet)
	router.HandleFunc("/companies", h.CreateCompany()).Methods(http.MethodPost)
	router.HandleFunc("/companies/{id}", h.UpdateCompany()).Methods(http.MethodPatch)
	router.HandleFunc("/companies/{id}", h.DeleteCompany()).Methods(http.MethodDelete)
}

// CompanyServicer represents necessary company service implementation for company handler.
type CompanyServicer interface {
	GetCompanyByID(ctx context.Context, companyID uuid.UUID) (model.Company, error)
	CreateCompany(ctx context.Context, company model.CompanyCreate) (uuid.UUID, error)
	UpdateCompany(ctx context.Context, company model.CompanyCreate, companyID uuid.UUID) error
	DeleteCompany(ctx context.Context, companyID uuid.UUID) error
}

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

// GetCompanyById handles retrieveing company by the id.
func (h *CompanyHandler) GetCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidString := mux.Vars(r)["id"]
		companyID, err := uuid.Parse(uuidString)
		if err != nil {
			h.encodeErrorWithStatus404(err, w)
		}

		res, err := h.companyService.GetCompanyByID(r.Context(), companyID)
		if err != nil {
			h.encodeErrorWithStatus500(err, w)
			return
		}

		json.NewEncoder(w).Encode(res)
	}
}

// CreateCompany handles creating a company.
func (h *CompanyHandler) CreateCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var company model.CompanyCreate

		err := json.NewDecoder(r.Body).Decode(&company)
		if err != nil {
			h.encodeErrorWithStatus404(err, w)
			return
		}

		err = h.ValidateModel(company)
		if err != nil {
			h.encodeErrorWithStatus404(err, w)
			return
		}

		ok, validTypes := company.IsValidType()
		if !ok {
			h.encodeErrorWithStatus404(fmt.Errorf("invalid company type, must be one of these: %s", validTypes), w)
			return
		}

		companyID, err := h.companyService.CreateCompany(r.Context(), company)
		if err != nil {
			h.encodeErrorWithStatus500(err, w)
			return
		}

		json.NewEncoder(w).Encode(fmt.Sprintf("Successfully created a company with the id: %s", companyID))
	}
}

// UpdateCompany handles updating a company.
func (h *CompanyHandler) UpdateCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var company model.CompanyCreate

		uuidString := mux.Vars(r)["id"]
		companyID, err := uuid.Parse(uuidString)
		if err != nil {
			h.encodeErrorWithStatus404(err, w)
		}

		err = json.NewDecoder(r.Body).Decode(&company)
		if err != nil {
			h.encodeErrorWithStatus404(err, w)
			return
		}

		err = h.ValidateModel(company)
		if err != nil {
			h.encodeErrorWithStatus404(err, w)
			return
		}

		ok, validTypes := company.IsValidType()
		if !ok {
			h.encodeErrorWithStatus404(fmt.Errorf("invalid company type, must be one of these: %s", validTypes), w)
			return
		}

		err = h.companyService.UpdateCompany(r.Context(), company, companyID)
		if err != nil {
			h.encodeErrorWithStatus500(err, w)
			return
		}

		json.NewEncoder(w).Encode("Successfully updated a company")
	}
}

// DeleteCompany handles deleting a company.
func (h *CompanyHandler) DeleteCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		uuidString := mux.Vars(r)["id"]
		companyID, err := uuid.Parse(uuidString)
		if err != nil {
			h.encodeErrorWithStatus404(err, w)
		}

		err = h.companyService.DeleteCompany(r.Context(), companyID)
		if err != nil {
			h.encodeErrorWithStatus500(err, w)
			return
		}

		json.NewEncoder(w).Encode("Successfully deleted a company")
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

func (h *CompanyHandler) encodeErrorWithStatus500(err error, w http.ResponseWriter) {
	errorResponse := fmt.Sprintf("error: %s", err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(errorResponse))
}

func (h *CompanyHandler) encodeErrorWithStatus404(err error, w http.ResponseWriter) {
	errorResponse := fmt.Sprintf("error: %s", err.Error())
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(errorResponse))
}
