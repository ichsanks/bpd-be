package auth

import "gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"

type DashboardService interface {
	GetAll() (newDashboard []Dashboard, err error)
	GetDataDashboardBpd(req DashboardRequest) (data DashboardBpd, err error)
	GetDataBpd(req DashboardRequest) (data []DataAktifBpd, err error)
	GetDataDashboardSppd(req DashboardRequest) (data DashboardSppd, err error)
	GetDataDashboardBpdNew(req DashboardRequest) (data DashboardSppd, err error)
	GetDataSppd(req DashboardRequest) (data []DataAktifSppd, err error)
	GetDataBpdNew(req DashboardRequest) (data []DataAktifBpdNew, err error)
	GetJumlahSppd(req DashboardRequest) (data JumlahSppd, err error)
	GetJumlahBpd(req DashboardRequest) (data JumlahSppd, err error)
}

type DashboardServiceImpl struct {
	DashboardRepository DashboardRepository
	Config              *configs.Config
}

func ProvideDashboardServiceImpl(repository DashboardRepository, config *configs.Config) *DashboardServiceImpl {
	s := new(DashboardServiceImpl)
	s.DashboardRepository = repository
	s.Config = config
	return s
}

func (s *DashboardServiceImpl) GetAll() (newDashboard []Dashboard, err error) {
	return s.DashboardRepository.GetAll()
}

func (s *DashboardServiceImpl) GetDataDashboardBpd(req DashboardRequest) (data DashboardBpd, err error) {
	return s.DashboardRepository.GetDataDashboardBpd(req)
}

func (s *DashboardServiceImpl) GetDataBpd(req DashboardRequest) (data []DataAktifBpd, err error) {
	return s.DashboardRepository.GetDataBpd(req)
}

func (s *DashboardServiceImpl) GetDataDashboardSppd(req DashboardRequest) (data DashboardSppd, err error) {
	return s.DashboardRepository.GetDataDashboardSppd(req)
}

func (s *DashboardServiceImpl) GetDataSppd(req DashboardRequest) (data []DataAktifSppd, err error) {
	return s.DashboardRepository.GetDataSppd(req)
}

func (s *DashboardServiceImpl) GetDataDashboardBpdNew(req DashboardRequest) (data DashboardSppd, err error) {
	return s.DashboardRepository.GetDataDashboardBpdNew(req)
}

func (s *DashboardServiceImpl) GetDataBpdNew(req DashboardRequest) (data []DataAktifBpdNew, err error) {
	return s.DashboardRepository.GetDataBpdNew(req)
}

func (s *DashboardServiceImpl) GetJumlahSppd(req DashboardRequest) (data JumlahSppd, err error) {
	return s.DashboardRepository.GetJumlahSppd(req)
}

func (s *DashboardServiceImpl) GetJumlahBpd(req DashboardRequest) (data JumlahSppd, err error) {
	return s.DashboardRepository.GetJumlahBpd(req)
}
