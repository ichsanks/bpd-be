package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type JenisBiayaService interface {
	Create(req JenisBiayaFormat, userId uuid.UUID, tenantId uuid.UUID) (newJenisBiaya JenisBiaya, err error)
	GetAll(req model.StandardRequest) (data []JenisBiaya, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	Update(req JenisBiayaFormat, userId uuid.UUID, tenantId uuid.UUID) (newJenisBiaya JenisBiaya, err error)
	ResolveByID(id uuid.UUID) (data JenisBiaya, err error)
	DeleteSoft(id uuid.UUID, userId uuid.UUID) error
	GetJumlahBiayaByIdBod(idBod uuid.UUID, ket string) (data JumlahBiaya, err error)
	GetAllDto(req model.StandardRequest) (data []JenisBiayaDto, err error)
	GetAllHeader() (data []JenisBiayaHeader, err error)
}

type JenisBiayaServiceImpl struct {
	JenisBiayaRepository JenisBiayaRepository
	Config               *configs.Config
}

func ProvideJenisBiayaServiceImpl(repository JenisBiayaRepository) *JenisBiayaServiceImpl {
	s := new(JenisBiayaServiceImpl)
	s.JenisBiayaRepository = repository
	return s
}

func (s *JenisBiayaServiceImpl) Create(req JenisBiayaFormat, userId uuid.UUID, tenantId uuid.UUID) (newJenisBiaya JenisBiaya, err error) {
	exist, err := s.JenisBiayaRepository.ExistByNama(req.Nama, *req.IdBranch)
	if exist {
		x := errors.New("Nama JenisBiaya sudah dipakai")
		return JenisBiaya{}, x
	}
	if err != nil {
		return JenisBiaya{}, err
	}
	newJenisBiaya, _ = newJenisBiaya.JenisBiayaFormat(req, userId, tenantId)
	err = s.JenisBiayaRepository.Create(newJenisBiaya)
	if err != nil {
		return JenisBiaya{}, err
	}
	return newJenisBiaya, nil
}

func (s *JenisBiayaServiceImpl) GetAll(req model.StandardRequest) (data []JenisBiaya, err error) {
	return s.JenisBiayaRepository.GetAll(req)
}

func (s *JenisBiayaServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.JenisBiayaRepository.ResolveAll(request)
}

func (s *JenisBiayaServiceImpl) DeleteByID(id uuid.UUID) error {
	_, err := s.JenisBiayaRepository.ResolveByID(id)
	if err != nil {
		return errors.New("Data JenisBiaya dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.JenisBiayaRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data JenisBiaya dengan ID: " + id.String())
	}
	return nil
}

func (s *JenisBiayaServiceImpl) Update(req JenisBiayaFormat, userId uuid.UUID, tenantId uuid.UUID) (newJenisBiaya JenisBiaya, err error) {
	exist, err := s.JenisBiayaRepository.ExistByNamaID(req.ID, req.Nama, *req.IdBranch)
	if exist {
		x := errors.New("Nama JenisBiaya sudah dipakai")
		return JenisBiaya{}, x
	}
	if err != nil {
		return JenisBiaya{}, err
	}
	newJenisBiaya, _ = newJenisBiaya.JenisBiayaFormat(req, userId, tenantId)
	err = s.JenisBiayaRepository.Update(newJenisBiaya)
	if err != nil {
		return JenisBiaya{}, err
	}
	return newJenisBiaya, nil
}

func (s *JenisBiayaServiceImpl) ResolveByID(id uuid.UUID) (newJenisBiaya JenisBiaya, err error) {
	newJenisBiaya, err = s.JenisBiayaRepository.ResolveByID(id)
	if err != nil {
		return JenisBiaya{}, errors.New("Data jenis biaya tidak ditemukan")
	}

	// komponenBiaya, err := s.JenisBiayaRepository.GetAllKomponenBiaya(id.String())
	// if err != nil {
	// 	return JenisBiaya{}, errors.New("Data komponen biaya tidak ditemukan")
	// }

	// newJenisBiaya.Detail = komponenBiaya

	return
}

func (s *JenisBiayaServiceImpl) DeleteSoft(id uuid.UUID, userId uuid.UUID) error {
	newJenisBiaya, err := s.JenisBiayaRepository.ResolveByID(id)

	if err != nil {
		return errors.New("Data JenisBiaya dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.JenisBiayaRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data Jenis Biaya dengan Nama :" + *&newJenisBiaya.Nama + " sedang digunakan")
	}

	newJenisBiaya.SoftDelete(userId)
	err = s.JenisBiayaRepository.Update(newJenisBiaya)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data JenisBiaya dengan ID: " + id.String())
	}
	return nil
}

func (s *JenisBiayaServiceImpl) GetJumlahBiayaByIdBod(idBod uuid.UUID, ket string) (data JumlahBiaya, err error) {
	return s.JenisBiayaRepository.GetJumlahBiayaByIdBod(idBod, ket)
}

func (s *JenisBiayaServiceImpl) GetAllDto(req model.StandardRequest) (data []JenisBiayaDto, err error) {
	return s.JenisBiayaRepository.GetAllDto(req)
}

func (s *JenisBiayaServiceImpl) GetAllHeader() (data []JenisBiayaHeader, err error) {
	return s.JenisBiayaRepository.GetAllHeader()
}
