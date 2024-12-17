package controller

import (
	"errors"
	"meals-app/exception"
	"meals-app/helper"
	"meals-app/model/entity"
	"meals-app/model/web"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MealControllerImpl struct {
	DB       *gorm.DB
	Validate *validator.Validate
	Cld      *cloudinary.Cloudinary
}

func NewMealControllerImpl(DB *gorm.DB, validate *validator.Validate, cld *cloudinary.Cloudinary) MealController {
	return &MealControllerImpl{
		DB:       DB,
		Validate: validate,
		Cld:      cld,
	}
}

func (controller *MealControllerImpl) CreateMealCtrl(c *fiber.Ctx) error {
	request := new(web.CreateMealReq)
	err := c.BodyParser(request)
	if err != nil {
		return exception.ErrorHandler(400, "BAD REQUEST", err)(c)
	}

	err = controller.Validate.Struct(request)
	if err != nil {
		return exception.ErrorHandler(422, "VALIDATION ERROR", err)(c)
	}

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

	isGlutenFree, err := strconv.ParseBool(request.IsGlutenFree)
	helper.PanicError(err)
	isLactoseFree, err := strconv.ParseBool(request.IsLactoseFree)
	helper.PanicError(err)
	isVegan, err := strconv.ParseBool(request.IsVegan)
	helper.PanicError(err)

	user := c.Locals("currentUser").(entity.User)

	var mealRecipe entity.MealRecipe

	err = controller.DB.Transaction(func(tx *gorm.DB) error {
		mealRecipe = entity.MealRecipe{
			UserId:        user.ID,
			Name:          request.Name,
			Category:      request.Category,
			ImageUrl:      uploadResult.SecureURL,
			Duration:      request.Duration,
			Complexity:    request.Complexity,
			Affordability: request.Affordability,
			IsGlutenFree:  isGlutenFree,
			IsLactoseFree: isLactoseFree,
			IsVegan:       isVegan,
		}

		err := tx.Create(&mealRecipe).Error
		if err != nil {
			return err
		}

		for _, ingredient := range request.Ingredients {
			mealIngredient := entity.MealIngredient{
				MealRecipeId: mealRecipe.ID,
				Ingredient:   ingredient,
			}
			err = tx.Create(&mealIngredient).Error
			if err != nil {
				return err
			}
		}

		for _, step := range request.Steps {
			mealStep := entity.MealRecipeStep{
				MealRecipeId: mealRecipe.ID,
				Step:         step,
			}
			err = tx.Create(&mealStep).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	helper.PanicError(err)

	response := web.MealResponse{
		ID:            mealRecipe.ID,
		UserId:        mealRecipe.UserId,
		Name:          mealRecipe.Name,
		Category:      mealRecipe.Category,
		ImageUrl:      mealRecipe.ImageUrl,
		Duration:      mealRecipe.Duration,
		Complexity:    mealRecipe.Complexity,
		Affordability: mealRecipe.Affordability,
		IsGlutenFree:  mealRecipe.IsGlutenFree,
		IsLactoseFree: mealRecipe.IsLactoseFree,
		IsVegan:       mealRecipe.IsVegan,
		Ingredients:   request.Ingredients,
		Steps:         request.Steps,
		CreatedAt:     mealRecipe.CreatedAt,
		UpdatedAt:     mealRecipe.UpdatedAt,
	}

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   response,
	})
}

func (controller *MealControllerImpl) GetAllMealCtrl(c *fiber.Ctx) error {
	var mealRecipes []entity.MealRecipe
	var err error

	name := c.Query("name")
	if name != "" {
		err = controller.DB.Raw(`
			SELECT 
				id,
				user_id,
				name,
				category,
				image_url,
				duration,
				complexity,
				affordability,
				is_gluten_free,
				is_lactose_free,
				is_vegan,
				created_at,
				updated_at
			FROM 
				meal_recipes
			WHERE
				MATCH(name) AGAINST (? IN NATURAL LANGUAGE MODE)
		`, name).Preload("Ingredients").Preload("Steps").Find(&mealRecipes).Error
		helper.PanicError(err)
	} else {
		err = controller.DB.Preload("Ingredients").Preload("Steps").Find(&mealRecipes).Error
		helper.PanicError(err)
	}

	var responses []web.MealResponse
	for _, meal := range mealRecipes {
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

func (controller *MealControllerImpl) GetMealByIDCtrl(c *fiber.Ctx) error {
	mealID, err := c.ParamsInt("id")
	helper.PanicError(err)

	meal := entity.MealRecipe{}
	err = controller.DB.Preload("Ingredients").Preload("Steps").Take(&meal, "id = ?", mealID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.ErrorHandler(404, "NOT FOUND", errors.New("meal recipe not found"))(c)
		}

		return exception.ErrorHandler(500, "INTERNAL SERVER ERROR", err)(c)
	}

	var ingredients []string
	for _, ingredient := range meal.Ingredients {
		ingredients = append(ingredients, ingredient.Ingredient)
	}

	var steps []string
	for _, step := range meal.Steps {
		steps = append(steps, step.Step)
	}

	response := web.MealResponse{
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

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   response,
	})
}

func (controller *MealControllerImpl) UpdateMealCtrl(c *fiber.Ctx) error {
	mealID, err := c.ParamsInt("id")
	helper.PanicError(err)

	user := c.Locals("currentUser").(entity.User)

	meal := entity.MealRecipe{}
	err = controller.DB.Preload("Ingredients").Preload("Steps").Take(&meal, "id = ?", mealID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.ErrorHandler(404, "NOT FOUND", errors.New("meal recipe not found"))(c)
		}

		return exception.ErrorHandler(500, "INTERNAL SERVER ERROR", err)(c)
	}

	if meal.UserId != user.ID {
		return exception.ErrorHandler(403, "FORBIDDEN", errors.New("forbidden, you are not allowed"))(c)
	}

	request := new(web.UpdateMealReq)
	err = c.BodyParser(request)
	if err != nil {
		return exception.ErrorHandler(400, "BAD REQUEST", err)(c)
	}

	err = controller.Validate.Struct(request)
	if err != nil {
		return exception.ErrorHandler(422, "VALIDATION ERROR", err)(c)
	}

	if request.Name != "" {
		meal.Name = request.Name
	}

	if request.Category != "" {
		meal.Category = request.Category
	}

	if request.Duration != "" {
		meal.Duration = request.Duration
	}

	if request.Complexity != "" {
		meal.Complexity = request.Complexity
	}

	if request.Affordability != "" {
		meal.Affordability = request.Affordability
	}

	if request.IsGlutenFree != "" {
		isGlutenFree, err := strconv.ParseBool(request.IsGlutenFree)
		helper.PanicError(err)
		meal.IsGlutenFree = isGlutenFree
	}

	if request.IsLactoseFree != "" {
		isLactoseFree, err := strconv.ParseBool(request.IsLactoseFree)
		helper.PanicError(err)
		meal.IsLactoseFree = isLactoseFree
	}

	if request.IsVegan != "" {
		isVegan, err := strconv.ParseBool(request.IsVegan)
		helper.PanicError(err)
		meal.IsVegan = isVegan
	}

	err = controller.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Save(&meal).Error
		if err != nil {
			return err
		}

		if len(request.Ingredients) > 0 {
			tx.Where("meal_recipe_id = ?", mealID).Delete(&entity.MealIngredient{})
			var newIngredients []entity.MealIngredient
			for _, ingredient := range request.Ingredients {
				newIngredients = append(newIngredients, entity.MealIngredient{
					MealRecipeId: mealID,
					Ingredient:   ingredient,
				})
			}
			err = tx.Create(&newIngredients).Error
			if err != nil {
				return err
			}
		}

		if len(request.Steps) > 0 {
			tx.Where("meal_recipe_id = ?", mealID).Delete(&entity.MealRecipeStep{})
			var newSteps []entity.MealRecipeStep
			for _, step := range request.Steps {
				newSteps = append(newSteps, entity.MealRecipeStep{
					MealRecipeId: mealID,
					Step:         step,
				})
			}
			err = tx.Create(&newSteps).Error
			if err != nil {
				return err
			}
		}

		err = tx.Preload("Ingredients").Preload("Steps").Take(&meal, "id = ?", mealID).Error
		if err != nil {
			return err
		}

		return nil
	})

	helper.PanicError(err)

	var ingredients []string
	for _, ingredient := range meal.Ingredients {
		ingredients = append(ingredients, ingredient.Ingredient)
	}

	var steps []string
	for _, step := range meal.Steps {
		steps = append(steps, step.Step)
	}

	response := web.MealResponse{
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

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   response,
	})
}

