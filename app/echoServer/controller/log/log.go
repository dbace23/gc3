// @Summary Create exercise log
// @Tags Logs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body object{ exercise_id=int, weight=int, repition_count=int, set_count=int, logged_at=string } true "Log payload"
// @Success 201 {object} map[string]any
// @Failure 403 {object} map[string]any
// @Router /api/logs [post]

// @Summary List exercise logs
// @Tags Logs
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]any
// @Router /api/logs [get]

package logcontroller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	mw "gym/app/echoServer/middleware"
	logsrv "gym/service/log"
)

type Controller struct {
	Svc logsrv.Service
}

// POST /api/logs
func (c *Controller) Create(ctx echo.Context) error {
	uid, errResp := mw.RequireUserID(ctx)
	if errResp != nil {
		return errResp
	}

	var req struct {
		ExerciseID uint      `json:"exercise_id" validate:"required"`
		SetCount   int       `json:"set_count" validate:"required,gt=0"`
		RepCount   int       `json:"repition_count" validate:"required,gt=0"`
		Weight     int       `json:"weight" validate:"required,gte=0"`
		LoggedAt   time.Time `json:"logged_at" validate:"required"`
	}
	if err := mw.BindAndValidate(ctx, &req); err != nil {
		return err
	}

	err := c.Svc.Create(uid, req.ExerciseID, req.SetCount, req.RepCount, req.Weight, req.LoggedAt)
	if err != nil {
		if err.Error() == "forbidden: exercise not owned by user" {
			return ctx.JSON(http.StatusForbidden, echo.Map{"message": err.Error()})
		}
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, echo.Map{"message": "log created"})
}

// GET /api/logs
func (c *Controller) List(ctx echo.Context) error {
	uid, errResp := mw.RequireUserID(ctx)
	if errResp != nil {
		return errResp
	}

	data, err := c.Svc.ListMy(uid)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"message": "ok", "data": data})
}
