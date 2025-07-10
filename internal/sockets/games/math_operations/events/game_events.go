package math_operations_events

type MathOperationsEvent string

var (
	MathOperationsEventQuestion MathOperationsEvent = "Question"
)

func (e MathOperationsEvent) String() string {
	switch e {
	case MathOperationsEventQuestion:
		return "Question"
	default:
		return "Unknown"
	}
}
