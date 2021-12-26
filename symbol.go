package kiwi

type symbol int

const (
	INVALID symbol = iota
	EXTERNAL
	SLACK
	ERROR
	DUMMY
)

func newSymbol(s symbol) *symbol {
	return &s
}

func (s symbol) is(kind symbol) bool {
	return s == kind
}
