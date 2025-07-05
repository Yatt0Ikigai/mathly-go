package sockets

import (
	"encoding/json"
	"fmt"
	"mathly/internal/models"
	"mathly/internal/service"
	"mathly/internal/shared"
	"mathly/internal/sockets/games/math_operations"
	math_operations_events "mathly/internal/sockets/games/math_operations/events"
	"time"

	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var (
	serviceCtrl *gomock.Controller
	serviceMock *service.MockService

	clientOneCtrl *gomock.Controller
	clientTwoCtrl *gomock.Controller
	clientOneMock *MockClient
	clientTwoMock *MockClient
)

var _ = Describe("Lobby", Ordered, func() {
	type player struct {
		id       string
		nickname string
	}

	playerOne := player{
		id:       "95f3cec5-ca92-45a4-a3b7-b5001eaad1b4",
		nickname: "clientOne",
	}
	playerTwo := player{
		id:       "b1eb7314-459d-4e4f-90ed-b79e1f8add5a",
		nickname: "clientTwo",
	}

	clientOneId, _ := uuid.Parse(playerOne.id)
	clientTwoId, _ := uuid.Parse(playerTwo.id)
	var L Lobby

	BeforeAll(func() {
		clientOneCtrl = gomock.NewController(GinkgoT())
		clientTwoCtrl = gomock.NewController(GinkgoT())
		serviceCtrl = gomock.NewController(GinkgoT())

		clientOneMock = NewMockClient(clientOneCtrl)
		clientTwoMock = NewMockClient(clientTwoCtrl)
		serviceMock = service.NewMockService(serviceCtrl)

		L = NewLobby(serviceMock)
	})

	BeforeEach(func() {
		clientOneCtrl.Finish()
		clientTwoCtrl.Finish()
	})

	AfterEach(func() {
		time.Sleep(100 * time.Millisecond)
	})

	It("should make first joiner a lobby owner", func() {
		// given
		clientOneMock.EXPECT().GetID().AnyTimes().Return(clientOneId)
		clientOneMock.EXPECT().GetNickname().AnyTimes().Return(playerOne.nickname)
		clientOneMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventLobby, shared.LobbyEventPlayerID, playerOne.id)).Times(1)
		clientOneMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventLobby, shared.LobbyEventPlayerJoined, playerOne.nickname)).Times(1)

		// when
		L.handleJoin(clientOneMock)

		// then
		Expect(L.GetOwnerID()).To(Equal(clientOneId))
	})

	It("new player joined lobby", func() {
		// given
		clientOneMock.EXPECT().GetID().AnyTimes().Return(clientOneId)

		clientOneMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventLobby, shared.LobbyEventPlayerJoined, playerTwo.nickname)).Times(1)

		clientTwoMock.EXPECT().GetID().AnyTimes().Return(clientTwoId)
		clientTwoMock.EXPECT().GetNickname().AnyTimes().Return(playerTwo.nickname)
		clientTwoMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventLobby, shared.LobbyEventPlayerID, playerTwo.id)).Times(1)
		clientTwoMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventLobby, shared.LobbyEventPlayerJoined, playerTwo.nickname)).Times(1)

		// when
		L.handleJoin(clientTwoMock)

		// then
		//	Expect(L.GetOwnerID()).To(Equal(clientTwoId))
	})

	It("second player shouldn't start a game", func() {
		// given && when && then
		L.handleLobbyMessage(models.Message{
			SenderID: clientTwoId,
			MessageDetails: models.MessageDetails{
				Type:   models.MessageTypeLobby,
				Action: models.ActionTypeStartGame,
			},
		})
	})

	It("first player should start a game", func() {
		// given
		clientOneMock.EXPECT().GetReceiver().Times(1)
		clientTwoMock.EXPECT().GetReceiver().Times(1)

		clientOneMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventLobby, shared.LobbyEventStartOfGame, "")).Times(1)
		clientOneMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventGame, shared.CommonGameEventScoreboard, "{\"clientOne\":0,\"clientTwo\":0}")).Times(1)
		clientOneMock.EXPECT().SendMessage(gomock.Any()).Times(1)
		clientTwoMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventLobby, shared.LobbyEventStartOfGame, "")).Times(1)
		clientTwoMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventGame, shared.CommonGameEventScoreboard, "{\"clientOne\":0,\"clientTwo\":0}")).Times(1)
		clientTwoMock.EXPECT().SendMessage(gomock.Any()).Times(1)

		// when
		L.handleLobbyMessage(models.Message{
			SenderID: clientOneId,
			MessageDetails: models.MessageDetails{
				Type:   models.MessageTypeLobby,
				Action: models.ActionTypeStartGame,
			},
		})
	})

	It("first player correct guessed answers", func() {
		// given
		for i := range 9 {
			clientOneMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventGame, math_operations_events.MathOperationsEventQuestion, gomock.Any().String())).Times(1)
			clientOneMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventGame, shared.CommonGameEventScoreboard, fmt.Sprintf("{\"clientOne\":%d,\"clientTwo\":0}", i+1))).Times(1)
			clientTwoMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventGame, math_operations_events.MathOperationsEventQuestion, gomock.Any().String())).Times(1)
			clientTwoMock.EXPECT().SendMessage(shared.CreateSocketResponse(shared.EventGame, shared.CommonGameEventScoreboard, fmt.Sprintf("{\"clientOne\":%d,\"clientTwo\":0}", i+1))).Times(1)

			answer := L.GetGame().GetRightAnswer(&i)
			data := math_operations.UserMessageData{
				Answer: answer,
			}
			byteData, _ := json.Marshal(data)

			// when
			L.handleMessage(models.Message{
				SenderID: clientOneId,
				MessageDetails: models.MessageDetails{
					Type:   models.MessageTypeGame,
					Action: models.ActionTypeGuessAnswer,
					Data:   string(byteData),
				},
			})
		}
	})

	// It("first player should miss answer", func() {
	// 	clientOneMock.EXPECT().SendMessage([]byte(`{"Type":"Scoreboard","Message":"{\"clientOne\":8,\"clientTwo\":0}"}`)).Times(1)
	// 	clientOneMock.EXPECT().SendMessage([]byte(`{"Type":"FinishedGame","Message":""}`)).Times(1)
	// 	clientTwoMock.EXPECT().SendMessage([]byte(`{"Type":"Scoreboard","Message":"{\"clientOne\":8,\"clientTwo\":0}"}`)).Times(1)

	// 	data := math_operations.UserMessageData{
	// 		Answer: "some-random-answer",
	// 	}
	// 	byteData, _ := json.Marshal(data)

	// 	// when
	// 	L.handleMessage(models.Message{
	// 		SenderID: clientOneId,
	// 		MessageDetails: models.MessageDetails{
	// 			Type:   models.MessageTypeGame,
	// 			Action: models.ActionTypeGuessAnswer,
	// 			Data:   string(byteData),
	// 		},
	// 	})
	// })
})
