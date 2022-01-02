package kiwi

const (
	INVALID kind = iota
	EXTERNAL
	SLACK
	ERROR
	DUMMY
)

type kind int

func (k kind) is(other kind) bool {
	return k == other
}

func (k kind) String() string {
	switch k {
	case INVALID:
		return "i"
	case EXTERNAL:
		return "v"
	case SLACK:
		return "s"
	case ERROR:
		return "e"
	case DUMMY:
		return "d"
	default:
		return "?"
	}
}
