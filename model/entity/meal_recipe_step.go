package entity

import "time"

type MealRecipeStep struct {
	ID           int        `json:"id"`
	MealRecipeId int        `json:"meal_recipe_id"`
	Step         string     `json:"step"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	MealRecipe   MealRecipe `gorm:"foreignKey:MealRecipeId;references:ID;OnDelete:CASCADE"`
}
