package echoServer

import (
	"github.com/labstack/echo/v4"

	ctrlAuth "gym/app/echoServer/controller/auth"
	ctrlExercise "gym/app/echoServer/controller/exercise"
	ctrlLog "gym/app/echoServer/controller/log"
	ctrlUser "gym/app/echoServer/controller/user"
	ctrlWorkout "gym/app/echoServer/controller/workout"
	mw "gym/app/echoServer/middleware"
	"gym/config"
	authsvc "gym/service/auth"
	usersvc "gym/service/user"
	"gym/util/rapid"
)

type Controllers struct {
	Auth     ctrlAuth.Controller
	User     ctrlUser.Controller
	Workout  ctrlWorkout.Controller
	Exercise ctrlExercise.Controller
	Log      ctrlLog.Controller
}

func NewAuthController(svc authsvc.Service, jwtSecret string) ctrlAuth.Controller {
	return ctrlAuth.Controller{Svc: svc, JWTSecret: jwtSecret}
}
func NewUserController(svc usersvc.Service, bmi *rapid.BMIClient) ctrlUser.Controller {
	return ctrlUser.Controller{Svc: svc, BMI: bmi}
}

func Mount(e *echo.Echo, cfg *config.Config, c Controllers) {
	api := e.Group("/api")

	// public
	api.POST("/users/register", c.Auth.Register)
	api.POST("/users/login", c.Auth.Login)

	// protected
	auth := api.Group("")
	auth.Use(mw.JWTAuth(cfg.JWTSecret))

	auth.GET("/users", c.User.Me)
	auth.GET("/workouts", c.Workout.List)
	auth.GET("/workouts/:id", c.Workout.Get)
	auth.POST("/workouts", c.Workout.Create)
	auth.PUT("/workouts/:id", c.Workout.Update)
	auth.DELETE("/workouts/:id", c.Workout.Delete)

	auth.POST("/exercises", c.Exercise.Create)
	auth.DELETE("/exercises/:id", c.Exercise.Delete)

	auth.POST("/logs", c.Log.Create)
	auth.GET("/logs", c.Log.List)

}
