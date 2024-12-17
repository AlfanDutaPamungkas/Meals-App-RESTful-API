package web

import "time"

type MealResponse struct {
	ID            int       `json:"id"`
	UserId        int       `json:"user_id"`
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	ImageUrl      string    `json:"image_url"`
	Duration      string    `json:"duration"`
	Complexity    string    `json:"complexity"`
	Affordability string    `json:"affordability"`
	IsGlutenFree  bool      `json:"is_gluten_free"`
	IsLactoseFree bool      `json:"is_lactose_free"`
	IsVegan       bool      `json:"is_vegan"`
	Ingredients   []string  `json:"ingredients"`
	Steps         []string  `json:"steps"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
