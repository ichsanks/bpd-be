package bpd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	// "io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type PerjalananDinasService interface {
	Create(reqFormat PerjalananDinasRequest, userID string) (data PerjalananDinas, err error)
	Update(reqFormat PerjalananDinasRequest, userID string) (data PerjalananDinas, err error)
	ResolveByIDDTO(id string) (data ListBpdNew, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveAllApproval(req model.StandardRequest) (data pagination.Response, err error)
	ResolveAllPenyelesaian(req model.StandardRequest) (data pagination.Response, err error)
	DeleteByID(id string, userID string) error
	UpdateFile(id string, userID string) error
	ResolveByID(id string) (data PerjalananDinas, err error)
	ResolveByDetailID(req FilterDetailBPD) (data PerjalananDinasDTO, err error)
	ResolveBpdPegawaiByID(id string) (data PerjalananDinasPegawaiDetail, err error)
	ResolveBpdPegawaiByIDDTO(req FilterDetailBPD) (data BpdPegawaiDTO, err error)
	UpdateBpdPegawai(req BpdPegawaiRequest, userID string) (pd PerjalananDinasPegawaiDetail, err error)
	UploadFilePerjalananDinas(req FilesPerjalananDinas, userID string) (data FilesPerjalananDinas, err error)
	UploadFile(w http.ResponseWriter, r *http.Request, formValue string, pathFile string, id string) (path string, err error)
	Upload(w http.ResponseWriter, r *http.Request, formValue string) (path string, err error)
	BsreResolveSignPdf(reqFormat EsignRequest) (err error)
	ExistByIdSppd(id string) (bool, string)
	ExistReimbursement(idBpdPegawai uuid.UUID, idPegawai uuid.UUID) (int64, string)
	GetDetailBiaya(idBpd string, idPegawai string) (data []BiayaPegawai, err error)
	DeleteFile(path string) error
	DeleteDokumen(id string) error
	UpdateSppBpd(req ResponseSpp) (data ResponseSpp, err error)
	ResolveByDetailHistori(req FilterDetailBPD) (data PerjalananDinasDTO, err error)
}

type PerjalananDinasServiceImpl struct {
	PerjalananDinasRepository          PerjalananDinasRepository
	PerjalananDinasKendaraanRepository PerjalananDinasKendaraanRepository
	LogKegiatanRepository              LogKegiatanRepository
	SuratPerjalananDinasRepository     SuratPerjalananDinasRepository
	Config                             *configs.Config
}

func ProvidePerjalananDinasServiceImpl(repository PerjalananDinasRepository, bpdKendaraanRepository PerjalananDinasKendaraanRepository, logRepository LogKegiatanRepository, suratPerjalanan SuratPerjalananDinasRepository, config *configs.Config) *PerjalananDinasServiceImpl {
	s := new(PerjalananDinasServiceImpl)
	s.PerjalananDinasRepository = repository
	s.PerjalananDinasKendaraanRepository = bpdKendaraanRepository
	s.LogKegiatanRepository = logRepository
	s.SuratPerjalananDinasRepository = suratPerjalanan
	s.Config = config
	return s
}

func (s *PerjalananDinasServiceImpl) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	return s.PerjalananDinasRepository.ResolveAll(req)
}

func (s *PerjalananDinasServiceImpl) ResolveAllApproval(req model.StandardRequest) (data pagination.Response, err error) {
	return s.PerjalananDinasRepository.ResolveAllApproval(req)
}

func (s *PerjalananDinasServiceImpl) ResolveAllPenyelesaian(req model.StandardRequest) (data pagination.Response, err error) {
	return s.PerjalananDinasRepository.ResolveAllPenyelesaian(req)
}

func (s *PerjalananDinasServiceImpl) Create(reqFormat PerjalananDinasRequest, userID string) (data PerjalananDinas, err error) {
	data, _ = data.NewPerjalananDinasFormat(reqFormat, userID)
	now := time.Now()
	formattedTime := now.Format("02/01/2006")
	nomor, err := s.GetNoBpd(formattedTime)
	if err != nil {
		x := errors.New("Kendala Teknis, Silahkan Hubungi Administrasi")
		return PerjalananDinas{}, x
	}

	if nomor.Status == "gagal" {
		return PerjalananDinas{}, errors.New("Get API Nomor BPD GAGAL")
	}

	data.Nomor = nomor.NomorSurat
	err = s.PerjalananDinasRepository.Create(data)
	if err != nil {
		return PerjalananDinas{}, err
	}
	return data, nil
}

