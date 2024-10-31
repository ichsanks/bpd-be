package bpd

import (
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/files"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type LogKegiatanService interface {
	Create(reqFormat LogKegiatanRequest, userID string) (data LogKegiatan, err error)
	GetAll(idPerjalananDinas string, idPegawai string) (data []LogKegiatan, err error)
	UploadFile(w http.ResponseWriter, r *http.Request, formValue string, pathFile string) (path string, err error)
	DeleteByID(id string) error
	// BPD Dokumen
	GetAllDokumen(idBpdPegawai string) (data []PerjalananDinasDokumen, err error)
	CreateDokumen(reqFormat PerjalananDinasDokumenRequest, userID string) (data PerjalananDinasDokumen, err error)
	DeleteDokumen(id string) error
	UploadFileDokumen(w http.ResponseWriter, r *http.Request, formValue string, pathFile string) (path string, err error)
	GenerateTextToImage(reqFormat LogKegiatanRequest) (img image.Image, err error)
}

type LogKegiatanServiceImpl struct {
	LogKegiatanRepository LogKegiatanRepository
	Config                *configs.Config
}

func ProvideLogKegiatanServiceImpl(repository LogKegiatanRepository, config *configs.Config) *LogKegiatanServiceImpl {
	s := new(LogKegiatanServiceImpl)
	s.LogKegiatanRepository = repository
	s.Config = config
	return s
}

func (s *LogKegiatanServiceImpl) GetAll(idPerjalananDinas string, idPegawai string) (data []LogKegiatan, err error) {
	return s.LogKegiatanRepository.GetAll(idPerjalananDinas, idPegawai)
}

func (s *LogKegiatanServiceImpl) Create(reqFormat LogKegiatanRequest, userID string) (data LogKegiatan, err error) {
	data, _ = data.NewLogKegiatanFormat(reqFormat, userID)
	err = s.LogKegiatanRepository.Create(data)
	if err != nil {
		return LogKegiatan{}, err
	}

	_, err = s.GenerateTextToImage(reqFormat)
	if err != nil {
		return LogKegiatan{}, err
	}

	return data, nil
}

func (s *LogKegiatanServiceImpl) DeleteByID(id string) error {
	log, err := s.LogKegiatanRepository.ResolveByID(id)
	if err != nil || (LogKegiatan{}) == log {
		return errors.New("Data dokumen perjalanan dinas dengan ID :" + id + " tidak ditemukan")
	}

	err = s.LogKegiatanRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data log kegiatan dengan ID: " + id)
	}

	err = s.DeleteFile(log.Foto)
	if err != nil {
		return errors.New("File tidak ditemukan")
	}

	return nil
}

func (s *LogKegiatanServiceImpl) GetAllDokumen(idBpdPegawai string) (data []PerjalananDinasDokumen, err error) {
	return s.LogKegiatanRepository.GetAllDokumen(idBpdPegawai)
}

func (s *LogKegiatanServiceImpl) CreateDokumen(reqFormat PerjalananDinasDokumenRequest, userID string) (data PerjalananDinasDokumen, err error) {
	data, _ = data.NewPerjalananDinasDokumenFormat(reqFormat, userID)
	err = s.LogKegiatanRepository.CreateDokumen(data)
	if err != nil {
		return PerjalananDinasDokumen{}, err
	}
	return data, nil
}

func (s *LogKegiatanServiceImpl) DeleteDokumen(id string) error {
	dokumen, err := s.LogKegiatanRepository.ResolveByIDDokumen(id)
	if err != nil || (PerjalananDinasDokumen{}) == dokumen {
		return errors.New("Data dokumen perjalanan dinas dengan ID :" + id + " tidak ditemukan")
	}

	err = s.LogKegiatanRepository.DeleteDokumen(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data perjalanan dinas dokumen dengan ID: " + id)
	}

	err = s.DeleteFile(dokumen.File)
	if err != nil {
		return errors.New("File tidak ditemukan")
	}

	return nil
}

func (s *LogKegiatanServiceImpl) UploadFile(w http.ResponseWriter, r *http.Request, formValue string, pathFile string) (path string, err error) {
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
	filename := fmt.Sprintf("%s%s", "foto_dinas_"+newID.String(), filepath.Ext(handler.Filename))
	dir := s.Config.App.File.Dir
	DokumenDir := s.Config.App.File.FotoDinas

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

func (s *LogKegiatanServiceImpl) UploadFileDokumen(w http.ResponseWriter, r *http.Request, formValue string, pathFile string) (path string, err error) {
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
	filename := fmt.Sprintf("%s%s", "dokumen_bpd_"+newID.String(), filepath.Ext(handler.Filename))
	dir := s.Config.App.File.Dir
	DokumenDir := s.Config.App.File.DokumenBpd

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

func (s *LogKegiatanServiceImpl) DeleteFile(path string) (err error) {
	dir := s.Config.App.File.Dir
	DokumenDir := path
	fileLocation := filepath.Join(dir, DokumenDir)
	err = os.Remove(fileLocation)
	return
}

func (s *LogKegiatanServiceImpl) GenerateTextToImage(reqFormat LogKegiatanRequest) (img image.Image, err error) {
	now := time.Now()
	cTime := now.Format("02-01-2006 15:04:05")
	dir := s.Config.App.File.Dir
	fileLocation := filepath.Join(dir, reqFormat.Foto)
	latlong := model.ParseString(reqFormat.Lat) + "," + model.ParseString(reqFormat.Long)
	address := model.ParseString(reqFormat.Address)
	txts := []string{cTime, latlong, address}
	req := files.RequestTextOnImg{
		BgImgPath: fileLocation,
		FontPath:  "",
		FontSize:  17,
		Text:      txts,
	}
	img, err = files.TextOnImg(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}