func (controller *MealControllerImpl) UpdateMealImageCtrl(c *fiber.Ctx) error {
	mealID, err := c.ParamsInt("id")
	helper.PanicError(err)

	user := c.Locals("currentUser").(entity.User)

	meal := entity.MealRecipe{}
	err = controller.DB.Preload("Ingredients").Preload("Steps").Take(&meal, "id = ?", mealID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.ErrorHandler(404, "NOT FOUND", errors.New("meal recipe not found"))(c)
		}

		return exception.ErrorHandler(500, "INTERNAL SERVER ERROR", err)(c)
	}

	if meal.UserId != user.ID {
		return exception.ErrorHandler(403, "FORBIDDEN", errors.New("forbidden, you are not allowed"))(c)
	}

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

	meal.ImageUrl = uploadResult.SecureURL

	err = controller.DB.Save(&meal).Error
	helper.PanicError(err)

	var ingredients []string
	for _, ingredient := range meal.Ingredients {
		ingredients = append(ingredients, ingredient.Ingredient)
	}

	var steps []string
	for _, step := range meal.Steps {
		steps = append(steps, step.Step)
	}

	response := web.MealResponse{
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

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   response,
	})
}

func (controller *MealControllerImpl) DeleteMealCtrl(c *fiber.Ctx) error {
	mealID, err := c.ParamsInt("id")
	helper.PanicError(err)

	user := c.Locals("currentUser").(entity.User)

	meal := entity.MealRecipe{}
	err = controller.DB.Preload("Ingredients").Preload("Steps").Take(&meal, "id = ?", mealID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.ErrorHandler(404, "NOT FOUND", errors.New("meal recipe not found"))(c)
		}

		return exception.ErrorHandler(500, "INTERNAL SERVER ERROR", err)(c)
	}

	if meal.UserId == user.ID || user.Role == "admin" {
		err = controller.DB.Delete(&meal).Error
		helper.PanicError(err)

		return c.Status(200).JSON(fiber.Map{
			"code":   200,
			"status": "success",
			"data":   "meal recipe deleted successfully",
		})
	}

	return exception.ErrorHandler(403, "FORBIDDEN", errors.New("forbidden, you are not allowed"))(c)
}

