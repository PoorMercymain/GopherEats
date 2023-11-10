package domain

type Ingredient struct {
	Id   uint64
	Name string
	Unit string
}

type DishIngredient struct {
	Id   uint64
	Name string
	Unit string
	Qty  uint64
}
