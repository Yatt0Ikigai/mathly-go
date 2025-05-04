//go:generate mockgen -source=rand.go -package utils -destination=rand_mock.go
package utils

type Random interface {
    Intn(n int) int
}