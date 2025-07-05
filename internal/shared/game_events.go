package shared

type CommonGameEvent string

var (
	CommonGameEventCorrectAnswer CommonGameEvent = "CorrectAnswer"
	CommonGameEventWrongAnswer   CommonGameEvent = "WrongAnswer"
	CommonGameEventFinishedGame  CommonGameEvent = "FinishedGame"
	CommonGameEventScoreboard    CommonGameEvent = "Scoreboard"
)

func (e CommonGameEvent) String() string {
	switch e {
	case CommonGameEventCorrectAnswer:
		return "CorrectAnswer"
	case CommonGameEventWrongAnswer:
		return "WrongAnswer"
	case CommonGameEventFinishedGame:
		return "FinishedGame"
	case CommonGameEventScoreboard:
		return "Scoreboard"
	default:
		return "Unknown"
	}
}
