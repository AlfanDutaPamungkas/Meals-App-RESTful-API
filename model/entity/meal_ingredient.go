package entity

import "time"

type MealIngredient struct {
	ID           int        `json:"id"`
	MealRecipeId int        `json:"meal_recipe_id"`
	Ingredient   string     `json:"ingredient"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	MealRecipe   MealRecipe `gorm:"foreignKey:MealRecipeId;references:ID;OnDelete:CASCADE"`
}
