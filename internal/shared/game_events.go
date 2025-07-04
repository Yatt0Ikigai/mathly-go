package shared

type GameEvent string

var (
	GameEventCorrectAnswer GameEvent = "CorrectAnswer"
	GameEventWrongAnswer   GameEvent = "WrongAnswer"
	GameEventFinishedGame  GameEvent = "FinishedGame"
)

func (g GameEvent) String() string {
	switch g {
	case GameEventCorrectAnswer:
		return "CorrectAnswer"
	case GameEventWrongAnswer:
		return "WrongAnswer"
	case GameEventFinishedGame:
		return "FinishedGame"
	default:
		return "Unknown"
	}
}
