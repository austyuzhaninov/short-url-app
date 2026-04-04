package service

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) Shorten() int64 {

	return 1
}

func (s *Service) GetShortCode() int64 {

	return 1
}

func (s *Service) GetStatus() int64 {

	return 1
}
