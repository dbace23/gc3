package middlewarex

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

func BindAndValidate(c echo.Context, dst any) error {
	if err := c.Bind(dst); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}
	if err := validate.Struct(dst); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			out := echo.Map{}
			for _, fe := range ve {
				out[fe.Field()] = fe.Tag()
			}
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"message": "validation error",
				"errors":  out,
			})
		}
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"error":   err.Error(),
		})
	}
	return nil
}
