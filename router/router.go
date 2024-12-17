package router

import (
	"meals-app/controller"
	"meals-app/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRouter(app *fiber.App, db *gorm.DB, userCtrl controller.UserController, mealCtrl controller.MealController) {
	api := app.Group("/api")
	api.Post("/register", userCtrl.RegisterCtrl)
	api.Post("/login", userCtrl.LoginCtrl)

	user := api.Group("/users")
	user.Get("/", middleware.Protected(db), userCtrl.ProfileCtrl)
	user.Put("/", middleware.Protected(db), userCtrl.UpdateProfileCtrl)
	user.Put("/change-password", middleware.Protected(db), userCtrl.UpdatePasswordCtrl)
	user.Put("/image", middleware.Protected(db), userCtrl.UpdateImgCtrl)
	user.Get("/favorites", middleware.Protected(db), userCtrl.GetAllFavoriteCtrl)

	meal := api.Group("/meals")
	meal.Post("/", middleware.Protected(db), mealCtrl.CreateMealCtrl)
	meal.Get("/", middleware.Protected(db), mealCtrl.GetAllMealCtrl)
	meal.Get("/:id", middleware.Protected(db), mealCtrl.GetMealByIDCtrl)
	meal.Put("/:id", middleware.Protected(db), mealCtrl.UpdateMealCtrl)
	meal.Put("/:id/image", middleware.Protected(db), mealCtrl.UpdateMealImageCtrl)
	meal.Delete("/:id", middleware.Protected(db), mealCtrl.DeleteMealCtrl)
	meal.Post("/:id/favorites", middleware.Protected(db), mealCtrl.AddToFavoriteCtrl)
	meal.Delete("/:id/favorites", middleware.Protected(db), mealCtrl.DeleteFromFavoriteCtrl)
}
