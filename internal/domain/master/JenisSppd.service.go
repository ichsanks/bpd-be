package master

import (
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
)

type JenisSppdService interface {
	GetAll() (data []JenisSppd, err error)
}

type JenisSppdServiceImpl struct {
	JenisSppdRepository JenisSppdRepository
	Config              *configs.Config
}

func ProvideJenisSppdServiceImpl(repository JenisSppdRepository) *JenisSppdServiceImpl {
	s := new(JenisSppdServiceImpl)
	s.JenisSppdRepository = repository
	return s
}

func (s *JenisSppdServiceImpl) GetAll() (data []JenisSppd, err error) {
	return s.JenisSppdRepository.GetAll()
}
