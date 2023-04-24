package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/djurica-surla/golang-exercise/internal/entity"
)

// RegisterRoutes links routes with the handler.
func (h *CompanyHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/companies/{id}", h.GetCompanyByID()).Methods(http.MethodGet)
}

// CompanyServicer represents necessary company service implementation for company handler.
type CompanyServicer interface {
	GetCompanyByID(ctx context.Context, companyID uuid.UUID) (entity.Company, error)
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
func (h *CompanyHandler) GetCompanyByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidString := mux.Vars(r)["id"]
		uuid, err := uuid.Parse(uuidString)
		if err != nil {
			h.encodeErrorWithStatus404(err, w)
		}

		res, err := h.companyService.GetCompanyByID(r.Context(), uuid)
		if err != nil {
			h.encodeErrorWithStatus500(err, w)
			return
		}

		json.NewEncoder(w).Encode(res)
	}
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
