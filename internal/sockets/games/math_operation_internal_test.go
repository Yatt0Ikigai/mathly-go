package games

import (
	"mathly/internal/models"
	"mathly/internal/utils"
	gameUtils "mathly/internal/sockets/games/utils"

	// "go.uber.org/mock/gomock"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func forwardMessageCreation(channel chan models.Message) func(models.Message) {
	return func(m models.Message) {
		channel <- m
	}
}

func broadcastMessagefuncCreation(channel chan []byte) func([]byte) {
	return func(m []byte) {
		channel <- m
	}
}

var _ = Describe("User", Ordered, func() {
	playerOneId, _ := uuid.Parse("10zcxzxc-3144-43fe-ab86-aab9f5c8de10")
	playerTwoId, _ := uuid.Parse("80asdasd-3144-43fe-ab86-aab9f5c8de80")
	playerOne := models.Player{
		Nickname:     "Test",
		ConnectionID: playerOneId,
	}
	playerTwo := models.Player{
		Nickname:     "TestTwo",
		ConnectionID: playerTwoId,
	}

	players := []models.Player{playerOne, playerTwo}
	var (
		mathOperationGame       MathOperations
		forwardMessageChannel   chan models.Message
		broadcastMessageChannel chan []byte

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

		forwardMessageChannel = make(chan models.Message)
		broadcastMessageChannel = make(chan []byte)
		fM := forwardMessageCreation(forwardMessageChannel)
		bM := broadcastMessagefuncCreation(broadcastMessageChannel)

		mathOperationGame = InitMathOperationsGame(config, players, fM, bM)
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
			Expect(question.correctAnswer).To(Equal(200))
			expectedAnswers := []int{160, 200, 210, 230}
			for _, item := range expectedAnswers {
				Expect(question.Answers).To(ContainElement(item))
			}
		})
	})

	Describe("HandleMessage", func() {
		It("", func() {

		})
	})
})
