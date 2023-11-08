package domain

type Bundle struct {
	id          uint64
	name        string
	description string
	price       uint64
	dishes      []*Dish
}
