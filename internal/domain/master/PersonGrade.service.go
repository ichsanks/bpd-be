package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type PersonGradeService interface {
	Create(req RequestPersonGrade, userId uuid.UUID, tenanId uuid.UUID) (data PersonGrade, err error)
	Update(req RequestPersonGrade, userId uuid.UUID, tenanId uuid.UUID) (data PersonGrade, err error)
	GetAll(req model.StandardRequest) (data []PersonGrade, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data PersonGrade, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID, idBranch string) error
}

type PersonGradeServiceImpl struct {
	PersonGradeRepository PersonGradeRepository
	Config                *configs.Config
}

func ProvidePersonGradeServiceImpl(repository PersonGradeRepository, config *configs.Config) *PersonGradeServiceImpl {
	s := new(PersonGradeServiceImpl)
	s.PersonGradeRepository = repository
	s.Config = config
	return s
}

func (s *PersonGradeServiceImpl) Create(req RequestPersonGrade, userId uuid.UUID, tenanId uuid.UUID) (data PersonGrade, err error) {
	existKode, err := s.PersonGradeRepository.ExistByKode(req.Kode, "", *req.IdBranch)
	if existKode {
		return PersonGrade{}, errors.New("Kode PersonGrade sudah dipakai")
	}

	existNama, err := s.PersonGradeRepository.ExistByNama(req.Nama, "", *req.IdBranch)
	if existNama {
		return PersonGrade{}, errors.New("Nama PersonGrade sudah dipakai")
	}

	data, _ = data.PersonGradeFormatRequest(req, userId, tenanId)
	if err != nil {
		return
	}

	err = s.PersonGradeRepository.Create(data)
	if err != nil {
		return PersonGrade{}, err
	}
	return data, nil
}

func (s *PersonGradeServiceImpl) GetAll(req model.StandardRequest) (data []PersonGrade, err error) {
	return s.PersonGradeRepository.GetAll(req)
}

func (s *PersonGradeServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.PersonGradeRepository.ResolveAll(request)
}

func (s *PersonGradeServiceImpl) DeleteByID(id uuid.UUID) error {
	newPersonGrade, err := s.PersonGradeRepository.ResolveByID(id)

	if err != nil || (PersonGrade{}) == newPersonGrade {
		return errors.New("Data PersonGrade dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.PersonGradeRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data PersonGrade dengan ID: " + id.String())
	}
	return nil
}

func (s *PersonGradeServiceImpl) Update(req RequestPersonGrade, userId uuid.UUID, tenanId uuid.UUID) (data PersonGrade, err error) {
	existKode, err := s.PersonGradeRepository.ExistByKode(req.Kode, req.ID.String(), *req.IdBranch)
	if existKode {
		return PersonGrade{}, errors.New("Kode PersonGrade sudah dipakai")
	}

	existNama, err := s.PersonGradeRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	if existNama {
		return PersonGrade{}, errors.New("Nama PersonGrade sudah dipakai")
	}

	data, _ = data.PersonGradeFormatRequest(req, userId, tenanId)
	err = s.PersonGradeRepository.Update(data)
	if err != nil {
		return PersonGrade{}, err
	}
	return data, nil
}

func (s *PersonGradeServiceImpl) ResolveByID(id uuid.UUID) (data PersonGrade, err error) {
	return s.PersonGradeRepository.ResolveByID(id)
}

func (s *PersonGradeServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID, idBranch string) error {
	newPersonGrade, err := s.PersonGradeRepository.ResolveByID(id)

	if err != nil || (PersonGrade{}) == newPersonGrade {
		return errors.New("Data PersonGrade dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.PersonGradeRepository.ExistRelasiStatus(id, idBranch)
	if exist {
		return errors.New(" Data PersonGrade dengan Nama :" + *&newPersonGrade.Nama + " sedang digunakan")
	}

	newPersonGrade.SoftDelete(userId)
	err = s.PersonGradeRepository.Update(newPersonGrade)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data PersonGrade dengan ID: " + id.String())
	}
	return nil
}
