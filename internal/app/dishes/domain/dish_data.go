package domain

type Dish struct {
	ID          uint64
	Name        string
	Description string
	Ingredients []*DishIngredient
}
