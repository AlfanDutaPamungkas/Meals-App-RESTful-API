package controller

import "github.com/gofiber/fiber/v2"

type MealController interface {
	CreateMealCtrl(c *fiber.Ctx) error
	GetAllMealCtrl(c *fiber.Ctx) error
	GetMealByIDCtrl(c *fiber.Ctx) error
	UpdateMealCtrl(c *fiber.Ctx) error
	UpdateMealImageCtrl(c *fiber.Ctx) error
	DeleteMealCtrl(c *fiber.Ctx) error
	AddToFavoriteCtrl(c *fiber.Ctx) error
	DeleteFromFavoriteCtrl(c *fiber.Ctx) error
}
