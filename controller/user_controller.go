package controller

import "github.com/gofiber/fiber/v2"

type UserController interface {
	RegisterCtrl(c *fiber.Ctx) error
	LoginCtrl(c *fiber.Ctx) error
	ProfileCtrl(c *fiber.Ctx) error
	UpdateProfileCtrl(c *fiber.Ctx) error
	UpdatePasswordCtrl(c *fiber.Ctx) error
	UpdateImgCtrl(c *fiber.Ctx) error
	GetAllFavoriteCtrl(c *fiber.Ctx) error
}
