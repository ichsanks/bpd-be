package bpd

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type PerjalananDinasKendaraanService interface {
	Create(req RequestPerjalananDinasKendaraan, userId uuid.UUID) (data PerjalananDinasKendaraan, err error)
	Update(req RequestPerjalananDinasKendaraan, userId uuid.UUID) (data PerjalananDinasKendaraan, err error)
	GetAll(idPerjalananDinas string) (data []PerjalananDinasKendaraanDTO, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data PerjalananDinasKendaraan, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type PerjalananDinasKendaraanServiceImpl struct {
	PerjalananDinasKendaraanRepository PerjalananDinasKendaraanRepository
	Config                             *configs.Config
}

func ProvidePerjalananDinasKendaraanServiceImpl(repository PerjalananDinasKendaraanRepository, config *configs.Config) *PerjalananDinasKendaraanServiceImpl {
	s := new(PerjalananDinasKendaraanServiceImpl)
	s.PerjalananDinasKendaraanRepository = repository
	s.Config = config
	return s
}

func (s *PerjalananDinasKendaraanServiceImpl) Create(req RequestPerjalananDinasKendaraan, userId uuid.UUID) (data PerjalananDinasKendaraan, err error) {
	data, _ = data.PerjalananDinasKendaraanFormatRequest(req, userId)
	if err != nil {
		return
	}

	err = s.PerjalananDinasKendaraanRepository.Create(data)
	if err != nil {
		return PerjalananDinasKendaraan{}, err
	}
	return data, nil
}

func (s *PerjalananDinasKendaraanServiceImpl) GetAll(idPerjalananDinas string) (data []PerjalananDinasKendaraanDTO, err error) {
	return s.PerjalananDinasKendaraanRepository.GetAll(idPerjalananDinas)
}

func (s *PerjalananDinasKendaraanServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.PerjalananDinasKendaraanRepository.ResolveAll(request)
}

func (s *PerjalananDinasKendaraanServiceImpl) DeleteByID(id uuid.UUID) error {
	newKendaraan, err := s.PerjalananDinasKendaraanRepository.ResolveByID(id)

	if err != nil || (PerjalananDinasKendaraan{}) == newKendaraan {
		return errors.New("Data Perjalanan dinas kendaraan dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.PerjalananDinasKendaraanRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Perjalanan Dinas Kendaraan dengan ID: " + id.String())
	}
	return nil
}

func (s *PerjalananDinasKendaraanServiceImpl) Update(req RequestPerjalananDinasKendaraan, userId uuid.UUID) (data PerjalananDinasKendaraan, err error) {
	data, _ = data.PerjalananDinasKendaraanFormatRequest(req, userId)
	err = s.PerjalananDinasKendaraanRepository.Update(data)
	if err != nil {
		return PerjalananDinasKendaraan{}, err
	}
	return data, nil
}

func (s *PerjalananDinasKendaraanServiceImpl) ResolveByID(id uuid.UUID) (data PerjalananDinasKendaraan, err error) {
	return s.PerjalananDinasKendaraanRepository.ResolveByID(id)
}

func (s *PerjalananDinasKendaraanServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newKendaraan, err := s.PerjalananDinasKendaraanRepository.ResolveByID(id)

	if err != nil || (PerjalananDinasKendaraan{}) == newKendaraan {
		return errors.New("Data Perjalanan Dinas Kendaraan dengan ID :" + id.String() + " tidak ditemukan")
	}

	newKendaraan.SoftDelete(userId)
	err = s.PerjalananDinasKendaraanRepository.Update(newKendaraan)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Perjalanan dinas kendaraan dengan ID: " + id.String())
	}
	return nil
}
