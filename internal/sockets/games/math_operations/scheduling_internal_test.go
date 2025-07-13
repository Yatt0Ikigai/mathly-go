package math_operations

import (
	"mathly/internal/mocks"
	"mathly/internal/models"
	"mathly/internal/shared"
	common_games "mathly/internal/sockets/games/common"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("User", Ordered, func() {
	playerOneReceiver := make(chan shared.SocketResponse, 10)
	playerTwoReceiver := make(chan shared.SocketResponse, 10)

	jobId, _ := uuid.Parse("89f5b4df-d552-4f6f-9046-64b79f47cd7d")
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
		mathOperationGame mathOperations

		schedulerCtrl *gomock.Controller
		schedulerMock *mocks.MockScheduler
	)

	BeforeEach(func() {
		schedulerCtrl = gomock.NewController(GinkgoT())
		schedulerMock = mocks.NewMockScheduler(schedulerCtrl)

		config := common_games.GameConfig{
			Settings:  common_games.GameSettings{},
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
			events: scheduledEvents{
				PlayerTurns: make(map[uuid.UUID]*uuid.UUID),
			},
		}
	})

	It("addEventGameEvent", func() {
		jobCtrl := gomock.NewController(GinkgoT())
		jobMock := mocks.NewMockJob(jobCtrl)

		schedulerMock.EXPECT().NewJob(
			gocron.DurationJob(3*time.Minute),
			gomock.Any(),
			gomock.AssignableToTypeOf(gocron.WithLimitedRuns(1)),
		).Return(jobMock, nil)

		jobMock.EXPECT().ID().Return(jobId)

		// when
		mathOperationGame.addEventGameEvent()

		// then
		Expect(mathOperationGame.events.EndOfGame).To(Equal(jobId))
	})

	It("removePlayerTurnJob", func() {
		mathOperationGame.events.PlayerTurns[playerOneId] = &jobId

		schedulerMock.EXPECT().RemoveJob(jobId)

		// when
		mathOperationGame.removePlayerTurnJob(playerOneId)

		// then
		Expect(mathOperationGame.events.PlayerTurns[playerOneId]).To(BeNil())
	})

	It("addPlayerTurnJob", func() {
		jobCtrl := gomock.NewController(GinkgoT())
		jobMock := mocks.NewMockJob(jobCtrl)

		schedulerMock.EXPECT().NewJob(
			gocron.DurationJob(30*time.Second),
			gomock.Any(),
			gomock.AssignableToTypeOf(gocron.WithLimitedRuns(1)),
		).Return(jobMock, nil)

		jobMock.EXPECT().ID().Return(jobId)

		// when
		mathOperationGame.addPlayerTurnJob(playerOneId)

		// then
		Expect(mathOperationGame.events.PlayerTurns[playerOneId]).To(Equal(&jobId))
	})
})
