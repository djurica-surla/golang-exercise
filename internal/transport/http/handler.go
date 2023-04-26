package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	httpErrors "github.com/djurica-surla/golang-exercise/internal/errors"
	"github.com/djurica-surla/golang-exercise/internal/middleware"
	"github.com/djurica-surla/golang-exercise/internal/model"
)

// RegisterRoutes links routes with the handler.
func (h *CompanyHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", middleware.HandleError(h.Login)).Methods(http.MethodPost)
	router.HandleFunc("/companies/{id}",
		middleware.Chain(middleware.HandleError(h.GetCompany), middleware.Authenticate(h.tokenService))).Methods(http.MethodGet)
	router.HandleFunc("/companies",
		middleware.Chain(middleware.HandleError(h.CreateCompany), middleware.Authenticate(h.tokenService))).Methods(http.MethodPost)
	router.HandleFunc("/companies/{id}",
		middleware.Chain(middleware.HandleError(h.UpdateCompany), middleware.Authenticate(h.tokenService))).Methods(http.MethodPatch)
	router.HandleFunc("/companies/{id}",
		middleware.Chain(middleware.HandleError(h.DeleteCompany), middleware.Authenticate(h.tokenService))).Methods(http.MethodDelete)
}

// CompanyServicer represents necessary company service implementation for company handler.
type CompanyServicer interface {
	GetCompanyByID(ctx context.Context, companyID uuid.UUID) (model.Company, error)
	CreateCompany(ctx context.Context, company model.CompanyCreate) (uuid.UUID, error)
	UpdateCompany(ctx context.Context, company model.CompanyCreate, companyID uuid.UUID) error
	DeleteCompany(ctx context.Context, companyID uuid.UUID) error
}

// TokenServicer represents necessary token service implementation for company handler.
type TokenServicer interface {
	CreateAccessToken(accessToken model.Login) (string, error)
	VerifyAccessToken(token string) error
}

// CompanyHandler handles http requests for companies.
type CompanyHandler struct {
	companyService CompanyServicer
	tokenService   TokenServicer
}

// NewCompanyHandler creates a new instance of a company handler.
func NewCompanyHandler(companyService CompanyServicer, tokenService TokenServicer) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
		tokenService:   tokenService,
	}
}

// GetCompanyById handles retrieveing company by the id.
func (h *CompanyHandler) Login(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var loginInfo model.Login

	err := json.NewDecoder(r.Body).Decode(&loginInfo)
	if err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	err = h.ValidateModel(loginInfo)
	if err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	accessToken, err := h.tokenService.CreateAccessToken(loginInfo)
	if err != nil {
		return nil, httpErrors.NewInternalServerErrorWrapped(err, "failed to create access token")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		HttpOnly: true,
		Expires:  time.Now().Add(5 * time.Minute),
		Path:     "/",
	})

	return "Logged in! Cookie has been set!", nil
}

// GetCompanyById handles retrieveing company by the id.
func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uuidString := mux.Vars(r)["id"]
	companyID, err := uuid.Parse(uuidString)
	if err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	res, err := h.companyService.GetCompanyByID(r.Context(), companyID)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}

	return res, nil
}

// CreateCompany handles creating a company.
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var company model.CompanyCreate

	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	err = h.ValidateModel(company)
	if err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	ok, validTypes := company.IsValidType()
	if !ok {
		return nil, httpErrors.NewBadRequestError(fmt.Errorf("invalid company type, must be one of these: %s", validTypes))
	}

	companyID, err := h.companyService.CreateCompany(r.Context(), company)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}

	return companyID, nil
}

// UpdateCompany handles updating a company.
func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var company model.CompanyCreate

	uuidString := mux.Vars(r)["id"]
	companyID, err := uuid.Parse(uuidString)
	if err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	err = json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	err = h.ValidateModel(company)
	if err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	ok, validTypes := company.IsValidType()
	if !ok {
		return nil, httpErrors.NewBadRequestError(fmt.Errorf("invalid company type, must be one of these: %s", validTypes))
	}

	err = h.companyService.UpdateCompany(r.Context(), company, companyID)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}

	return "Successfully updated company", nil
}

// DeleteCompany handles deleting a company.
func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uuidString := mux.Vars(r)["id"]
	companyID, err := uuid.Parse(uuidString)
	if err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	err = h.companyService.DeleteCompany(r.Context(), companyID)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}

	return "Successfully deleted a company", nil
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
