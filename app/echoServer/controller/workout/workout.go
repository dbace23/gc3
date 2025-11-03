// @Summary List workouts
// @Tags Workouts
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]any
// @Router /api/workouts [get]

// @Summary Get workout by ID
// @Tags Workouts
// @Security BearerAuth
// @Param id path int true "Workout ID"
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /api/workouts/{id} [get]

// @Summary Create new workout
// @Tags Workouts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body object{ name=string, description=string } true "Workout payload"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /api/workouts [post]

// @Summary Update workout
// @Tags Workouts
// @Security BearerAuth
// @Param id path int true "Workout ID"
// @Accept json
// @Produce json
// @Param body body object{ name=string, description=string } true "Updated workout"
// @Success 200 {object} map[string]any
// @Failure 403 {object} map[string]any
// @Router /api/workouts/{id} [put]

// @Summary Delete workout
// @Tags Workouts
// @Security BearerAuth
// @Param id path int true "Workout ID"
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 403 {object} map[string]any
// @Router /api/workouts/{id} [delete]

package workout

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	mw "gym/app/echoServer/middleware"
	workoutsrv "gym/service/workout"
)

type Controller struct {
	Svc workoutsrv.Service
}

// GET /api/workouts
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

// GET /api/workouts/:id
func (c *Controller) Get(ctx echo.Context) error {
	uid, errResp := mw.RequireUserID(ctx)
	if errResp != nil {
		return errResp
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	data, err := c.Svc.GetMy(uid, uint(id))
	if err != nil {
		if err.Error() == "forbidden: not owner" {
			return ctx.JSON(http.StatusForbidden, echo.Map{"message": err.Error()})
		}
		return ctx.JSON(http.StatusNotFound, echo.Map{"message": "workout not found"})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"message": "ok", "data": data})
}

// POST /api/workouts
func (c *Controller) Create(ctx echo.Context) error {
	uid, errResp := mw.RequireUserID(ctx)
	if errResp != nil {
		return errResp
	}

	var req struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}
	if err := mw.BindAndValidate(ctx, &req); err != nil {
		return err
	}

	w, err := c.Svc.Create(uid, req.Name, req.Description)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, echo.Map{"message": "workout created", "data": w})
}

// PUT /api/workouts/:id
func (c *Controller) Update(ctx echo.Context) error {
	uid, errResp := mw.RequireUserID(ctx)
	if errResp != nil {
		return errResp
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	var req struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}
	if err := mw.BindAndValidate(ctx, &req); err != nil {
		return err
	}

	w, err := c.Svc.Update(uid, uint(id), req.Name, req.Description)
	if err != nil {
		if err.Error() == "forbidden: not owner" {
			return ctx.JSON(http.StatusForbidden, echo.Map{"message": err.Error()})
		}
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"message": "workout updated", "data": w})
}

// DELETE /api/workouts/:id
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
		return ctx.JSON(http.StatusNotFound, echo.Map{"message": "workout not found"})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"message": "workout deleted"})
}
