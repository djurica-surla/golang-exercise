package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/djurica-surla/golang-exercise/internal/mock"
	"github.com/djurica-surla/golang-exercise/internal/model"
	transporthttp "github.com/djurica-surla/golang-exercise/internal/transport/http"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

var (
	cookie = &http.Cookie{Name: "accessToken", Value: "eyb4dasVBSd"}
	errFoo = errors.New("some-error")
)

func TestHandler_GetCompany(t *testing.T) {
	t.Run("should successfully get a company", func(t *testing.T) {
		var (
			ctrl            = gomock.NewController(t)
			companyServicer = mock.NewMockCompanyServicer(ctrl)
			tokenServicer   = mock.NewMockTokenServicer(ctrl)
		)

		companyID, _ := uuid.NewUUID()
		companyPath := fmt.Sprintf("/companies/%s", companyID)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, companyPath, nil)
		req.AddCookie(cookie)

		gomock.InOrder(
			tokenServicer.EXPECT().VerifyAccessToken(cookie.Value).Return(nil),
			companyServicer.EXPECT().GetCompanyByID(gomock.Any(), companyID).Return(model.Company{}, nil),
		)

		handler := transporthttp.NewCompanyHandler(companyServicer, tokenServicer)
		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("should fail when company id is invalid", func(t *testing.T) {
		var (
			ctrl            = gomock.NewController(t)
			companyServicer = mock.NewMockCompanyServicer(ctrl)
			tokenServicer   = mock.NewMockTokenServicer(ctrl)
		)

		companyPath := fmt.Sprintf("/companies/%s", "1234")

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, companyPath, nil)
		req.AddCookie(cookie)

		gomock.InOrder(
			tokenServicer.EXPECT().VerifyAccessToken(cookie.Value).Return(nil),
		)

		handler := transporthttp.NewCompanyHandler(companyServicer, tokenServicer)
		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("should fail when service fails", func(t *testing.T) {
		var (
			ctrl            = gomock.NewController(t)
			companyServicer = mock.NewMockCompanyServicer(ctrl)
			tokenServicer   = mock.NewMockTokenServicer(ctrl)
		)

		companyID, _ := uuid.NewUUID()
		companyPath := fmt.Sprintf("/companies/%s", companyID)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, companyPath, nil)
		req.AddCookie(cookie)

		gomock.InOrder(
			tokenServicer.EXPECT().VerifyAccessToken(cookie.Value).Return(nil),
			companyServicer.EXPECT().GetCompanyByID(gomock.Any(), companyID).Return(model.Company{}, errFoo),
		)

		handler := transporthttp.NewCompanyHandler(companyServicer, tokenServicer)
		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func TestHandler_CreateCompany(t *testing.T) {
	t.Run("should successfully create a company", func(t *testing.T) {
		var (
			ctrl            = gomock.NewController(t)
			companyServicer = mock.NewMockCompanyServicer(ctrl)
			tokenServicer   = mock.NewMockTokenServicer(ctrl)
		)

		company := model.CompanyCreate{
			Name:       "test",
			Employees:  23,
			Registered: true,
			Type:       "Corporation",
		}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(company)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/companies", &b)
		req.AddCookie(cookie)

		accountID, _ := uuid.NewUUID()

		gomock.InOrder(
			tokenServicer.EXPECT().VerifyAccessToken(cookie.Value).Return(nil),
			companyServicer.EXPECT().CreateCompany(gomock.Any(), company).Return(accountID, nil),
		)

		handler := transporthttp.NewCompanyHandler(companyServicer, tokenServicer)
		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("should fail when fields are missing", func(t *testing.T) {
		var (
			ctrl            = gomock.NewController(t)
			companyServicer = mock.NewMockCompanyServicer(ctrl)
			tokenServicer   = mock.NewMockTokenServicer(ctrl)
		)

		company := model.CompanyCreate{}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(company)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/companies", &b)
		req.AddCookie(cookie)

		gomock.InOrder(
			tokenServicer.EXPECT().VerifyAccessToken(cookie.Value).Return(nil),
		)

		handler := transporthttp.NewCompanyHandler(companyServicer, tokenServicer)
		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusBadRequest, res.Code)
	})
	t.Run("should fail when company type is wrong", func(t *testing.T) {
		var (
			ctrl            = gomock.NewController(t)
			companyServicer = mock.NewMockCompanyServicer(ctrl)
			tokenServicer   = mock.NewMockTokenServicer(ctrl)
		)

		company := model.CompanyCreate{
			Name:       "test",
			Employees:  23,
			Registered: true,
			Type:       "Corp",
		}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(company)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/companies", &b)
		req.AddCookie(cookie)

		gomock.InOrder(
			tokenServicer.EXPECT().VerifyAccessToken(cookie.Value).Return(nil),
		)

		handler := transporthttp.NewCompanyHandler(companyServicer, tokenServicer)
		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("should fail when service fails", func(t *testing.T) {
		var (
			ctrl            = gomock.NewController(t)
			companyServicer = mock.NewMockCompanyServicer(ctrl)
			tokenServicer   = mock.NewMockTokenServicer(ctrl)
		)

		company := model.CompanyCreate{
			Name:       "test",
			Employees:  23,
			Registered: true,
			Type:       "Corporation",
		}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(company)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/companies", &b)
		req.AddCookie(cookie)

		gomock.InOrder(
			tokenServicer.EXPECT().VerifyAccessToken(cookie.Value).Return(nil),
			companyServicer.EXPECT().CreateCompany(gomock.Any(), company).Return(uuid.Nil, errFoo),
		)

		handler := transporthttp.NewCompanyHandler(companyServicer, tokenServicer)
		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusInternalServerError, res.Code)
	})
}
