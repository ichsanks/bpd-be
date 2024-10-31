package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type FasilitasTransportService interface {
	Create(req RequestFasilitasTransport, userId uuid.UUID, tenantId uuid.UUID) (data FasilitasTransport, err error)
	Update(req RequestFasilitasTransport, userId uuid.UUID, tenantId uuid.UUID) (data FasilitasTransport, err error)
	GetAll(req model.StandardRequest) (data []FasilitasTransport, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data FasilitasTransport, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type FasilitasTransportServiceImpl struct {
	FasilitasTransportRepository FasilitasTransportRepository
	Config                       *configs.Config
}

func ProvideFasilitasTransportServiceImpl(repository FasilitasTransportRepository, config *configs.Config) *FasilitasTransportServiceImpl {
	s := new(FasilitasTransportServiceImpl)
	s.FasilitasTransportRepository = repository
	s.Config = config
	return s
}

func (s *FasilitasTransportServiceImpl) Create(req RequestFasilitasTransport, userId uuid.UUID, tenantId uuid.UUID) (data FasilitasTransport, err error) {
	existNama, err := s.FasilitasTransportRepository.ExistByNama(req.Nama, "", *req.IdBranch)
	if existNama {
		return FasilitasTransport{}, errors.New("Nama FasilitasTransport sudah dipakai")
	}

	data, _ = data.FasilitasTransportFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.FasilitasTransportRepository.Create(data)
	if err != nil {
		return FasilitasTransport{}, err
	}
	return data, nil
}

func (s *FasilitasTransportServiceImpl) GetAll(req model.StandardRequest) (data []FasilitasTransport, err error) {
	return s.FasilitasTransportRepository.GetAll(req)
}

func (s *FasilitasTransportServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.FasilitasTransportRepository.ResolveAll(request)
}

func (s *FasilitasTransportServiceImpl) DeleteByID(id uuid.UUID) error {
	newFasilitasTransport, err := s.FasilitasTransportRepository.ResolveByID(id)

	if err != nil || (FasilitasTransport{}) == newFasilitasTransport {
		return errors.New("Data FasilitasTransport dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.FasilitasTransportRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data FasilitasTransport dengan ID: " + id.String())
	}
	return nil
}

func (s *FasilitasTransportServiceImpl) Update(req RequestFasilitasTransport, userId uuid.UUID, tenantId uuid.UUID) (data FasilitasTransport, err error) {
	existNama, err := s.FasilitasTransportRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	if existNama {
		return FasilitasTransport{}, errors.New("Nama FasilitasTransport sudah dipakai")
	}

	data, _ = data.FasilitasTransportFormatRequest(req, userId, tenantId)
	err = s.FasilitasTransportRepository.Update(data)
	if err != nil {
		return FasilitasTransport{}, err
	}
	return data, nil
}

func (s *FasilitasTransportServiceImpl) ResolveByID(id uuid.UUID) (data FasilitasTransport, err error) {
	return s.FasilitasTransportRepository.ResolveByID(id)
}

func (s *FasilitasTransportServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newFasilitasTransport, err := s.FasilitasTransportRepository.ResolveByID(id)

	if err != nil || (FasilitasTransport{}) == newFasilitasTransport {
		return errors.New("Data FasilitasTransport dengan ID :" + id.String() + " tidak ditemukan")
	}

	newFasilitasTransport.SoftDelete(userId)
	err = s.FasilitasTransportRepository.Update(newFasilitasTransport)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data FasilitasTransport dengan ID: " + id.String())
	}
	return nil
}
