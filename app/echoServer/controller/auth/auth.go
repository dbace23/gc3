// @Summary Register user
// @Description Creates a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body registerReq true "User registration data"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /api/users/register [post]

// @Summary Login user
// @Description Logs user in and returns JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body loginReq true "User login credentials"
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/users/login [post]
package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	mw "gym/app/echoServer/middleware"
	authsvc "gym/service/auth"
)

type Controller struct {
	Svc       authsvc.Service
	JWTSecret string
}

type registerReq struct {
	FullName string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
	WeightKG int    `json:"weight" validate:"required,gt=0"`
	HeightCM int    `json:"height" validate:"required,gt=0"`
}

func (c *Controller) Register(ctx echo.Context) error {
	var req registerReq
	if err := mw.BindAndValidate(ctx, &req); err != nil {
		return err
	}

	res, err := c.Svc.Register(authsvc.RegisterParams{
		FullName: req.FullName, Email: req.Email, Username: req.Username,
		Password: req.Password, WeightKG: req.WeightKG, HeightCM: req.HeightCM,
	})
	if err != nil {
		if err.Error() == "email already registered" || err.Error() == "username already taken" {
			return ctx.JSON(http.StatusConflict, echo.Map{"message": err.Error()})
		}
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, echo.Map{
		"message": "user registered",
		"data":    res,
	})
}

type loginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (c *Controller) Login(ctx echo.Context) error {
	var req loginReq
	if err := mw.BindAndValidate(ctx, &req); err != nil {
		return err
	}

	res, err := c.Svc.Login(authsvc.LoginParams{
		Email: req.Email, Password: req.Password,
	}, c.JWTSecret)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"message": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "login success",
		"data":    res,
	})
}
