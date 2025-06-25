package sockets

import (
	"encoding/json"
	"fmt"
	"mathly/internal/models"
	"mathly/internal/service"
	"mathly/internal/sockets/games/math_operations"
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
	clientOneId, _ := uuid.Parse("10zcxzxc-3144-43fe-ab86-aab9f5c8de10")
	clientTwoId, _ := uuid.Parse("asdaxzxc-3144-43fe-ab86-aab9f5c8de10")
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
		time.Sleep(1 * time.Second)
	})

	It("should make first joiner a lobby owner", func() {
		// given
		clientOneMock.EXPECT().GetNickname().AnyTimes().Return("clientOne")
		clientOneMock.EXPECT().GetID().AnyTimes().Return(clientOneId)
		clientOneMock.EXPECT().SendMessage([]byte("10000000-0000-0000-0000-000000000000")).Times(1)
		clientOneMock.EXPECT().SendMessage([]byte("New Player clientOne Joined")).Times(1)

		// when
		L.handleJoin(clientOneMock)

		// then
		Expect(L.GetOwnerID()).To(Equal(clientOneId))
	})

	It("new player joined lobby", func() {
		// given
		clientOneMock.EXPECT().SendMessage([]byte("New Player clientTwo Joined")).Times(1)

		clientTwoMock.EXPECT().GetID().AnyTimes().Return(clientTwoId)
		clientTwoMock.EXPECT().GetNickname().AnyTimes().Return("clientTwo")
		clientTwoMock.EXPECT().SendMessage([]byte("00000000-0000-0000-0000-000000000000")).Times(1)
		clientTwoMock.EXPECT().SendMessage([]byte("New Player clientTwo Joined")).Times(1)

		// when
		L.handleJoin(clientTwoMock)

		// then
		Expect(L.GetOwnerID()).To(Equal(clientOneId))
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

		clientOneMock.EXPECT().SendMessage([]byte(`{"Type":"StartOfGame","Message":""}`)).Times(1)
		clientOneMock.EXPECT().SendMessage([]byte(`{"Type":"Scoreboard","Message":"{\"clientOne\":0,\"clientTwo\":0}"}`)).Times(1)
	//	clientOneMock.EXPECT().SendMessage(gomock.Any()).Times(1)
		clientTwoMock.EXPECT().SendMessage([]byte(`{"Type":"StartOfGame","Message":""}`)).Times(1)
		clientTwoMock.EXPECT().SendMessage([]byte(`{"Type":"Scoreboard","Message":"{\"clientOne\":0,\"clientTwo\":0}"}`)).Times(1)
	//	clientTwoMock.EXPECT().SendMessage(gomock.Any()).Times(1)

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
			clientOneMock.EXPECT().SendMessage([]byte(fmt.Sprintf(`{"Type":"Scoreboard","Message":"{\"clientOne\":%d,\"clientTwo\":0}"}`, i+1))).Times(1)
			clientOneMock.EXPECT().SendMessage(gomock.Any()).Times(1)
			clientTwoMock.EXPECT().SendMessage([]byte(fmt.Sprintf(`{"Type":"Scoreboard","Message":"{\"clientOne\":%d,\"clientTwo\":0}"}`, i+1))).Times(1)

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

	It("first player should miss answer", func() {
		clientOneMock.EXPECT().SendMessage([]byte(`{"Type":"Scoreboard","Message":"{\"clientOne\":8,\"clientTwo\":0}"}`)).Times(1)
		clientOneMock.EXPECT().SendMessage([]byte(`{"Type":"FinishedGame","Message":""}`)).Times(1)
		clientTwoMock.EXPECT().SendMessage([]byte(`{"Type":"Scoreboard","Message":"{\"clientOne\":8,\"clientTwo\":0}"}`)).Times(1)

		data := math_operations.UserMessageData{
			Answer: "some-random-answer",
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
	})
})
