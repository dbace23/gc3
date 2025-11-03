// @Summary Get current user
// @Description Returns authenticated user profile and BMI
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/users [get]

package user

import (
	"net/http"

	"github.com/labstack/echo/v4"

	mw "gym/app/echoServer/middleware"
	usersvc "gym/service/user"
	"gym/util/rapid"
)

type Controller struct {
	Svc usersvc.Service
	BMI *rapid.BMIClient
}

func (c *Controller) Me(ctx echo.Context) error {
	uid, errResp := mw.RequireUserID(ctx)
	if errResp != nil {
		return errResp
	}

	u, err := c.Svc.FindByID(uid)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"message": "user not found"})
	}

	bmi, err := c.BMI.Calculate(u.WeightKG, u.HeightCM)
	if err != nil {
		return ctx.JSON(http.StatusOK, echo.Map{
			"message": "ok (BMI lookup failed, returning profile only)",
			"data":    echo.Map{"user": u},
			"error":   err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "ok",
		"data":    echo.Map{"user": u, "bmi": bmi},
	})
}
