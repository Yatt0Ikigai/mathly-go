package math_operations

import (
	"encoding/json"
	"mathly/internal/mocks"
	"mathly/internal/models"
	"mathly/internal/service"
	"mathly/internal/shared"
	common_games "mathly/internal/sockets/games/common"
	math_operations_events "mathly/internal/sockets/games/math_operations/events"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func mockSender(receiver chan shared.SocketResponse) func(shared.SocketResponse) {
	return func(b shared.SocketResponse) {
		receiver <- b
	}
}

var _ = Describe("Questions", Ordered, func() {
	playerOneReceiver := make(chan shared.SocketResponse, 10)
	playerTwoReceiver := make(chan shared.SocketResponse, 10)

	playerOneId, _ := uuid.Parse("95f3cec5-ca92-45a4-a3b7-b5001eaad1b4")
	playerTwoId, _ := uuid.Parse("80asdasd-3144-43fe-ab86-aab9f5c8de80")
	playerOne := models.Player{
		Nickname:     "Test",
		ConnectionID: playerOneId,
		SendMessage:  mockSender(playerOneReceiver),
	}
	playerTwo := models.Player{
		Nickname:     "TestTwo",
		ConnectionID: playerTwoId,
		SendMessage:  mockSender(playerTwoReceiver),
	}

	players := make(map[uuid.UUID]models.Player)
	players[playerOneId] = playerOne
	players[playerTwoId] = playerTwo

	var (
		broadcast         chan shared.SocketResponse
		mathOperationGame mathOperations

		randomCtrl    *gomock.Controller
		schedulerCtrl *gomock.Controller

		randomMock    *service.MockRandom
		schedulerMock *mocks.MockScheduler
	)

	BeforeEach(func() {
		randomCtrl = gomock.NewController(GinkgoT())
		schedulerCtrl = gomock.NewController(GinkgoT())
		randomMock = service.NewMockRandom(randomCtrl)
		schedulerMock = mocks.NewMockScheduler(schedulerCtrl)

		broadcast = make(chan shared.SocketResponse, 10)

		config := common_games.GameConfig{
			Services: common_games.GameServices{
				Random: randomMock,
			},
			Settings:  common_games.GameSettings{},
			Broadcast: broadcast,
			Players:   players,
			Scheduler: schedulerMock,
		}

		scoreboard := make(map[uuid.UUID]int)
		scoreboard[playerOneId] = 0
		scoreboard[playerTwoId] = 0

		playerScoreboard := make(map[uuid.UUID]int)
		playerScoreboard[playerOneId] = 0
		playerScoreboard[playerTwoId] = 0

		mathOperationGame = mathOperations{
			config:         config,
			scoreBoard:     scoreboard,
			playerQuestion: playerScoreboard,
		}
	})

	Describe("generateAdditionQuestion", func() {
		It("should return correct question", func() {
			// given
			gomock.InOrder(
				randomMock.EXPECT().GenerateRandomNumber(1000).Return(600, nil),
				randomMock.EXPECT().GenerateRandomNumber(1000).Return(600, nil),
				randomMock.EXPECT().GenerateRandomNumber(100).Return(10, nil),
				randomMock.EXPECT().GenerateRandomNumber(100).Return(60, nil),
				randomMock.EXPECT().GenerateRandomNumber(100).Return(80, nil),
			)

			// when
			question := mathOperationGame.generateAdditionQuestion()

			// then
			Expect(question.Question).To(Equal(`What's the sum of 100 + 100 ?`))
			Expect(len(question.Answers)).To(Equal(4))
			Expect(question.correctAnswer).To(Equal("200"))
			expectedAnswers := []string{"160", "200", "210", "230"}
			for _, item := range expectedAnswers {
				Expect(question.Answers).To(ContainElement(item))
			}
		})
	})

	Describe("handleAnswer", func() {
		var (
			correctAnswer   = "255"
			incorrectAnswer = "2137"
			jobId, _        = uuid.Parse("89f5b4df-d552-4f6f-9046-64b79f47cd7d")
			secondJobId, _  = uuid.Parse("e9b4451e-b002-4636-867a-19032a6fac7c")
		)

		BeforeEach(func() {
			mathOperationGame.questions = []MathQuestion{
				{
					Question:      "Q1",
					Answers:       []string{"A", "B"},
					correctAnswer: correctAnswer,
				},
				{
					Question:      "Q2",
					Answers:       []string{"D", "C"},
					correctAnswer: correctAnswer},
			}
			mathOperationGame.events.PlayerTurns = make(map[uuid.UUID]*uuid.UUID)
			mathOperationGame.events.PlayerTurns[playerOneId] = &jobId
		})

		It("should increase user score and send new question", func() {
			// given
			jobCtrl := gomock.NewController(GinkgoT())
			jobMock := mocks.NewMockJob(jobCtrl)
			jobMock.EXPECT().ID().Return(secondJobId)

			schedulerMock.EXPECT().RemoveJob(jobId)
			schedulerMock.EXPECT().NewJob(
				gocron.DurationJob(30*time.Second),
				gomock.Any(),
				gomock.AssignableToTypeOf(gocron.WithLimitedRuns(1)),
			).Return(jobMock, nil)

			answer := UserMessageData{
				Answer: correctAnswer,
			}
			data, _ := json.Marshal(answer)

			// when
			mathOperationGame.handleAnswerMessage(models.Message{
				SenderID: playerOneId,
				MessageDetails: models.MessageDetails{
					Type:   models.MessageTypeGame,
					Action: models.ActionTypeGuessAnswer,
					Data:   string(data),
				},
			})

			// then
			msg := <-playerOneReceiver
			Expect(msg).To(Equal(
				shared.CreateSocketResponse(
					shared.EventGame,
					math_operations_events.MathOperationsEventQuestion,
					"{\"Question\":\"Q2\",\"Answers\":[\"D\",\"C\"]}",
				),
			))
			msg = <-broadcast
			Expect(msg).To(Equal(
				shared.CreateSocketResponse(
					shared.EventGame,
					shared.CommonGameEventScoreboard,
					"{\"Test\":1,\"TestTwo\":0}",
				),
			))

			Expect(mathOperationGame.scoreBoard[playerOne.ConnectionID]).To(Equal(1))
			Expect(mathOperationGame.playerQuestion[playerOne.ConnectionID]).To(Equal(1))
			Expect(mathOperationGame.events.PlayerTurns[playerOneId]).To(Equal(&secondJobId))
		})

		It("should decrease user score and send new question", func() {
			// given
			jobCtrl := gomock.NewController(GinkgoT())
			jobMock := mocks.NewMockJob(jobCtrl)
			jobMock.EXPECT().ID().Return(secondJobId)

			schedulerMock.EXPECT().RemoveJob(jobId)
			schedulerMock.EXPECT().NewJob(
				gocron.DurationJob(30*time.Second),
				gomock.Any(),
				gomock.AssignableToTypeOf(gocron.WithLimitedRuns(1)),
			).Return(jobMock, nil)

			answer := UserMessageData{
				Answer: incorrectAnswer,
			}
			data, _ := json.Marshal(answer)

			// when
			mathOperationGame.handleAnswerMessage(models.Message{
				SenderID: playerOneId,
				MessageDetails: models.MessageDetails{
					Type:   models.MessageTypeGame,
					Action: models.ActionTypeGuessAnswer,
					Data:   string(data),
				},
			})

			// then
			msg := <-playerOneReceiver
			Expect(msg).To(Equal(shared.CreateSocketResponse(shared.EventGame, math_operations_events.MathOperationsEventQuestion, "{\"Question\":\"Q2\",\"Answers\":[\"D\",\"C\"]}")))
			msg = <-broadcast
			Expect(msg).To(Equal(shared.CreateSocketResponse(shared.EventGame, shared.CommonGameEventScoreboard, "{\"Test\":-1,\"TestTwo\":0}")))

			Expect(mathOperationGame.scoreBoard[playerOne.ConnectionID]).To(Equal(-1))
			Expect(mathOperationGame.playerQuestion[playerOne.ConnectionID]).To(Equal(1))
			Expect(mathOperationGame.events.PlayerTurns[playerOneId]).To(Equal(&secondJobId))
		})
	})
})
