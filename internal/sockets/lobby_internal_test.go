package sockets

import (
	"github.com/google/uuid"
	"mathly/internal/service"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var (
	serviceCtrl *gomock.Controller
	serviceMock *service.MockService

	clientOneCtrl *gomock.Controller
	clientOneMock *MockClient
)

var _ = Describe("Lobby", Ordered, func() {
	clientOneId, _ := uuid.Parse("10zcxzxc-3144-43fe-ab86-aab9f5c8de10")
	var L Lobby

	BeforeAll(func() {
		clientOneCtrl = gomock.NewController(GinkgoT())
		clientOneMock = NewMockClient(clientOneCtrl)

		serviceCtrl = gomock.NewController(GinkgoT())
		serviceMock = service.NewMockService(serviceCtrl)

		L = NewLobby(serviceMock)
	})

	It("should make first joiner a lobby owner", func() {
		// given
		clientOneMock.EXPECT().GetNickname().AnyTimes().Return("clientOne")
		clientOneMock.EXPECT().GetID().AnyTimes().Return(clientOneId)
		clientOneMock.EXPECT().SendMessage(gomock.Any()).AnyTimes()

		// when
		L.handleJoin(clientOneMock)

		// then
		Expect(L.GetOwnerID()).To(Equal(clientOneId))
	})
})
