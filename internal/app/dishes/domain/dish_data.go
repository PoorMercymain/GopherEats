package domain

type Dish struct {
	Id          uint64
	Name        string
	Description string
	Ingredients []*DishIngredient
}
