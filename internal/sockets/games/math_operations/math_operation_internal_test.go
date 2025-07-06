package math_operations

import (
	"encoding/json"
	"mathly/internal/models"
	"mathly/internal/service"
	"mathly/internal/shared"
	"mathly/internal/sockets/games/common"
	math_operations_events "mathly/internal/sockets/games/math_operations/events"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func mockSender(receiver chan shared.SocketReponse) func(shared.SocketReponse) {
	return func(b shared.SocketReponse) {
		receiver <- b
	}
}

var _ = Describe("User", Ordered, func() {
	playerOneReceiver := make(chan shared.SocketReponse, 10)
	playerTwoReceiver := make(chan shared.SocketReponse, 10)

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
		broadcast         chan shared.SocketReponse
		mathOperationGame mathOperations

		randomCtrl *gomock.Controller

		randomMock *service.MockRandom
	)

	BeforeEach(func() {
		randomCtrl = gomock.NewController(GinkgoT())
		randomMock = service.NewMockRandom(randomCtrl)

		broadcast = make(chan shared.SocketReponse, 10)

		config := common_games.GameConfig{
			Services: common_games.GameServices{
				Random: randomMock,
			},
			Settings:        common_games.GameSettings{},
			MessageListener: nil,
			Broadcast:       broadcast,
			Players:         players,
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
		})

		It("should increase user score and send new question", func() {
			// given
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
		})

		It("should decrease user score and send new question", func() {
			// given
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
		})
	})
})
