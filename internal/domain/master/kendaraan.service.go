package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type KendaraanService interface {
	Create(req RequestKendaraan, userId uuid.UUID) (data Kendaraan, err error)
	Update(req RequestKendaraan, userId uuid.UUID) (data Kendaraan, err error)
	GetAll() (data []Kendaraan, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequestKendaraan) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Kendaraan, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type KendaraanServiceImpl struct {
	KendaraanRepository KendaraanRepository
	Config              *configs.Config
}

func ProvideKendaraanServiceImpl(repository KendaraanRepository, config *configs.Config) *KendaraanServiceImpl {
	s := new(KendaraanServiceImpl)
	s.KendaraanRepository = repository
	s.Config = config
	return s
}

func (s *KendaraanServiceImpl) Create(req RequestKendaraan, userId uuid.UUID) (data Kendaraan, err error) {
	existNopol, err := s.KendaraanRepository.ExistByNopol(req.Nopol, "")
	if existNopol {
		return Kendaraan{}, errors.New("Nopol Kendaraan sudah dipakai")
	}

	data, _ = data.KendaraanFormatRequest(req, userId)
	if err != nil {
		return
	}

	err = s.KendaraanRepository.Create(data)
	if err != nil {
		return Kendaraan{}, err
	}
	return data, nil
}

func (s *KendaraanServiceImpl) GetAll() (data []Kendaraan, err error) {
	return s.KendaraanRepository.GetAll()
}

func (s *KendaraanServiceImpl) ResolveAll(request model.StandardRequestKendaraan) (orders pagination.Response, err error) {
	return s.KendaraanRepository.ResolveAll(request)
}

func (s *KendaraanServiceImpl) DeleteByID(id uuid.UUID) error {
	newKendaraan, err := s.KendaraanRepository.ResolveByID(id)

	if err != nil || (Kendaraan{}) == newKendaraan {
		return errors.New("Data Kendaraan dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.KendaraanRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Kendaraan dengan ID: " + id.String())
	}
	return nil
}

func (s *KendaraanServiceImpl) Update(req RequestKendaraan, userId uuid.UUID) (data Kendaraan, err error) {
	existNopol, err := s.KendaraanRepository.ExistByNopol(req.Nopol, req.ID.String())
	if existNopol {
		return Kendaraan{}, errors.New("Nopol Kendaraan sudah dipakai")
	}
	data, _ = data.KendaraanFormatRequest(req, userId)
	err = s.KendaraanRepository.Update(data)
	if err != nil {
		return Kendaraan{}, err
	}
	return data, nil
}

func (s *KendaraanServiceImpl) ResolveByID(id uuid.UUID) (data Kendaraan, err error) {
	return s.KendaraanRepository.ResolveByID(id)
}

func (s *KendaraanServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newKendaraan, err := s.KendaraanRepository.ResolveByID(id)

	if err != nil || (Kendaraan{}) == newKendaraan {
		return errors.New("Data Kendaraan dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.KendaraanRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data Kendaraan dengan Nama :" + *&newKendaraan.Nama + " sedang digunakan")
	}

	newKendaraan.SoftDelete(userId)
	err = s.KendaraanRepository.Update(newKendaraan)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Kendaraan dengan ID: " + id.String())
	}
	return nil
}
