package web

type CreateMealReq struct {
	Name          string   `form:"name" validate:"required,max=100"`
	Category      string   `form:"category" validate:"required"`
	Duration      string   `form:"duration" validate:"required"`
	Complexity    string   `form:"complexity" validate:"required"`
	Affordability string   `form:"affordability" validate:"required"`
	IsGlutenFree  string   `form:"is_gluten_free" validate:"required"`
	IsLactoseFree string   `form:"is_lactose_free" validate:"required"`
	IsVegan       string   `form:"is_vegan" validate:"required"`
	Ingredients   []string `form:"ingredients[]" validate:"required"`
	Steps         []string `form:"steps[]" validate:"required"`
}