func (s *PerjalananDinasServiceImpl) Update(reqFormat PerjalananDinasRequest, userID string) (data PerjalananDinas, err error) {
	data, _ = data.NewPerjalananDinasFormat(reqFormat, userID)
	err = s.PerjalananDinasRepository.UpdatePerjalananDinas(data)
	if err != nil {
		return PerjalananDinas{}, err
	}
	return data, nil
}

func (s *PerjalananDinasServiceImpl) ResolveByIDDTO(id string) (data ListBpdNew, err error) {
	data, err = s.PerjalananDinasRepository.ResolveByIDDTO(id)
	if err != nil {
		return ListBpdNew{}, errors.New("Data Perjalanan Dinas dengan ID :" + id + " tidak ditemukan")
	}

	PDinasKendaraan, err := s.PerjalananDinasKendaraanRepository.GetAll(id)
	if err != nil {
		return
	}
	data.DetailKendaraan = PDinasKendaraan

	return
}

func (s *PerjalananDinasServiceImpl) ResolveByDetailID(req FilterDetailBPD) (data PerjalananDinasDTO, err error) {
	data, err = s.PerjalananDinasRepository.ResolveByDetailID(req)
	if err != nil {
		return PerjalananDinasDTO{}, errors.New("Data Perjalanan Dinas dengan ID :" + req.ID + " tidak ditemukan")
	}

	PDokumen, err := s.SuratPerjalananDinasRepository.GetSppdDokumenDto(*data.IdSppd)
	if err != nil {
		return
	}
	data.DetailDokumen = PDokumen

	PDinasKendaraan, err := s.PerjalananDinasKendaraanRepository.GetAll(req.ID)
	if err != nil {
		return
	}
	data.DetailKendaraan = PDinasKendaraan

	return
}

func (s *PerjalananDinasServiceImpl) ResolveByID(id string) (data PerjalananDinas, err error) {
	data, err = s.PerjalananDinasRepository.ResolveByID(id)
	if err != nil {
		return PerjalananDinas{}, errors.New("Data Perjalanan Dinas dengan ID :" + id + " tidak ditemukan")
	}

	return
}

func (s *PerjalananDinasServiceImpl) ResolveBpdPegawaiByID(id string) (data PerjalananDinasPegawaiDetail, err error) {
	data, err = s.PerjalananDinasRepository.ResolveBpdPegawaiByID(id)
	if err != nil {
		return PerjalananDinasPegawaiDetail{}, errors.New("Data Perjalanan Dinas pegawai dengan ID :" + id + " tidak ditemukan")
	}

	return
}

func (s *PerjalananDinasServiceImpl) ResolveBpdPegawaiByIDDTO(req FilterDetailBPD) (data BpdPegawaiDTO, err error) {
	data, err = s.PerjalananDinasRepository.ResolveBpdPegawaiByIDDTO(req)
	if err != nil {
		return BpdPegawaiDTO{}, errors.New("Data Perjalanan Dinas pegawai dengan ID :" + req.ID + " tidak ditemukan")
	}

	PDinasKendaraan, err := s.PerjalananDinasKendaraanRepository.GetAll(data.IdPerjalananDinas)
	if err != nil {
		return
	}
	data.DetailKendaraan = PDinasKendaraan

	logKegiatan, err := s.LogKegiatanRepository.GetAll("", req.ID)
	if err != nil {
		return
	}
	data.DetailLogKegiatan = logKegiatan

	dokumen, err := s.LogKegiatanRepository.GetAllDokumen(req.ID)
	if err != nil {
		return
	}
	data.DetailDokumen = dokumen

	return
}

