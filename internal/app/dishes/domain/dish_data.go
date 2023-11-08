package domain

type Dish struct {
	id          uint64
	name        string
	description string
	ingredients []*Ingredient
}
