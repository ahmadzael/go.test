package middleware

import (
	"net/http"

	util "go.test/utils"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing or invalid token"})
		}

		tokenString := authHeader[len("Bearer "):]

		claims, err := util.ParseJWT(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		return next(c)
	}
}