func (s *PerjalananDinasServiceImpl) UpdateBpdPegawai(req BpdPegawaiRequest, userID string) (pd PerjalananDinasPegawaiDetail, err error) {
	pd, err = s.PerjalananDinasRepository.ResolveBpdPegawaiByID(req.ID)
	if err != nil {
		return PerjalananDinasPegawaiDetail{}, errors.New("Data Perjalanan Dinas Pegawai dengan ID :" + req.ID + " tidak ditemukan")
	}

	pd.FormatUpdateBpdPegawai(req, userID)
	err = s.PerjalananDinasRepository.UpdateBpdPegawai(pd)
	if err != nil {
		return PerjalananDinasPegawaiDetail{}, err
	}

	return pd, nil
}

func (s *PerjalananDinasServiceImpl) DeleteByID(id string, userID string) error {
	pd, err := s.PerjalananDinasRepository.ResolveByID(id)

	if err != nil {
		return errors.New("Data Perjalanan Dinas dengan ID :" + id + " tidak ditemukan")
	}

	now := time.Now()
	pd.IsDeleted = true
	pd.UpdatedBy = &userID
	pd.UpdatedAt = &now
	err = s.PerjalananDinasRepository.Update(pd)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Perjalanan Dinas dengan ID: " + id)
	}

	return nil
}

func (s *PerjalananDinasServiceImpl) UploadFilePerjalananDinas(req FilesPerjalananDinas, userID string) (data FilesPerjalananDinas, err error) {
	err = s.PerjalananDinasRepository.UpdateFilePerjalananDinas(req)
	if err != nil {
		return FilesPerjalananDinas{}, err
	}

	data = req
	return data, nil
}

func (s *PerjalananDinasServiceImpl) UploadFile(w http.ResponseWriter, r *http.Request, formValue string, pathFile string, id string) (path string, err error) {
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

	filename := fmt.Sprintf("%s%s", id, filepath.Ext(handler.Filename))
	dir := s.Config.App.File.Dir
	DokumenDir := s.Config.App.File.PerjalananDinas

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

func (s *PerjalananDinasServiceImpl) BsreResolveSignPdf(reqFormat EsignRequest) (err error) {
	IpEsign := s.Config.App.APIExternal.IpEsign
	pathPdf := filepath.Join(".", reqFormat.Pdf)
	pathTTD := filepath.Join(".", *reqFormat.Ttd)

	if reqFormat.Pdf == "" {
		fmt.Println("file must not be empty")
		return nil
	}
	reqFormat.Pdf = pathPdf
	reqFormat.Ttd = &pathTTD
	reqFormat.Addr = IpEsign

	err = s.doUpload(reqFormat)
	if err != nil {
		fmt.Printf("upload file [%s] error: %s", reqFormat.Pdf, err)
		fmt.Printf("upload image [%s] error: %s", reqFormat.Ttd, err)
		return err
	}
	fmt.Printf("upload file [%s] ok\n", reqFormat.Pdf)
	fmt.Printf("upload image [%s] ok\n", reqFormat.Ttd)

	return nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createReqBody(dataRequest EsignRequest) (string, io.Reader, error) {

	var err error
	pr, pw := io.Pipe()
	bw := multipart.NewWriter(pw) // body writer
	f, err := os.Open(dataRequest.Pdf)
	t, err := os.Open(*dataRequest.Ttd)

	if err != nil {
		return "", nil, err
	}
	go func() {
		defer f.Close()
		defer t.Close()
		// text part1
		// p0w, _ := bw.CreateFormField("halaman")
		// p0w.Write([]byte("PERTAMA"))

		p1w, _ := bw.CreateFormField("jenis_response")
		p1w.Write([]byte("BYTE"))

		// text part2
		// p2w, _ := bw.CreateFormField("linkQR")
		// p2w.Write([]byte("string"))

		// text part2
		p3w, _ := bw.CreateFormField("nik")
		p3w.Write([]byte(dataRequest.Nik))

		// text part2
		// p4w, _ := bw.CreateFormField("page")
		// p4w.Write([]byte("1"))

		// text part2
		p5w, _ := bw.CreateFormField("passphrase")
		p5w.Write([]byte(dataRequest.Passphrase))

		// text part2
		p6w, _ := bw.CreateFormField("tampilan")
		p6w.Write([]byte("visible"))

		p7w, _ := bw.CreateFormField("image")
		p7w.Write([]byte("true"))

		p8w, _ := bw.CreateFormField("height")
		p8w.Write([]byte("200"))

		p9w, _ := bw.CreateFormField("width")
		p9w.Write([]byte("200"))

		// p10w, _ := bw.CreateFormField("xAxis")
		// p10w.Write([]byte("300"))

		// p11w, _ := bw.CreateFormField("yAxis")
		// p11w.Write([]byte("150"))

		p12w, _ := bw.CreateFormField("tag_koordinat")
		p12w.Write([]byte("#"))

		// file part1
		_, fileName := filepath.Split(dataRequest.Pdf)
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				escapeQuotes("file"), escapeQuotes(fileName)))
		h.Set("Content-Type", "application/pdf")

		// file part1
		_, fileNameImg := filepath.Split(*dataRequest.Ttd)
		hImg := make(textproto.MIMEHeader)
		hImg.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				escapeQuotes("imageTTD"), escapeQuotes(fileNameImg)))
		hImg.Set("Content-Type", "image/png")

		fw1, _ := bw.CreatePart(h)
		var buf = make([]byte, 1024)
		cnt, _ := io.CopyBuffer(fw1, f, buf)

		fw1, _ = bw.CreatePart(hImg)
		var bufImg = make([]byte, 1024)
		cntr, _ := io.CopyBuffer(fw1, t, bufImg)

		log.Printf("copy %d bytes from file %s in total\n", cnt, fileName)
		log.Printf("copy %d bytes from file %s in total\n", cntr, fileName)
		bw.Close() //write the tail boundry
		pw.Close()
	}()
	return bw.FormDataContentType(), pr, nil
}

