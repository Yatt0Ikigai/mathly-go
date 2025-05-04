package math_operations

import (
	"encoding/json"
	"mathly/internal/models"
	gameUtils "mathly/internal/sockets/games/utils"
	"mathly/internal/utils"

	// "go.uber.org/mock/gomock"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("User", Ordered, func() {
	playerOneReceiver := make(chan []byte, 10)
	playerTwoReceiver := make(chan []byte, 10)

	playerOneId, _ := uuid.Parse("10zcxzxc-3144-43fe-ab86-aab9f5c8de10")
	playerTwoId, _ := uuid.Parse("80asdasd-3144-43fe-ab86-aab9f5c8de80")
	playerOne := models.Player{
		Nickname:     "Test",
		ConnectionID: playerOneId,
		Receiver:     playerOneReceiver,
	}
	playerTwo := models.Player{
		Nickname:     "TestTwo",
		ConnectionID: playerTwoId,
		Receiver:     playerTwoReceiver,
	}

	players := []models.Player{playerOne, playerTwo}
	var (
		broadcast         chan []byte
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

		broadcast = make(chan []byte, 10)

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
			mathOperationGame.handleAnswer(models.Message{
				SenderID: playerOneId,
				MessageDetails: models.MessageDetails{
					Type:   models.MessageTypeGame,
					Action: models.ActionTypeGuessAnswer,
					Data:   string(data),
				},
			})

			// then
			msg := <-playerOneReceiver
			Expect(string(msg)).To(Equal("{\"Question\":\"Q2\",\"Answers\":[\"D\",\"C\"]}"))
			msg = <-broadcast
			Expect(string(msg)).To(Equal("{\"Type\":\"Scoreboard\",\"Message\":\"🏆 Scoreboard:\\nTest: 1\\nTestTwo: 0\\n\"}"))

			Expect(mathOperationGame.scoreBoard[playerOne]).To(Equal(1))
			Expect(mathOperationGame.playerQuestion[playerOne]).To(Equal(1))
		})

		It("should decrease user score and send new question", func() {
			// given
			answer := UserMessageData{
				Answer: incorrectAnswer,
			}
			data, _ := json.Marshal(answer)

			// when
			mathOperationGame.handleAnswer(models.Message{
				SenderID: playerOneId,
				MessageDetails: models.MessageDetails{
					Type:   models.MessageTypeGame,
					Action: models.ActionTypeGuessAnswer,
					Data:   string(data),
				},
			})

			// then
			msg := <-playerOneReceiver
			Expect(string(msg)).To(Equal("{\"Question\":\"Q2\",\"Answers\":[\"D\",\"C\"]}"))
			msg = <-broadcast
			Expect(string(msg)).To(Equal("{\"Type\":\"Scoreboard\",\"Message\":\"🏆 Scoreboard:\\nTest: -1\\nTestTwo: 0\\n\"}"))

			Expect(mathOperationGame.scoreBoard[playerOne]).To(Equal(-1))
			Expect(mathOperationGame.playerQuestion[playerOne]).To(Equal(1))
		})
	})
})
