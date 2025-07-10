//go:generate mockgen -source=random.go -package service -destination=random_mock.go
package service

import (
	"errors"
	"math/rand"
)

type Random interface {
	GenerateRandomNumber(nums ...int) (int, error)
}

type random struct{}

func newRandom() Random {
	return &random{}
}

func (r random) GenerateRandomNumber(nums ...int) (int, error) {
	switch len(nums) {
	case 1:
		if nums[0] <= 0 {
			return 0, errors.New("number must be greater than 0")
		}
		return rand.Intn(nums[0]), nil
	case 2:
		min, max := nums[0], nums[1]
		if min >= max {
			return 0, errors.New("min must be less than max")
		}
		return rand.Intn(max-min) + min, nil
	default:
		return 0, errors.New("must provide 1 or 2 numbers")
	}
}
