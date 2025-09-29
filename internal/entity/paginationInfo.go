package entity

type PaginationInfo struct {
	Comments       []Comment
	HasPrev        bool
	HasNext        bool
	NextOffset     string
	PrevOffset     string
	SortFromOldest bool
}
