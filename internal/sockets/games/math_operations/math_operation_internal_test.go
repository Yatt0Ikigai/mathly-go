package math_operations

import (
	"encoding/json"
	"mathly/internal/models"
	"mathly/internal/shared"
	math_operations_events "mathly/internal/sockets/games/math_operations/events"
	gameUtils "mathly/internal/sockets/games/utils"
	"mathly/internal/utils"

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

	playerOneId, _ := uuid.Parse("10zcxzxc-3144-43fe-ab86-aab9f5c8de10")
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
		randomMock *utils.MockRandom
	)

	BeforeEach(func() {
		randomCtrl = gomock.NewController(GinkgoT())
		randomMock = utils.NewMockRandom(randomCtrl)

		listener := make(chan models.Message)
		config := gameUtils.GameConfig{
			MessageListener: listener,
			Random:          randomMock,
		}

		broadcast = make(chan shared.SocketReponse, 10)

		mathOperationGame = InitMockOperationsGame(config, players, broadcast)
	})

	Describe("generateAdditionQuestion", func() {
		It("should return correct question", func() {
			// given
			gomock.InOrder(
				randomMock.EXPECT().Intn(1000).Return(600),
				randomMock.EXPECT().Intn(1000).Return(600),
				randomMock.EXPECT().Intn(100).Return(10),
				randomMock.EXPECT().Intn(100).Return(60),
				randomMock.EXPECT().Intn(100).Return(80),
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
