//go:generate mockgen -source=rand.go -package utils -destination=rand_mock.go
package utils

import (
	"math/rand"
)

type Random interface {
	Intn(n int) int
}

type random struct{}

func (r random) Intn(n int) int {
	return rand.Intn(n)
}

func NewRandom() Random {
	return random{}
}
