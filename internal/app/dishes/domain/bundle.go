package domain

type Bundle struct {
	ID    uint64
	Name  string
	Price uint64
}

type WeeklyBundle struct {
	ID         uint64
	Name       string
	Price      uint64
	WeekNumber uint64
	Dishes     []*Dish
}