func (controller *MealControllerImpl) AddToFavoriteCtrl(c *fiber.Ctx) error {
	mealID, err := c.ParamsInt("id")
	helper.PanicError(err)

	meal := entity.MealRecipe{}
	err = controller.DB.Preload("Ingredients").Preload("Steps").Take(&meal, "id = ?", mealID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.ErrorHandler(404, "NOT FOUND", errors.New("meal recipe not found"))(c)
		}

		return exception.ErrorHandler(500, "INTERNAL SERVER ERROR", err)(c)
	}

	user := c.Locals("currentUser").(entity.User)

	err = controller.DB.Model(&user).Association("FavoriteMeals").Append(&meal)
	helper.PanicError(err)

	var ingredients []string
	for _, ingredient := range meal.Ingredients {
		ingredients = append(ingredients, ingredient.Ingredient)
	}

	var steps []string
	for _, step := range meal.Steps {
		steps = append(steps, step.Step)
	}

	response := web.MealResponse{
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

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   response,
	})
}

func (controller *MealControllerImpl) DeleteFromFavoriteCtrl(c *fiber.Ctx) error {
	mealID, err := c.ParamsInt("id")
	helper.PanicError(err)

	user := c.Locals("currentUser").(entity.User)

	var meal entity.MealRecipe
	err = controller.DB.
		Joins("JOIN favorite_user_meal ON meal_recipes.id = favorite_user_meal.meal_recipe_id").Where("meal_recipes.id = ? AND favorite_user_meal.user_id = ?", mealID, user.ID).First(&meal).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.ErrorHandler(404, "NOT FOUND", errors.New("meal recipe is not in favorites"))(c)
		}

		return exception.ErrorHandler(500, "INTERNAL SERVER ERROR", err)(c)
	}

	err = controller.DB.Model(&user).Association("FavoriteMeals").Delete(&meal)
	helper.PanicError(err)

	return c.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "success",
		"data":   "meal recipe removed from favorites successfully",
	})
}