func (s *PerjalananDinasServiceImpl) doUpload(reqFormat EsignRequest) (err error) {
	// fmt.Println("file", reqFormat)
	usernameEsign := s.Config.App.APIExternal.UsernameEsign
	passwordEsign := s.Config.App.APIExternal.PasswordEsign
	auth := base64.StdEncoding.EncodeToString([]byte(usernameEsign + ":" + passwordEsign))
	// create body
	contType, reader, err := createReqBody(reqFormat)
	if err != nil {
		return err
	}

	log.Printf("createReqBody ok\n")
	url := fmt.Sprintf("http://%s/api/sign/pdf", reqFormat.Addr)
	req, err := http.NewRequest("POST", url, reader)

	//add headers
	req.Header.Add("Content-Type", contType)

	// req.Header.Add("Authorization", "Basic ZXNpZ246cXdlcnR5")
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Cookie", "JSESSIONID=7B116A4846189118C533EF636616D16A")

	client := &http.Client{}
	log.Printf("upload %s...\n", reqFormat.Pdf)
	log.Printf("uploadimg %s...\n", reqFormat.Ttd)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request send error:", err)
		return err
	}
	// fmt.Println("response", resp.StatusCode)
	defer resp.Body.Close()
	log.Printf("upload %s ok\n", reqFormat.Pdf)
	log.Printf("uploadimg %s ok\n", reqFormat.Ttd)

	// Save the response to a file (outFile)

	// pathSave := filepath.Join(".", fmt.Sprintf("/files/perjalanan_dinas/%s_suratPerjalananDinasSigned.pdf", time.Now().Format("20060102150405.000")))
	pathSave := filepath.Join(".", reqFormat.Pdf)
	if resp.StatusCode == 400 {
		err = os.Remove(pathSave)
		fmt.Println("Proses signing gagal : Passphrase anda salah 2031,")
		return errors.New("Proses signing gagal : Passphrase anda salah 2031")
	} else {

		outFile, err := os.Create(pathSave)
		if err != nil {
			fmt.Println("Error creating outFile:", err)
			return err
		}
		defer outFile.Close()

		// Copy the response body to outFile
		_, err = io.Copy(outFile, resp.Body)
	}
	if err != nil {
		fmt.Println("Error copying response body to surat_perjalanan_dinas:", err)
		return err
	}

	fmt.Println("Request successful. Response saved to 'surat_perjalanan_dinas'.")
	return nil
}
func (s *PerjalananDinasServiceImpl) UpdateFile(id string, userID string) error {
	pd, err := s.PerjalananDinasRepository.ResolveByID(id)

	if err != nil {
		return errors.New("Data Perjalanan Dinas dengan ID :" + id + " tidak ditemukan")
	}

	now := time.Now()
	*pd.File = ""
	pd.UpdatedBy = &userID
	pd.UpdatedAt = &now
	err = s.PerjalananDinasRepository.Update(pd)
	if err != nil {
		return errors.New("Ada kesalahan dalam mengupdate data Perjalanan Dinas dengan ID: " + id)
	}

	return nil
}

