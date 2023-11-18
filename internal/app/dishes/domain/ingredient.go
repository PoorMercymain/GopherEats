package domain

type Ingredient struct {
	ID   uint64
	Name string
	Unit string
}

type DishIngredient struct {
	ID   uint64
	Name string
	Unit string
	Qty  uint64
}
