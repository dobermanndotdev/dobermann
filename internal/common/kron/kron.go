package kron

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s Service) Stop() error {
	return nil
}

func (s Service) Start() error {
	return nil
}
