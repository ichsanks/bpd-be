package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type UnitKerjaService interface {
	Create(req RequestUnor, userId uuid.UUID, tenantId uuid.UUID) (data UnitOrganisasiKerja, err error)
	Update(req RequestUnor, userId uuid.UUID, tenantId uuid.UUID) (data UnitOrganisasiKerja, err error)
	GetAll(req model.StandardRequest) (data []UnitOrganisasiKerjaDTO, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data UnitOrganisasiKerja, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type UnitKerjaServiceImpl struct {
	UnitKerjaRepository UnitKerjaRepository
	Config              *configs.Config
}

func ProvideUnitKerjaServiceImpl(repository UnitKerjaRepository, config *configs.Config) *UnitKerjaServiceImpl {
	s := new(UnitKerjaServiceImpl)
	s.UnitKerjaRepository = repository
	s.Config = config
	return s
}

func (s *UnitKerjaServiceImpl) Create(req RequestUnor, userId uuid.UUID, tenantId uuid.UUID) (data UnitOrganisasiKerja, err error) {
	existKode, _ := s.UnitKerjaRepository.ExistByKode(req.Kode, "", *req.IdBranch)
	if existKode {
		return UnitOrganisasiKerja{}, errors.New("Kode Unor sudah dipakai")
	}

	existNama, _ := s.UnitKerjaRepository.ExistByNama(req.Nama, "", *req.IdBranch)
	if existNama {
		return UnitOrganisasiKerja{}, errors.New("Nama Unor sudah dipakai")
	}

	data, _ = data.UnorFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.UnitKerjaRepository.Create(data)
	if err != nil {
		return UnitOrganisasiKerja{}, err
	}
	return data, nil
}

func (s *UnitKerjaServiceImpl) GetAll(req model.StandardRequest) (data []UnitOrganisasiKerjaDTO, err error) {
	return s.UnitKerjaRepository.GetAll(req)
}

func (s *UnitKerjaServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.UnitKerjaRepository.ResolveAll(request)
}

func (s *UnitKerjaServiceImpl) DeleteByID(id uuid.UUID) error {
	newUnor, err := s.UnitKerjaRepository.ResolveByID(id)

	if err != nil || (UnitOrganisasiKerja{}) == newUnor {
		return errors.New("Data Unor dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.UnitKerjaRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Unor dengan ID: " + id.String())
	}
	return nil
}

func (s *UnitKerjaServiceImpl) Update(req RequestUnor, userId uuid.UUID, tenantId uuid.UUID) (data UnitOrganisasiKerja, err error) {
	existKode, _ := s.UnitKerjaRepository.ExistByKode(req.Kode, req.ID.String(), *req.IdBranch)
	if existKode {
		return UnitOrganisasiKerja{}, errors.New("Kode Unor sudah dipakai")
	}

	existNama, _ := s.UnitKerjaRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	if existNama {
		return UnitOrganisasiKerja{}, errors.New("Nama Unor sudah dipakai")
	}

	data, _ = data.UnorFormatRequest(req, userId, tenantId)
	err = s.UnitKerjaRepository.Update(data)
	if err != nil {
		return UnitOrganisasiKerja{}, err
	}
	return data, nil
}

func (s *UnitKerjaServiceImpl) ResolveByID(id uuid.UUID) (data UnitOrganisasiKerja, err error) {
	return s.UnitKerjaRepository.ResolveByID(id)
}

func (s *UnitKerjaServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newUnor, err := s.UnitKerjaRepository.ResolveByID(id)

	if err != nil || (UnitOrganisasiKerja{}) == newUnor {
		return errors.New("Data Unor dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.UnitKerjaRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data Unor dengan Nama :" + *&newUnor.Nama + " sedang digunakan")
	}

	newUnor.SoftDelete(userId)
	err = s.UnitKerjaRepository.Update(newUnor)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Unor dengan ID: " + id.String())
	}
	return nil
}