func (s *PerjalananDinasServiceImpl) Upload(w http.ResponseWriter, r *http.Request, path_file string) (path string, err error) {
	if err = r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploadedFile, handler, err := r.FormFile("file")
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
	filename := fmt.Sprintf("%s%s", "perjalanan_dinas"+newID.String(), filepath.Ext(handler.Filename))
	dir := s.Config.App.File.Dir
	ImageUrlDir := s.Config.App.File.PerjalananDinas

	if path_file == "" {
		path = filepath.Join(ImageUrlDir, filename)
	} else {
		path = path_file
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

func (s *PerjalananDinasServiceImpl) ExistByIdSppd(id string) (bool, string) {
	return s.PerjalananDinasRepository.ExistByIdSppd(id)
}

func (s *PerjalananDinasServiceImpl) ExistReimbursement(idBpdPegawai uuid.UUID, idPegawai uuid.UUID) (int64, string) {
	return s.PerjalananDinasRepository.ExistReimbursement(idBpdPegawai, idPegawai)
}

func (s *PerjalananDinasServiceImpl) GetDetailBiaya(idBpd string, idPegawai string) (data []BiayaPegawai, err error) {
	return s.PerjalananDinasRepository.GetDetailBiaya(idBpd, idPegawai)
}

func (s *PerjalananDinasServiceImpl) DeleteDokumen(id string) error {
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

func (s *PerjalananDinasServiceImpl) DeleteFile(path string) (err error) {
	dir := s.Config.App.File.Dir
	DokumenDir := path
	fileLocation := filepath.Join(dir, DokumenDir)
	err = os.Remove(fileLocation)
	return
}

func (s *PerjalananDinasServiceImpl) GetNoBpd(tanggal string) (data ResponseDataNomor, err error) {
	url := s.Config.App.APIExternal.IpNadineBpd
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("tanggal", tanggal)
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	responseData, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Response Data (raw):", string(responseData))

	// Cek apakah response data valid atau kosong
	if len(responseData) == 0 {
		fmt.Println(err)
		return
	}

	// Struktur untuk memetakan response JSON
	rest := ResponseDataNomor{}
	// Unmarshal JSON dari responseData
	err = json.Unmarshal([]byte(string(responseData)), &rest)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Jika berhasil
	data = rest
	fmt.Println("Data ter-unmarshal:", data)
	// responseData, err := ioutil.ReadAll(res.Body)

	// fmt.Println("SPPD Nomor", string(responseData))

	// rest := ResponseDataNomor{}
	// json.Unmarshal([]byte(string(responseData)), &rest)
	// data = rest

	// fmt.Println("data", data)

	return
}

func (s *PerjalananDinasServiceImpl) UpdateSppBpd(req ResponseSpp) (data ResponseSpp, err error) {
	err = s.PerjalananDinasRepository.UpdateSppBpd(req)
	if err != nil {
		return ResponseSpp{}, err
	}

	data = req
	return data, nil
}

func (s *PerjalananDinasServiceImpl) ResolveByDetailHistori(req FilterDetailBPD) (data PerjalananDinasDTO, err error) {
	data, err = s.PerjalananDinasRepository.ResolveByDetailHistori(req)
	if err != nil {
		return PerjalananDinasDTO{}, errors.New("Data Perjalanan Dinas dengan ID :" + req.ID + " tidak ditemukan")
	}

	return
}
