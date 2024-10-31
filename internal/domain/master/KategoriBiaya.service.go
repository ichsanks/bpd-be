package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type KategoriBiayaService interface {
	Create(req RequestKategoriBiaya, userId uuid.UUID, tenantId uuid.UUID) (data KategoriBiaya, err error)
	Update(req RequestKategoriBiaya, userId uuid.UUID, tenantId uuid.UUID) (data KategoriBiaya, err error)
	GetAll(req model.StandardRequest) (data []KategoriBiaya, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data KategoriBiaya, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type KategoriBiayaServiceImpl struct {
	KategoriBiayaRepository KategoriBiayaRepository
	Config                  *configs.Config
}

func ProvideKategoriBiayaServiceImpl(repository KategoriBiayaRepository, config *configs.Config) *KategoriBiayaServiceImpl {
	s := new(KategoriBiayaServiceImpl)
	s.KategoriBiayaRepository = repository
	s.Config = config
	return s
}

func (s *KategoriBiayaServiceImpl) Create(req RequestKategoriBiaya, userId uuid.UUID, tenantId uuid.UUID) (data KategoriBiaya, err error) {

	existNama, err := s.KategoriBiayaRepository.ExistByNama(req.Nama, "", *req.IdBranch)
	if existNama {
		return KategoriBiaya{}, errors.New("Nama KategoriBiaya sudah dipakai")
	}

	data, _ = data.KategoriBiayaFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.KategoriBiayaRepository.Create(data)
	if err != nil {
		return KategoriBiaya{}, err
	}
	return data, nil
}

func (s *KategoriBiayaServiceImpl) GetAll(req model.StandardRequest) (data []KategoriBiaya, err error) {
	return s.KategoriBiayaRepository.GetAll(req)
}

func (s *KategoriBiayaServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.KategoriBiayaRepository.ResolveAll(request)
}

func (s *KategoriBiayaServiceImpl) DeleteByID(id uuid.UUID) error {
	newKategoriBiaya, err := s.KategoriBiayaRepository.ResolveByID(id)

	if err != nil || (KategoriBiaya{}) == newKategoriBiaya {
		return errors.New("Data KategoriBiaya dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.KategoriBiayaRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data KategoriBiaya dengan ID: " + id.String())
	}
	return nil
}

func (s *KategoriBiayaServiceImpl) Update(req RequestKategoriBiaya, userId uuid.UUID, tenantId uuid.UUID) (data KategoriBiaya, err error) {

	existNama, err := s.KategoriBiayaRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	if existNama {
		return KategoriBiaya{}, errors.New("Nama KategoriBiaya sudah dipakai")
	}

	data, _ = data.KategoriBiayaFormatRequest(req, userId, tenantId)
	err = s.KategoriBiayaRepository.Update(data)
	if err != nil {
		return KategoriBiaya{}, err
	}
	return data, nil
}

func (s *KategoriBiayaServiceImpl) ResolveByID(id uuid.UUID) (data KategoriBiaya, err error) {
	return s.KategoriBiayaRepository.ResolveByID(id)
}

func (s *KategoriBiayaServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newKategoriBiaya, err := s.KategoriBiayaRepository.ResolveByID(id)

	if err != nil || (KategoriBiaya{}) == newKategoriBiaya {
		return errors.New("Data KategoriBiaya dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.KategoriBiayaRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data KategoriBiaya dengan Nama :" + *&newKategoriBiaya.Nama + " sedang digunakan")
	}

	newKategoriBiaya.SoftDelete(userId)
	err = s.KategoriBiayaRepository.Update(newKategoriBiaya)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data KategoriBiaya dengan ID: " + id.String())
	}
	return nil
}
