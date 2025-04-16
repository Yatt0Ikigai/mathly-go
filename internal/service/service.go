package service

type ServiceConfig struct {
	JWT JwtConfig
}

type Service interface {
	JWT() JWT
}

type service struct {
	jwt JWT
}

func NewService(c ServiceConfig) Service {
	return &service{
		jwt: newJWT(c.JWT),
	}
}

func (s service) JWT() JWT {
	return s.jwt
}
