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
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type SuratPerjalananDinasService interface {
	Create(reqFormat SuratPerjalananDinasRequest, userID string) (data SuratPerjalananDinas, err error)
	Update(reqFormat SuratPerjalananDinasRequest, userID string) (data SuratPerjalananDinas, err error)
	ResolveAllDto(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveDtoByID(id string) (data SuratPerjalananDinasDto, err error)
	DeleteByID(id string, userID string) error
	ResolveByDetailID(req FilterDetailSPPD) (data SuratPerjalananDinasListDto, err error)
	ResolveAllApproval(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id string) (data SuratPerjalananDinas, err error)
	UpdateFile(id string, userID string) error
	UploadFileSuratPerjalananDinas(req FilesSuratPerjalananDinas, userID string) (data FilesSuratPerjalananDinas, err error)
	UploadFile(w http.ResponseWriter, r *http.Request, formValue string, pathFile string, id string) (path string, err error)
	BsreResolveSignPdf(reqFormat EsignRequest) (err error)
	GetDetailBiaya(tglAwal string, tglAkhir string, idSppd string, idBodLevel string) (data []DetailBiaya, err error)
	UploadLinkFileSuratPerjalananDinas(req LinkFilesSuratPerjalananDinas, userID string) (data LinkFilesSuratPerjalananDinas, err error)
}

type SuratPerjalananDinasServiceImpl struct {
	SuratPerjalananDinasRepository SuratPerjalananDinasRepository
	Config                         *configs.Config
}

func ProvideSuratPerjalananDinasServiceImpl(repository SuratPerjalananDinasRepository, config *configs.Config) *SuratPerjalananDinasServiceImpl {
	s := new(SuratPerjalananDinasServiceImpl)
	s.SuratPerjalananDinasRepository = repository
	s.Config = config
	return s
}

func (s *SuratPerjalananDinasServiceImpl) Create(reqFormat SuratPerjalananDinasRequest, userID string) (data SuratPerjalananDinas, err error) {
	t, err := time.Parse("2006-01-02", reqFormat.TglSurat)
	if err != nil {
		fmt.Println("Error parsing tanggal:", err)
		return
	}
	formattedDate := t.Format("02/01/2006")

	nomor, err := s.GetNoSppd(formattedDate)

	fmt.Println("err", err)

	if err != nil {
		x := errors.New("Kendala Teknis, Silahkan Hubungi Administrasi")
		return SuratPerjalananDinas{}, x
		// return SuratPerjalananDinas{}, errors.New("Get API Nomor SPPD Tidak Ditemukan")
	}

	if nomor.Status == "gagal" {
		return SuratPerjalananDinas{}, errors.New("Get API Nomor SPPD GAGAL")
	}

	reqFormat.NomorSurat = nomor.NomorSurat
	data, _ = data.NewSuratPerjalananDinasFormat(reqFormat, userID)
	err = s.SuratPerjalananDinasRepository.Create(data)
	if err != nil {
		return SuratPerjalananDinas{}, err
	}
	return data, nil
}

func (s *SuratPerjalananDinasServiceImpl) Update(reqFormat SuratPerjalananDinasRequest, userID string) (data SuratPerjalananDinas, err error) {
	data, _ = data.NewSuratPerjalananDinasFormat(reqFormat, userID)
	err = s.SuratPerjalananDinasRepository.UpdateSuratPerjalananDinas(data)
	if err != nil {
		return SuratPerjalananDinas{}, err
	}
	return data, nil
}

func (s *SuratPerjalananDinasServiceImpl) ResolveAllDto(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.SuratPerjalananDinasRepository.ResolveAllDto(request)
}

func (s *SuratPerjalananDinasServiceImpl) ResolveDtoByID(id string) (data SuratPerjalananDinasDto, err error) {
	return s.SuratPerjalananDinasRepository.ResolveDtoByID(id)
}

func (s *SuratPerjalananDinasServiceImpl) ResolveByID(id string) (data SuratPerjalananDinas, err error) {
	return s.SuratPerjalananDinasRepository.ResolveByID(id)
}

func (s *SuratPerjalananDinasServiceImpl) DeleteByID(id string, userID string) error {
	pd, err := s.SuratPerjalananDinasRepository.ResolveByID(id)

	if err != nil {
		return errors.New("Data Perjalanan Dinas dengan ID :" + id + " tidak ditemukan")
	}

	now := time.Now()
	pd.IsDeleted = true
	pd.UpdatedBy = &userID
	pd.UpdatedAt = &now
	err = s.SuratPerjalananDinasRepository.Update(pd)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Perjalanan Dinas dengan ID: " + id)
	}

	return nil
}

