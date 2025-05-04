package math_operations

import (
	"encoding/json"
	"mathly/internal/log"
)

type MathQuestion struct {
	Question      string
	Answers       []int
	correctAnswer int
}

func (m MathQuestion) String() string {
	r, err := json.Marshal(m.Answers)
	if err != nil {
		log.Log.Errorf(`Failed to marshal Math Question: %v`, err)
		return ""
	}
	return string(r)
}
