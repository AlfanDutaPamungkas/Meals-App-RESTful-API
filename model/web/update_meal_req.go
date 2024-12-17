package web

type UpdateMealReq struct {
	Name          string   `json:"name" validate:"max=100"`
	Category      string   `json:"category"`
	Duration      string   `json:"duration"`
	Complexity    string   `json:"complexity"`
	Affordability string   `json:"affordability"`
	IsGlutenFree  string   `json:"is_gluten_free"`
	IsLactoseFree string   `json:"is_lactose_free"`
	IsVegan       string   `json:"is_vegan"`
	Ingredients   []string `json:"ingredients"`
	Steps         []string `json:"steps"`
}
