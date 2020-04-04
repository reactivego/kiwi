package kiwi

import "github.com/reactivego/kiwi/symbol"

type Symbol interface {
	Type() symbol.Type
	IsInvalid() bool
	IsExternal() bool
	IsSlack() bool
	IsError() bool
	IsDummy() bool
}

func NewInvalidSymbol() Symbol {
	return NewSymbol(symbol.INVALID)
}

func NewSymbol(symbolType symbol.Type) Symbol {
	return &symbolType
}
