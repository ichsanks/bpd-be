package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type StatusPegawaiService interface {
	Create(req RequestStatusPegawai, userId uuid.UUID, tenantId uuid.UUID) (data StatusPegawai, err error)
	Update(req RequestStatusPegawai, userId uuid.UUID, tenantId uuid.UUID) (data StatusPegawai, err error)
	GetAll(req model.StandardRequest) (data []StatusPegawai, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data StatusPegawai, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type StatusPegawaiServiceImpl struct {
	StatusPegawaiRepository StatusPegawaiRepository
	Config                  *configs.Config
}

func ProvideStatusPegawaiServiceImpl(repository StatusPegawaiRepository, config *configs.Config) *StatusPegawaiServiceImpl {
	s := new(StatusPegawaiServiceImpl)
	s.StatusPegawaiRepository = repository
	s.Config = config
	return s
}

func (s *StatusPegawaiServiceImpl) Create(req RequestStatusPegawai, userId uuid.UUID, tenantId uuid.UUID) (data StatusPegawai, err error) {

	existNama, err := s.StatusPegawaiRepository.ExistByNama(req.Nama, "", *req.IdBranch)
	if existNama {
		return StatusPegawai{}, errors.New("Nama StatusPegawai sudah dipakai")
	}

	data, _ = data.StatusPegawaiFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.StatusPegawaiRepository.Create(data)
	if err != nil {
		return StatusPegawai{}, err
	}
	return data, nil
}

func (s *StatusPegawaiServiceImpl) GetAll(req model.StandardRequest) (data []StatusPegawai, err error) {
	return s.StatusPegawaiRepository.GetAll(req)
}

func (s *StatusPegawaiServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.StatusPegawaiRepository.ResolveAll(request)
}

func (s *StatusPegawaiServiceImpl) DeleteByID(id uuid.UUID) error {
	newStatusPegawai, err := s.StatusPegawaiRepository.ResolveByID(id)

	if err != nil || (StatusPegawai{}) == newStatusPegawai {
		return errors.New("Data StatusPegawai dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.StatusPegawaiRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data StatusPegawai dengan ID: " + id.String())
	}
	return nil
}

func (s *StatusPegawaiServiceImpl) Update(req RequestStatusPegawai, userId uuid.UUID, tenantId uuid.UUID) (data StatusPegawai, err error) {

	existNama, err := s.StatusPegawaiRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	if existNama {
		return StatusPegawai{}, errors.New("Nama StatusPegawai sudah dipakai")
	}

	data, _ = data.StatusPegawaiFormatRequest(req, userId, tenantId)
	err = s.StatusPegawaiRepository.Update(data)
	if err != nil {
		return StatusPegawai{}, err
	}
	return data, nil
}

func (s *StatusPegawaiServiceImpl) ResolveByID(id uuid.UUID) (data StatusPegawai, err error) {
	return s.StatusPegawaiRepository.ResolveByID(id)
}

func (s *StatusPegawaiServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newStatusPegawai, err := s.StatusPegawaiRepository.ResolveByID(id)

	if err != nil || (StatusPegawai{}) == newStatusPegawai {
		return errors.New("Data StatusPegawai dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.StatusPegawaiRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data StatusPegawai dengan Nama :" + *&newStatusPegawai.Nama + " sedang digunakan")
	}

	newStatusPegawai.SoftDelete(userId)
	err = s.StatusPegawaiRepository.Update(newStatusPegawai)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data StatusPegawai dengan ID: " + id.String())
	}
	return nil
}
