package filter

type FilterOperation string

const (
	FilterEQ      FilterOperation = "="
	FilterNEQ     FilterOperation = "!="
	FilterLike    FilterOperation = "like"
	FilterNotLike FilterOperation = "nlike"
	FilterIn      FilterOperation = "in"
	FilterNotIn   FilterOperation = "nin"
)

type Filter struct {
	Field     string
	Operation FilterOperation
	Value     string
}
