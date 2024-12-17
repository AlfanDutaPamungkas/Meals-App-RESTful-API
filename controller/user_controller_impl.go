package controller

import (
	"errors"
	"meals-app/exception"
	"meals-app/helper"
	"meals-app/model/entity"
	"meals-app/model/web"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserControllerImpl struct {
	DB       *gorm.DB
	Validate *validator.Validate
	Cld      *cloudinary.Cloudinary
}

func NewUserControllerImpl(DB *gorm.DB, validate *validator.Validate, cld *cloudinary.Cloudinary) UserController {
	return &UserControllerImpl{
		DB:       DB,
		Validate: validate,
		Cld:      cld,
	}
}

func (controller *UserControllerImpl) RegisterCtrl(c *fiber.Ctx) error {
	request := new(web.UserRegisterReq)
	err := c.BodyParser(request)
	if err != nil {
		return exception.ErrorHandler(400, "BAD REQUEST", err)(c)
	}

	err = controller.Validate.Struct(request)
	if err != nil {
		return exception.ErrorHandler(422, "VALIDATION ERROR", err)(c)
	}

	hashedPassword, err := helper.HashPassword(request.Password)
	helper.PanicError(err)

	user := entity.User{
		Username: request.Username,
		Password: hashedPassword,
		Email:    request.Email,
		ImageUrl: "https://th.bing.com/th/id/OIP.R9HMSxN_IRyxw9-iE1usugAAAA?rs=1&pid=ImgDetMain",
	}

	err = controller.DB.Create(&user).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return exception.ErrorHandler(409, "DUPLICATE ENTRY", errors.New("email already exists"))(c)
		}

		helper.PanicError(err)
	}

	response := web.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Image:     user.ImageUrl,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   response,
	})
}

func (controller *UserControllerImpl) LoginCtrl(c *fiber.Ctx) error {
	request := new(web.LoginRequest)
	err := c.BodyParser(request)
	if err != nil {
		return exception.ErrorHandler(400, "BAD REQUEST", err)(c)
	}

	err = controller.Validate.Struct(request)
	if err != nil {
		return exception.ErrorHandler(422, "VALIDATION ERROR", err)(c)
	}

	user := entity.User{}
	err = controller.DB.Take(&user, "email = ?", request.Email).Error
	if err != nil {
		return exception.ErrorHandler(401, "INVALID CREDENTIALS", err)(c)
	}

	valid := helper.VerifyPassword(request.Password, user.Password)
	if !valid {
		return exception.ErrorHandler(401, "INVALID CREDENTIALS", errors.New("invalid credentials"))(c)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["user_id"] = user.ID
	claims["token_version"] = user.TokenVersion
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN_SECRET")))
	helper.PanicError(err)

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   t,
	})
}

func (controller *UserControllerImpl) ProfileCtrl(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(entity.User)

	response := web.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Image:     user.ImageUrl,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   response,
	})
}

func (controller *UserControllerImpl) UpdateProfileCtrl(c *fiber.Ctx) error {
	request := new(web.UserUpdateReq)
	err := c.BodyParser(request)
	if err != nil {
		return exception.ErrorHandler(400, "BAD REQUEST", err)(c)
	}

	err = controller.Validate.Struct(request)
	if err != nil {
		return exception.ErrorHandler(422, "VALIDATION ERROR", err)(c)
	}

	user := c.Locals("currentUser").(entity.User)

	if request.Email != "" {
		user.Email = request.Email
		user.TokenVersion++
	}

	if request.Username != "" {
		user.Username = request.Username
	}

	err = controller.DB.Save(&user).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return exception.ErrorHandler(409, "DUPLICATE ENTRY", errors.New("email already exists"))(c)
		}

		helper.PanicError(err)
	}

	response := web.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Image:     user.ImageUrl,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   response,
	})
}

func (controller *UserControllerImpl) UpdatePasswordCtrl(c *fiber.Ctx) error {
	request := new(web.ChangePassReq)
	err := c.BodyParser(request)
	if err != nil {
		return exception.ErrorHandler(400, "BAD REQUEST", err)(c)
	}

	err = controller.Validate.Struct(request)
	if err != nil {
		return exception.ErrorHandler(422, "VALIDATION ERROR", err)(c)
	}

	hash, err := helper.HashPassword(request.Password)
	helper.PanicError(err)

	user := c.Locals("currentUser").(entity.User)

	user.Password = hash
	user.TokenVersion++

	err = controller.DB.Save(&user).Error
	helper.PanicError(err)

	response := web.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Image:     user.ImageUrl,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   response,
	})
}

func (controller *UserControllerImpl) UpdateImgCtrl(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return exception.ErrorHandler(400, "BAD REQUEST", err)(c)
	}
	helper.PanicError(err)

	file, err := fileHeader.Open()
	helper.PanicError(err)

	param := uploader.UploadParams{
		PublicID:       fileHeader.Filename,
		Folder:         "meals-app",
		AllowedFormats: []string{"jpg", "png", "jpeg"},
	}

	uploadResult, err := controller.Cld.Upload.Upload(c.Context(), file, param)
	helper.PanicError(err)

	user := c.Locals("currentUser").(entity.User)
	user.ImageUrl = uploadResult.SecureURL

	err = controller.DB.Save(&user).Error
	helper.PanicError(err)

	response := web.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Image:     user.ImageUrl,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   response,
	})
}

func (controller *UserControllerImpl) GetAllFavoriteCtrl(c *fiber.Ctx) error{
	user := c.Locals("currentUser").(entity.User)

	var meals []entity.MealRecipe
	err := controller.DB.Model(&user).Association("FavoriteMeals").Find(&meals)
	helper.PanicError(err)

	var responses []web.MealResponse
	for _, meal := range meals {
		var ingredients []string
		var steps []string

		for _, ingredient := range meal.Ingredients {
			ingredients = append(ingredients, ingredient.Ingredient)
		}

		for _, step := range meal.Steps {
			steps = append(steps, step.Step)
		}

		mealResponse := web.MealResponse{
			ID:            meal.ID,
			UserId:        meal.UserId,
			Name:          meal.Name,
			Category:      meal.Category,
			ImageUrl:      meal.ImageUrl,
			Duration:      meal.Duration,
			Complexity:    meal.Complexity,
			Affordability: meal.Affordability,
			IsGlutenFree:  meal.IsGlutenFree,
			IsLactoseFree: meal.IsLactoseFree,
			IsVegan:       meal.IsVegan,
			Ingredients:   ingredients,
			Steps:         steps,
			CreatedAt:     meal.CreatedAt,
			UpdatedAt:     meal.UpdatedAt,
		}

		responses = append(responses, mealResponse)
	}

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   responses,
	})
}