func (s *SuratPerjalananDinasServiceImpl) ResolveByDetailID(req FilterDetailSPPD) (data SuratPerjalananDinasListDto, err error) {
	data, err = s.SuratPerjalananDinasRepository.ResolveByDetailID(req)
	if err != nil {
		return SuratPerjalananDinasListDto{}, errors.New("Data Surat Perjalanan Dinas dengan ID :" + req.ID + " tidak ditemukan")
	}

	return
}
func (s *SuratPerjalananDinasServiceImpl) ResolveAllApproval(req model.StandardRequest) (data pagination.Response, err error) {
	return s.SuratPerjalananDinasRepository.ResolveAllApproval(req)
}

func (s *SuratPerjalananDinasServiceImpl) UploadFileSuratPerjalananDinas(req FilesSuratPerjalananDinas, userID string) (data FilesSuratPerjalananDinas, err error) {
	err = s.SuratPerjalananDinasRepository.UpdateFileSuratPerjalananDinas(req)
	if err != nil {
		return FilesSuratPerjalananDinas{}, err
	}

	data = req
	return data, nil
}

func (s *SuratPerjalananDinasServiceImpl) UpdateFile(id string, userID string) error {
	pd, err := s.SuratPerjalananDinasRepository.ResolveByID(id)

	if err != nil {
		return errors.New("Data Surat Perjalanan Dinas dengan ID :" + id + " tidak ditemukan")
	}

	now := time.Now()
	*pd.File = ""
	pd.UpdatedBy = &userID
	pd.UpdatedAt = &now
	err = s.SuratPerjalananDinasRepository.Update(pd)
	if err != nil {
		return errors.New("Ada kesalahan dalam mengupdate data SuratPerjalanan Dinas dengan ID: " + id)
	}

	return nil
}

func (s *SuratPerjalananDinasServiceImpl) UploadFile(w http.ResponseWriter, r *http.Request, formValue string, pathFile string, id string) (path string, err error) {
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

	// newID, _ := uuid.NewV4()
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

func (s *SuratPerjalananDinasServiceImpl) BsreResolveSignPdf(reqFormat EsignRequest) (err error) {
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

var quoteEscapers = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotess(s string) string {
	return quoteEscapers.Replace(s)
}

func createReqBodys(dataRequest EsignRequest) (string, io.Reader, error) {

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

func (s *SuratPerjalananDinasServiceImpl) doUpload(reqFormat EsignRequest) (err error) {
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

func (s *SuratPerjalananDinasServiceImpl) GetDetailBiaya(tglAwal string, tglAkhir string, idSppd string, idBodLevel string) (data []DetailBiaya, err error) {
	return s.SuratPerjalananDinasRepository.GetDetailBiaya(tglAwal, tglAkhir, idSppd, idBodLevel)
}

func (s *SuratPerjalananDinasServiceImpl) UploadLinkFileSuratPerjalananDinas(req LinkFilesSuratPerjalananDinas, userID string) (data LinkFilesSuratPerjalananDinas, err error) {
	reqPayload := LinkFilesSuratPerjalananDinas{
		ID:       req.ID,
		LinkFile: s.Config.App.URL + "/v1/files?path=" + req.LinkFile,
	}
	err = s.SuratPerjalananDinasRepository.UpdateLinkFileSuratPerjalananDinas(reqPayload)
	if err != nil {
		return LinkFilesSuratPerjalananDinas{}, err
	}

	data = req
	return data, nil
}

func (s *SuratPerjalananDinasServiceImpl) GetNoSppd(tanggal string) (data ResponseDataNomor, err error) {
	url := s.Config.App.APIExternal.IpNadineSppd
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

	// fmt.Println("SPPD Nomor", string(responseData))

	// rest := ResponseDataNomor{}
	// json.Unmarshal([]byte(string(responseData)), &rest)
	// data = rest

	// fmt.Println("data", data)

	return
}
