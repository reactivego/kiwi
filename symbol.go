package kiwi

type Symbol int

const (
	INVALID Symbol = iota
	EXTERNAL
	SLACK
	ERROR
	DUMMY
)

func NewSymbol(symbol Symbol) *Symbol {
	return &symbol
}

func NewInvalidSymbol() *Symbol {
	return NewSymbol(INVALID)
}

func (s Symbol) IsInvalid() bool {
	return s == INVALID
}

func (s Symbol) IsExternal() bool {
	return s == EXTERNAL
}

func (s Symbol) IsSlack() bool {
	return s == SLACK
}

func (s Symbol) IsError() bool {
	return s == ERROR
}

func (s Symbol) IsDummy() bool {
	return s == DUMMY
}
