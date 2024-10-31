package master

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type PegawaiService interface {
	Create(req PegawaiFormat, userId uuid.UUID, tenantId uuid.UUID) (newPegawai Pegawai, err error)
	GetAll(req model.StandardRequest) (data []Pegawai, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	Update(req PegawaiFormat, userId uuid.UUID, tenantId uuid.UUID) (newPegawai Pegawai, err error)
	ResolveByID(id uuid.UUID) (newPegawai PegawaiDTO, err error)
	DeleteSoft(id uuid.UUID, userId uuid.UUID, idBranch string) error
	UploadFile(w http.ResponseWriter, r *http.Request, formValue string, pathFile string) (path string, err error)
	DeleteFile(path string) (err error)
}

type PegawaiServiceImpl struct {
	PegawaiRepository PegawaiRepository
	Config            *configs.Config
}

func ProvidePegawaiServiceImpl(repository PegawaiRepository, config *configs.Config) *PegawaiServiceImpl {
	s := new(PegawaiServiceImpl)
	s.PegawaiRepository = repository
	s.Config = config
	return s
}

func (s *PegawaiServiceImpl) Create(req PegawaiFormat, userId uuid.UUID, tenantId uuid.UUID) (newPegawai Pegawai, err error) {
	exist, err := s.PegawaiRepository.ExistByNip(req.Nip, *req.IdBranch)
	if exist {
		x := errors.New("Nip Pegawai sudah dipakai")
		return Pegawai{}, x
	}
	if err != nil {
		return Pegawai{}, err
	}

	existNama, err := s.PegawaiRepository.ExistByNama(req.Nama, *req.IdBranch)
	if existNama {
		x := errors.New("Nama Pegawai sudah dipakai")
		return Pegawai{}, x
	}
	if err != nil {
		return Pegawai{}, err
	}
	newPegawai, _ = newPegawai.PegawaiFormat(req, userId, tenantId)
	err = s.PegawaiRepository.Create(newPegawai)
	if err != nil {
		return Pegawai{}, err
	}
	return newPegawai, nil
}

func (s *PegawaiServiceImpl) GetAll(req model.StandardRequest) (data []Pegawai, err error) {
	return s.PegawaiRepository.GetAll(req)
}

func (s *PegawaiServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.PegawaiRepository.ResolveAll(request)
}

func (s *PegawaiServiceImpl) DeleteByID(id uuid.UUID) error {
	pegawai, err := s.PegawaiRepository.ResolveByID(id)

	if err != nil || (Pegawai{}) == pegawai {
		return errors.New("Data Pegawai dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.PegawaiRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Pegawai dengan ID: " + id.String())
	}
	return nil
}

func (s *PegawaiServiceImpl) Update(req PegawaiFormat, userId uuid.UUID, tenantId uuid.UUID) (newPegawai Pegawai, err error) {
	exist, err := s.PegawaiRepository.ExistByNipID(req.ID, req.Nip, *req.IdBranch)
	if exist {
		x := errors.New("Nip Pegawai sudah dipakai")
		return Pegawai{}, x
	}
	if err != nil {
		return Pegawai{}, err
	}

	existNama, err := s.PegawaiRepository.ExistByNamaID(req.ID, req.Nama, *req.IdBranch)
	if existNama {
		x := errors.New("Nama Pegawai sudah dipakai")
		return Pegawai{}, x
	}
	if err != nil {
		return Pegawai{}, err
	}
	newPegawai, _ = newPegawai.PegawaiFormat(req, userId, tenantId)
	err = s.PegawaiRepository.Update(newPegawai)
	if err != nil {
		return Pegawai{}, err
	}
	return newPegawai, nil
}

func (s *PegawaiServiceImpl) ResolveByID(id uuid.UUID) (newPegawai PegawaiDTO, err error) {
	return s.PegawaiRepository.ResolveByIDDTO(id)
}

func (s *PegawaiServiceImpl) DeleteSoft(id uuid.UUID, userId uuid.UUID, idBranch string) error {
	newPegawai, err := s.PegawaiRepository.ResolveByID(id)

	if err != nil || (Pegawai{}) == newPegawai {
		return errors.New("Data Pegawai dengan ID :" + id.String() + " tidak ditemukan")
	}
	exist := s.PegawaiRepository.ExistRelasiStatus(id, idBranch)
	if exist {
		return errors.New(" Data Pegawai dengan Nama :" + *&newPegawai.Nama + " sedang digunakan")
	}

	newPegawai.SoftDelete(userId)
	err = s.PegawaiRepository.Update(newPegawai)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Pegawai dengan ID: " + id.String())
	}
	return nil
}

func (s *PegawaiServiceImpl) UploadFile(w http.ResponseWriter, r *http.Request, formValue string, pathFile string) (path string, err error) {
	if err = r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploadedFile, handler, err := r.FormFile(formValue)
	// fmt.Println("upload", uploadedFile)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	newID, _ := uuid.NewV4()
	filename := fmt.Sprintf("%s%s", "ttd_pegawai_"+newID.String(), filepath.Ext(handler.Filename))
	dir := s.Config.App.File.Dir
	DokumenDir := s.Config.App.File.TtdPegawai

	if pathFile == "" {
		path = filepath.Join(DokumenDir, filename)
	} else {
		path = pathFile
	}
	fileLocation := filepath.Join(dir, path)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("ERROR FILE:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err = io.Copy(targetFile, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("ERROR COPY FILE:", err)
		return
	}
	return
}
func (s *PegawaiServiceImpl) DeleteFile(path string) (err error) {
	if path == "" {
		return nil
	}

	dir := s.Config.App.File.Dir
	fileLocation := filepath.Join(dir, path)
	err = os.Remove(fileLocation)
	return
}
