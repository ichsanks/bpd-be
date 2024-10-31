package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type STtdService interface {
	Create(req RequestSTtd, userId uuid.UUID, tenantId uuid.UUID) (data STtd, err error)
	Update(req RequestSTtd, userId uuid.UUID, tenantId uuid.UUID) (data STtd, err error)
	GetAll(req model.StandardRequest) (data []STtdDto, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data STtd, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type STtdServiceImpl struct {
	STtdRepository STtdRepository
	Config         *configs.Config
}

func ProvideSTtdServiceImpl(repository STtdRepository, config *configs.Config) *STtdServiceImpl {
	s := new(STtdServiceImpl)
	s.STtdRepository = repository
	s.Config = config
	return s
}

func (s *STtdServiceImpl) Create(req RequestSTtd, userId uuid.UUID, tenantId uuid.UUID) (data STtd, err error) {

	data, _ = data.STtdFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.STtdRepository.Create(data)
	if err != nil {
		return STtd{}, err
	}
	return data, nil
}

func (s *STtdServiceImpl) GetAll(req model.StandardRequest) (data []STtdDto, err error) {
	return s.STtdRepository.GetAll(req)
}

func (s *STtdServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.STtdRepository.ResolveAll(request)
}

func (s *STtdServiceImpl) DeleteByID(id uuid.UUID) error {
	newSTtd, err := s.STtdRepository.ResolveByID(id)

	if err != nil || (STtd{}) == newSTtd {
		return errors.New("Data STtd dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.STtdRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data STtd dengan ID: " + id.String())
	}
	return nil
}

func (s *STtdServiceImpl) Update(req RequestSTtd, userId uuid.UUID, tenantId uuid.UUID) (data STtd, err error) {

	data, _ = data.STtdFormatRequest(req, userId, tenantId)
	err = s.STtdRepository.Update(data)
	if err != nil {
		return STtd{}, err
	}
	return data, nil
}

func (s *STtdServiceImpl) ResolveByID(id uuid.UUID) (data STtd, err error) {
	return s.STtdRepository.ResolveByID(id)
}

func (s *STtdServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newSTtd, err := s.STtdRepository.ResolveByID(id)

	if err != nil || (STtd{}) == newSTtd {
		return errors.New("Data STtd dengan ID :" + id.String() + " tidak ditemukan")
	}

	newSTtd.SoftDelete(userId)
	err = s.STtdRepository.Update(newSTtd)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data STtd dengan ID: " + id.String())
	}
	return nil
}
