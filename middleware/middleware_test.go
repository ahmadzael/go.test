package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRoleBasedAccess(t *testing.T) {
	e := echo.New()

	// Handler to be protected by middleware
	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Access granted")
	}

	tests := []struct {
		name         string
		role         string
		requiredRole string
		expectedCode int
	}{
		{
			name:         "User trying to access supervisor route",
			role:         "user",
			requiredRole: "supervisor",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "Supervisor trying to access manager route",
			role:         "supervisor",
			requiredRole: "manager",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "Manager accessing supervisor route",
			role:         "manager",
			requiredRole: "supervisor",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Manager accessing user route",
			role:         "manager",
			requiredRole: "user",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Supervisor accessing supervisor route",
			role:         "supervisor",
			requiredRole: "supervisor",
			expectedCode: http.StatusOK,
		},
		{
			name:         "User accessing user route",
			role:         "user",
			requiredRole: "user",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("role", tt.role)

			// Apply middleware
			mw := RoleBasedAccess(handler, tt.requiredRole)
			mw(c)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
