package entity

import "time"

type MealRecipe struct {
	ID               int              `json:"id"`
	UserId           int              `json:"user_id"`
	Name             string           `json:"name"`
	Category         string           `json:"category"`
	ImageUrl         string           `json:"image_url"`
	Duration         string           `json:"duration"`
	Complexity       string           `json:"complexity"`
	Affordability    string           `json:"affordability"`
	IsGlutenFree     bool             `json:"is_gluten_free"`
	IsLactoseFree    bool             `json:"is_lactose_free"`
	IsVegan          bool             `json:"is_vegan"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	User             User             `gorm:"foreignKey:UserId;references:ID"`
	Ingredients      []MealIngredient `gorm:"foreignKey:MealRecipeId;references:ID"`
	Steps            []MealRecipeStep `gorm:"foreignKey:MealRecipeId;references:ID"`
	FavoritedByUsers []User `gorm:"many2many:favorite_user_meal;foreignKey:id;joinForeignKey:meal_recipe_id;references:id;joinReferences:user_id"`
}
