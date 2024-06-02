package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RBAC(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("role").(string)
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					return next(c)
				}
			}
			return c.JSON(http.StatusForbidden, map[string]string{"message": "Forbidden"})
		}
	}
}
