package health

type Service interface {
	Check() string
}

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) Check() string {
	return "service is healthy"
}
