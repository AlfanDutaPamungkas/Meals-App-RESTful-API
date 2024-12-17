package entity

import "time"

type User struct {
	ID            int          `json:"id"`
	Username      string       `json:"username"`
	Email         string       `json:"email"`
	Role          string       `json:"role" gorm:"default:user"`
	Password      string       `json:"password"`
	ImageUrl      string       `json:"image_url"`
	TokenVersion  int          `json:"token_version" gorm:"default:1"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"update_at"`
	MealRecipes   []MealRecipe `gorm:"foreignKey:UserId;references:ID"`
	FavoriteMeals []MealRecipe `gorm:"many2many:favorite_user_meal;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:meal_recipe_id"`
}
