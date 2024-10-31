package master

import "gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"

type JenisKelaminService interface {
	GetAll() (newJenisKelamin []JenisKelamin, err error)
}

type JenisKelaminServiceImpl struct {
	JenisKelaminRepository JenisKelaminRepository
	Config                 *configs.Config
}

func ProvideJenisKelaminServiceImpl(repository JenisKelaminRepository, config *configs.Config) *JenisKelaminServiceImpl {
	s := new(JenisKelaminServiceImpl)
	s.JenisKelaminRepository = repository
	s.Config = config
	return s
}

func (s *JenisKelaminServiceImpl) GetAll() (newJenisKelamin []JenisKelamin, err error) {
	return s.JenisKelaminRepository.GetAll()
}
