package master

import "gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"

type AgamaService interface {
	GetAll() (newAgama []Agama, err error)
}

type AgamaServiceImpl struct {
	AgamaRepository AgamaRepository
	Config          *configs.Config
}

func ProvideAgamaServiceImpl(repository AgamaRepository, config *configs.Config) *AgamaServiceImpl {
	s := new(AgamaServiceImpl)
	s.AgamaRepository = repository
	s.Config = config
	return s
}

func (s *AgamaServiceImpl) GetAll() (newAgama []Agama, err error) {
	return s.AgamaRepository.GetAll()
}
