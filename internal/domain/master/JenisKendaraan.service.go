package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type JenisKendaraanService interface {
	Create(req JenisKendaraanFormat, userId uuid.UUID) (newJenisKendaraan JenisKendaraan, err error)
	GetAll() (data []JenisKendaraan, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	Update(req JenisKendaraanFormat, userId uuid.UUID) (newJenisKendaraan JenisKendaraan, err error)
	ResolveByID(id uuid.UUID) (data JenisKendaraan, err error)
	DeleteSoft(id uuid.UUID, userId uuid.UUID) error
}

type JenisKendaraanServiceImpl struct {
	JenisKendaraanRepository JenisKendaraanRepository
	Config                   *configs.Config
}

func ProvideJenisKendaraanServiceImpl(repository JenisKendaraanRepository) *JenisKendaraanServiceImpl {
	s := new(JenisKendaraanServiceImpl)
	s.JenisKendaraanRepository = repository
	return s
}

func (s *JenisKendaraanServiceImpl) Create(req JenisKendaraanFormat, userId uuid.UUID) (newJenisKendaraan JenisKendaraan, err error) {
	exist, err := s.JenisKendaraanRepository.ExistByNama(req.Nama)
	if exist {
		x := errors.New("Nama JenisKendaraan sudah dipakai")
		return JenisKendaraan{}, x
	}
	if err != nil {
		return JenisKendaraan{}, err
	}
	newJenisKendaraan, _ = newJenisKendaraan.JenisKendaraanFormat(req, userId)
	err = s.JenisKendaraanRepository.Create(newJenisKendaraan)
	if err != nil {
		return JenisKendaraan{}, err
	}
	return newJenisKendaraan, nil
}

func (s *JenisKendaraanServiceImpl) GetAll() (data []JenisKendaraan, err error) {
	return s.JenisKendaraanRepository.GetAll()
}

func (s *JenisKendaraanServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.JenisKendaraanRepository.ResolveAll(request)
}

func (s *JenisKendaraanServiceImpl) DeleteByID(id uuid.UUID) error {
	jenisKendaraan, err := s.JenisKendaraanRepository.ResolveByID(id)

	if err != nil || (JenisKendaraan{}) == jenisKendaraan {
		return errors.New("Data JenisKendaraan dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.JenisKendaraanRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data JenisKendaraan dengan ID: " + id.String())
	}
	return nil
}

func (s *JenisKendaraanServiceImpl) Update(req JenisKendaraanFormat, userId uuid.UUID) (newJenisKendaraan JenisKendaraan, err error) {
	exist, err := s.JenisKendaraanRepository.ExistByNamaID(req.ID, req.Nama)
	if exist {
		x := errors.New("Nama JenisKendaraan sudah dipakai")
		return JenisKendaraan{}, x
	}
	if err != nil {
		return JenisKendaraan{}, err
	}
	newJenisKendaraan, _ = newJenisKendaraan.JenisKendaraanFormat(req, userId)
	err = s.JenisKendaraanRepository.Update(newJenisKendaraan)
	if err != nil {
		return JenisKendaraan{}, err
	}
	return newJenisKendaraan, nil
}

func (s *JenisKendaraanServiceImpl) ResolveByID(id uuid.UUID) (newJenisKendaraan JenisKendaraan, err error) {
	return s.JenisKendaraanRepository.ResolveByID(id)
}

func (s *JenisKendaraanServiceImpl) DeleteSoft(id uuid.UUID, userId uuid.UUID) error {
	newJenisKendaraan, err := s.JenisKendaraanRepository.ResolveByID(id)

	if err != nil || (JenisKendaraan{}) == newJenisKendaraan {
		return errors.New("Data JenisKendaraan dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.JenisKendaraanRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New("Data Jenis Kendaraan dengan Nama :" + *&newJenisKendaraan.Nama + " sedang digunakan")
	}

	newJenisKendaraan.SoftDelete(userId)
	err = s.JenisKendaraanRepository.Update(newJenisKendaraan)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data JenisKendaraan dengan ID: " + id.String())
	}
	return nil
}
