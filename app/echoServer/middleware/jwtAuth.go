package middlewarex

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			raw := c.Request().Header.Get("Authorization")
			if raw == "" || !strings.HasPrefix(raw, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthenticated: missing bearer token"})
			}
			token := strings.TrimPrefix(raw, "Bearer ")
			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthenticated: invalid token"})
			}
			c.Set("claims", claims)
			return next(c)
		}
	}
}
