package symbol

const (
	INVALID Type = iota
	EXTERNAL
	SLACK
	ERROR
	DUMMY
)

type Type int

func (s Type) Type() Type {
	return s
}

func (s Type) IsInvalid() bool {
	return s == INVALID
}

func (s Type) IsExternal() bool {
	return s == EXTERNAL
}

func (s Type) IsSlack() bool {
	return s == SLACK
}

func (s Type) IsError() bool {
	return s == ERROR
}

func (s Type) IsDummy() bool {
	return s == DUMMY
}
