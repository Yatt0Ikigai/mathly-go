package math_operations

import (
	"encoding/json"
	"mathly/internal/log"
)

type MathQuestion struct {
	Question      string
	Answers       []string
	correctAnswer string
}

func (m MathQuestion) String() string {
	r, err := json.Marshal(m)
	if err != nil {
		log.Log.Errorf(`Failed to marshal Math Question: %v`, err)
		return ""
	}
	return string(r)
}
