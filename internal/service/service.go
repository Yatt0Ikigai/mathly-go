//go:generate mockgen -source=service.go -package service -destination=service_mock.go
package service

type ServiceConfig struct {
	JWT JwtConfig
}

type Service interface {
	JWT() JWT
	LobbyHandler() LobbyHandler
}

type service struct {
	jwt          JWT
	lobbyHandler LobbyHandler
}

func NewService(c ServiceConfig) Service {
	return &service{
		jwt:          newJWT(c.JWT),
		lobbyHandler: newLobbyHandler(),
	}
}

func (s service) JWT() JWT {
	return s.jwt
}

func (s service) LobbyHandler() LobbyHandler {
	return s.lobbyHandler
}
