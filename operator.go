// SPDX-License-Identifier: BSD-3-Clause

package kiwi

type Operator int

const (
	LE Operator = iota
	GE
	EQ
)

func (o Operator) String() string {
	return [...]string{"<=", ">=", "=="}[o]
}
