package math_operations

// TODO: rename it later
type PlayerType string

var (
	PlayerTypeAnswer PlayerType = "Answer"
)

type UserMessage struct {
	Type   PlayerType `json:"type"`
	Answer int        `json:"answer"`
}
