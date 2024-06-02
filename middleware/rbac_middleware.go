package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RoleBasedAccess middleware checks if the user's role meets the required role for access.
func RoleBasedAccess(next echo.HandlerFunc, requiredRole string) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole := c.Get("role").(string)

		roleHierarchy := map[string]int{
			"user":       1,
			"supervisor": 2,
			"manager":    3,
		}

		if roleHierarchy[userRole] < roleHierarchy[requiredRole] {
			return echo.NewHTTPError(http.StatusForbidden, "You don't have the necessary permissions to access this resource.")
		}

		return next(c)
	}
}
