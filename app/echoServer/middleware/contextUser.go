package middlewarex

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RequireUserID(c echo.Context) (uint, error) {
	claims, _ := c.Get("claims").(map[string]any)
	if claims == nil {
		return 0, c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthenticated"})
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthenticated"})
	}
	return uint(sub), nil
}
