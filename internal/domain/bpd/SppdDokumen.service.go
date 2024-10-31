package bpd

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

type SppdDokumenService interface {
	Create(req SppdDokumenRequest, userId string) (data SppdDokumen, err error)
	Update(req SppdDokumenRequest, userId string) (data SppdDokumen, err error)
	GetAll() (data []SppdDokumen, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data SppdDokumen, err error)
	SoftDelete(id uuid.UUID, userId string) error
	UploadFileDokumen(w http.ResponseWriter, r *http.Request, formValue string, pathFile string) (path string, err error)
}

type SppdDokumenServiceImpl struct {
	SppdDokumenRepository SppdDokumenRepository
	Config                *configs.Config
}

func ProvideSppdDokumenServiceImpl(repository SppdDokumenRepository, config *configs.Config) *SppdDokumenServiceImpl {
	s := new(SppdDokumenServiceImpl)
	s.SppdDokumenRepository = repository
	s.Config = config
	return s
}

func (s *SppdDokumenServiceImpl) Create(req SppdDokumenRequest, userId string) (data SppdDokumen, err error) {
	data, _ = data.NewSppdDokumenFormat(req, userId)
	err = s.SppdDokumenRepository.Create(data)
	if err != nil {
		return SppdDokumen{}, err
	}
	return data, nil
}

func (s *SppdDokumenServiceImpl) GetAll() (data []SppdDokumen, err error) {
	return s.SppdDokumenRepository.GetAll()
}

func (s *SppdDokumenServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.SppdDokumenRepository.ResolveAll(request)
}

func (s *SppdDokumenServiceImpl) DeleteByID(id uuid.UUID) error {
	newBidang, err := s.SppdDokumenRepository.ResolveByID(id)

	if err != nil || (SppdDokumen{}) == newBidang {
		return errors.New("Data Sppd Dokumen dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.SppdDokumenRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Sppd Dokumen dengan ID: " + id.String())
	}
	return nil
}

func (s *SppdDokumenServiceImpl) Update(req SppdDokumenRequest, userId string) (data SppdDokumen, err error) {
	data, _ = data.NewSppdDokumenFormat(req, userId)
	err = s.SppdDokumenRepository.Update(data)
	if err != nil {
		return SppdDokumen{}, err
	}
	return data, nil
}

func (s *SppdDokumenServiceImpl) ResolveByID(id uuid.UUID) (data SppdDokumen, err error) {
	return s.SppdDokumenRepository.ResolveByID(id)
}

func (s *SppdDokumenServiceImpl) SoftDelete(id uuid.UUID, userId string) error {
	newBidang, err := s.SppdDokumenRepository.ResolveByID(id)

	if err != nil || (SppdDokumen{}) == newBidang {
		return errors.New("Data Sppd Dokumen dengan ID :" + id.String() + " tidak ditemukan")
	}

	newBidang.SoftDelete(userId)
	err = s.SppdDokumenRepository.Update(newBidang)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Sppd Dokumen dengan ID: " + id.String())
	}
	return nil
}

func (s *SppdDokumenServiceImpl) UploadFileDokumen(w http.ResponseWriter, r *http.Request, formValue string, pathFile string) (path string, err error) {
	if err = r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploadedFile, handler, err := r.FormFile(formValue)
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
	filename := fmt.Sprintf("%s%s", "dokumen_sppd_"+newID.String(), filepath.Ext(handler.Filename))
	dir := s.Config.App.File.Dir
	DokumenDir := s.Config.App.File.DokumenSppd

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
