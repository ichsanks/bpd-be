package master

import (
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type LogService interface {
	ResolveAll(request model.StandardRequestPegawai) (orders pagination.Response, err error)
}

type LogServiceImpl struct {
	LogRepository LogRepository
	Config        *configs.Config
}

func ProvideLogServiceImpl(repository LogRepository) *LogServiceImpl {
	s := new(LogServiceImpl)
	s.LogRepository = repository
	return s
}

func (s *LogServiceImpl) ResolveAll(request model.StandardRequestPegawai) (orders pagination.Response, err error) {
	return s.LogRepository.ResolveAll(request)
}
