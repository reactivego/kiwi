package kiwi

type kind int

const (
	INVALID kind = iota
	EXTERNAL
	SLACK
	ERROR
	DUMMY
)

func (k kind) is(other kind) bool {
	return k == other
}

func (k kind) String() string {
	return [...]string{"i", "v", "s", "e", "d"}[k]
}
