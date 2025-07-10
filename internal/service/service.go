//go:generate mockgen -source=service.go -package service -destination=service_mock.go
package service

type ServiceConfig struct {
	JWT JwtConfig
}

type Service interface {
	JWT() JWT
	LobbyHandler() LobbyHandler
	Random() Random
}

type service struct {
	jwt          JWT
	lobbyHandler LobbyHandler
	random       Random
}

func NewService(c ServiceConfig) Service {
	return &service{
		jwt:          newJWT(c.JWT),
		lobbyHandler: newLobbyHandler(),
		random:       newRandom(),
	}
}

func (s service) JWT() JWT {
	return s.jwt
}

func (s service) LobbyHandler() LobbyHandler {
	return s.lobbyHandler
}

func (s service) Random() Random {
	return s.random
}
