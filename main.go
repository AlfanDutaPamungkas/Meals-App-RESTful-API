package main

import (
	"log"
	"meals-app/config"
	"meals-app/controller"
	"meals-app/database"
	"meals-app/helper"
	"meals-app/router"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load(".env")
	helper.PanicError(envErr)

	cld := config.NewCloudinary()
	validate := validator.New()
	db := database.DatabaseInit()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":   fiber.StatusInternalServerError,
				"status": "Internal Server Error",
				"data": fiber.Map{
					"error": err.Error(),
				},
			})
		},
	})
	app.Use(recover.New())

	userController := controller.NewUserControllerImpl(db, validate, cld)
	mealController := controller.NewMealControllerImpl(db, validate, cld)

	router.SetupRouter(app, db, userController, mealController)

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
