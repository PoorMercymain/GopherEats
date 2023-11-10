package domain

type Bundle struct {
	Id    uint64
	Name  string
	Price uint64
}

type WeeklyBundle struct {
	Id         uint64
	Name       string
	Price      uint64
	WeekNumber uint64
	Dishes     []*Dish
}
