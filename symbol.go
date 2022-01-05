// SPDX-License-Identifier: BSD-3-Clause

package kiwi

import "fmt"

type symbol struct {
	kind
	id int
}

var _sid = 0

func newSymbol(k kind) *symbol {
	_sid++
	return &symbol{k, _sid}
}

func (s symbol) String() string {
	return fmt.Sprintf("%v%d", s.kind, s.id)
}
