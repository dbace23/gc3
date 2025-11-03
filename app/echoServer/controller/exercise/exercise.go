// @Summary Create exercise
// @Tags Exercises
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body object{ workout_id=int, name=string, description=string } true "Exercise payload"
// @Success 201 {object} map[string]any
// @Failure 403 {object} map[string]any
// @Router /api/exercises [post]

// @Summary Delete exercise
// @Tags Exercises
// @Security BearerAuth
// @Param id path int true "Exercise ID"
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 403 {object} map[string]any
// @Router /api/exercises/{id} [delete]

package exercise

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	mw "gym/app/echoServer/middleware"
	exercisesrv "gym/service/exercise"
)

type Controller struct {
	Svc exercisesrv.Service
}

// POST /api/exercises
func (c *Controller) Create(ctx echo.Context) error {
	uid, errResp := mw.RequireUserID(ctx)
	if errResp != nil {
		return errResp
	}

	var req struct {
		WorkoutID   uint   `json:"workout_id" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}
	if err := mw.BindAndValidate(ctx, &req); err != nil {
		return err
	}

	e, err := c.Svc.Create(uid, req.WorkoutID, req.Name, req.Description)
	if err != nil {
		if err.Error() == "forbidden: workout not owned by user" {
			return ctx.JSON(http.StatusForbidden, echo.Map{"message": err.Error()})
		}
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, echo.Map{"message": "exercise created", "data": e})
}

// DELETE /api/exercises/:id
func (c *Controller) Delete(ctx echo.Context) error {
	uid, errResp := mw.RequireUserID(ctx)
	if errResp != nil {
		return errResp
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := c.Svc.Delete(uid, uint(id))
	if err != nil {
		if err.Error() == "forbidden: not owner" {
			return ctx.JSON(http.StatusForbidden, echo.Map{"message": err.Error()})
		}
		return ctx.JSON(http.StatusNotFound, echo.Map{"message": "exercise not found"})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"message": "exercise deleted"})
}
