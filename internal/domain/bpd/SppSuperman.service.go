package bpd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
)

type SppSupermanService interface {
	GetDataVendor() (data []SapVendor, err error)
	GetDataGL() (data []SapGL, err error)
	GetDataCustomer() (data []SapCustomer, err error)
	GetDataCostCenter() (data []CostCenter, err error)
	GetDataProfitCenter() (data []ProfitCenter, err error)
	GetDataCashFlow() (data []CashFlow, err error)
	GetDataSumberDana() (data []SumberDana, err error)
	GetDataBagian() (data []Bagian, err error)
	GetDataInoKaryawan(nik string) (data []InoKaryawan, err error)
	GetDataNomorUrut() (data GetNomor, err error)
	GetDataListSpp() (data []ListDataSpp, err error)
	Create(req RequestSppb) (data DataSppb, err error)
	CreateNew(req RequestPayloadUraian) (data ResponseDataSuperman, err error)
	GetDataRekamJejak() (data GetRekamJejak, err error)
	GetDetailSpp() (data GetDetailSpp, err error)
	GetSppId(sppId int) (data []GetSpp, err error)
}

type SppSupermanServiceImpl struct {
	SppSupermanRepository SppSupermanRepository
	Config                *configs.Config
}

func ProvideSppSupermanServiceImpl(repository SppSupermanRepository, config *configs.Config) *SppSupermanServiceImpl {
	s := new(SppSupermanServiceImpl)
	s.SppSupermanRepository = repository
	s.Config = config
	return s
}

func (s *SppSupermanServiceImpl) GetDataVendor() (data []SapVendor, err error) {
	req, _ := http.NewRequest("GET", s.Config.App.APIExternal.IpSuperman+"/api/getSapVendor", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return data, errors.New("Gagal get data !")
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest ResponseDataVendor
	json.Unmarshal(responseData, &rest)

	for _, item := range rest.Data {
		var obj SapVendor = item
		data = append(data, obj)
	}

	return
}

func (s *SppSupermanServiceImpl) GetDataGL() (data []SapGL, err error) {
	req, _ := http.NewRequest("GET", s.Config.App.APIExternal.IpSuperman+"/api/getSapGL", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return data, errors.New("Gagal get data !")
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest ResponseDataGL
	json.Unmarshal(responseData, &rest)

	for _, item := range rest.Data {
		var obj SapGL = item
		data = append(data, obj)
	}

	return
}

func (s *SppSupermanServiceImpl) GetDataCustomer() (data []SapCustomer, err error) {
	req, _ := http.NewRequest("GET", s.Config.App.APIExternal.IpSuperman+"/api/getSapCustomer", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return data, errors.New("Gagal get data !")
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest ResponseDataCustomer
	json.Unmarshal(responseData, &rest)

	for _, item := range rest.Data {
		var obj SapCustomer = item
		data = append(data, obj)
	}

	return
}

func (s *SppSupermanServiceImpl) GetDataCostCenter() (data []CostCenter, err error) {
	req, _ := http.NewRequest("GET", s.Config.App.APIExternal.IpSuperman+"/api/getCostCenter", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return data, errors.New("Gagal get data !")
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest ResponseDataCostCenter
	json.Unmarshal(responseData, &rest)

	for _, item := range rest.Data {
		var obj CostCenter = item
		data = append(data, obj)
	}

	return
}

func (s *SppSupermanServiceImpl) GetDataProfitCenter() (data []ProfitCenter, err error) {

	req, _ := http.NewRequest("GET", s.Config.App.APIExternal.IpSuperman+"/api/getProfitCenter", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return data, errors.New("Gagal get data !")
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest ResponseDataProfitCenter
	json.Unmarshal(responseData, &rest)

	for _, item := range rest.Data {
		var obj ProfitCenter = item
		data = append(data, obj)
	}

	return
}

func (s *SppSupermanServiceImpl) GetDataCashFlow() (data []CashFlow, err error) {
	req, _ := http.NewRequest("GET", s.Config.App.APIExternal.IpSuperman+"/api/getCashFlow", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return data, errors.New("Gagal get data !")
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest ResponseDataCashFlow
	json.Unmarshal(responseData, &rest)

	for _, item := range rest.Data {
		var obj CashFlow = item
		data = append(data, obj)
	}

	return
}

func (s *SppSupermanServiceImpl) GetDataSumberDana() (data []SumberDana, err error) {
	req, _ := http.NewRequest("GET", s.Config.App.APIExternal.IpSuperman+"/api/getSumberDana", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return data, errors.New("Gagal get data !")
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest ResponseDataSumberDana
	json.Unmarshal(responseData, &rest)

	for _, item := range rest.Data {
		var obj SumberDana = item
		data = append(data, obj)
	}

	return
}

func (s *SppSupermanServiceImpl) GetDataBagian() (data []Bagian, err error) {
	req, _ := http.NewRequest("GET", s.Config.App.APIExternal.IpSuperman+"/api/getBagianCreateSPP", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return data, errors.New("Gagal get data !")
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest ResponseDataBagian
	json.Unmarshal(responseData, &rest)

	for _, item := range rest.Data {
		var obj Bagian = item
		data = append(data, obj)
	}

	return
}

func (s *SppSupermanServiceImpl) GetDataInoKaryawan(nik string) (data []InoKaryawan, err error) {

	req, _ := http.NewRequest("GET", s.Config.App.APIExternal.IpInoPtpn12+"/api/get_karyawan/4a685f78e08fb8037fb34905d8440be9225dcdeae25873ae0ae145d6ebd3ab3f7a80fcefb84cb2e460b2724182c2eb730b75570897d9893f48d6117582a17823T3kiHux2Py8pTxJ5fmiotFETRjSfjDFM", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return data, errors.New("Gagal get data !")
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest []InoKaryawan
	json.Unmarshal([]byte(responseData), &rest)

	// for _, reqItem := range rest {
	// 	item := InoKaryawan{}
	// 	item.KaryawanId = reqItem.KaryawanId
	// 	item.KaryawanNik = reqItem.KaryawanNik
	// 	item.KaryawanNama = reqItem.KaryawanNama
	// 	item.KaryawanNoVendor = reqItem.KaryawanNoVendor
	// 	item.KaryawanNoRekening = reqItem.KaryawanNoRekening
	// 	item.KaryawanNamaBank = reqItem.KaryawanNamaBank
	// 	data = append(data, item)
	// }

	var foundItem *InoKaryawan
	for _, item := range rest {
		if *item.KaryawanNik == nik {
			// Buat salinan objek agar tidak mengacu pada data slice
			foundItem = &InoKaryawan{
				KaryawanId:         item.KaryawanId,
				KaryawanNik:        item.KaryawanNik,
				KaryawanNama:       item.KaryawanNama,
				KaryawanNoVendor:   item.KaryawanNoVendor,
				KaryawanNoRekening: item.KaryawanNoRekening,
				KaryawanNamaBank:   item.KaryawanNamaBank,
			}
			break
		}
	}

	// Menampilkan hasil pencarian
	if foundItem != nil {
		item := InoKaryawan{}
		item.KaryawanId = foundItem.KaryawanId
		item.KaryawanNik = foundItem.KaryawanNik
		item.KaryawanNama = foundItem.KaryawanNama
		item.KaryawanNoVendor = foundItem.KaryawanNoVendor
		item.KaryawanNoRekening = foundItem.KaryawanNoRekening
		item.KaryawanNamaBank = foundItem.KaryawanNamaBank
		data = append(data, item)
		fmt.Println("Data  ditemukan.")
	} else {
		fmt.Println("Data tidak ditemukan.")
	}

	return
}

func (s *SppSupermanServiceImpl) GetDataNomorUrut() (data GetNomor, err error) {

	req, _ := http.NewRequest("GET", s.Config.App.APIExternal.IpSupermanDev+"/api/getNomorUrutSPP", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return data, errors.New("Gagal get data !")
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest GetNomor
	json.Unmarshal([]byte(string(responseData)), &rest)

	data = rest

	return
}

func (s *SppSupermanServiceImpl) GetDataNomorUrutSpp() (data GetNomor, err error) {

	url := s.Config.App.APIExternal.IpSupermanDev + "/api/getNomorUrutSPP"
	method := "POST"
	masterBagianID := 11
	// payload := strings.NewReader("master_bagian_id=11")
	payload := strings.NewReader(fmt.Sprintf("master_bagian_id=%f", masterBagianID))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	responseData, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(responseData))
	if err != nil {
		log.Fatal(err)
	}
	// Struct untuk menyimpan hasil unmarshalling
	rest := GetNomor{}
	// Unmarshal JSON ke dalam struct
	err = json.Unmarshal([]byte(responseData), &rest)
	if err != nil {
		log.Fatal(err) // Pastikan Anda menangani error dengan benar
	}
	// Menyimpan hasil ke variabel `data`
	data = rest

	return
}

func (s *SppSupermanServiceImpl) Create(request RequestSppb) (data DataSppb, err error) {
	newData1, err := s.GetDataNomorUrutSpp()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("data", newData1.Data)
	// newData2, err := s.GetDataBagian()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	kodeBagian := "11"
	Nomorurut := newData1.Data
	nomorSppb, err := s.SppSupermanRepository.GetNomorUrutSPPb(Nomorurut, kodeBagian)
	// nomorSppb, err := s.SppSupermanRepository.GetNomorUrutSPPb(Nomorurut, *kodeBagian.MasterBagianKode)
	fmt.Println("nosppb", nomorSppb)
	if err != nil {
		fmt.Println(err)
		return
	}
	// var now = time.Now()
	// bulan := fmt.Sprintf("%d\n", now.Month())
	// bulans, err := s.SppSupermanRepository.GetNomorRomawi(bulan)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// tahun := now.Year()
	// Tahun := strconv.Itoa(tahun)
	// request.SppbNo = &nomorSppb
	// request.SppbUrutan = strconv.Itoa(&Nomorurut)
	// request.SppbTahun = &Tahun
	// request.SppbBulan = &bulans
	// newData, err := s.createSPPB(request)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// data = newData

	inputReqSpp := RequestSppbdanSpp{
		// SppbId:             newData.SppbId,
		// SppTanggal:         request.SppTanggal,
		// SppJenisSumberdana: request.SppJenisSumberdana,
		// SppJenis:           request.SppJenis,
	}

	dataSppnSppb, err := s.CreateSPPBnSpp(inputReqSpp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("datasp dan sppb", dataSppnSppb)
	newDetail := RequestSppbIsiDanUraianKariyawan{
		// SppbId:             &dataSppnSppb.DataSppb.SppbID,
		// MasterKodeVendorId: d.MasterKodeVendorId,
		// MasterGlId:         d.MasterGlId,
		// MasterCustomerId:   d.MasterCustomerId,
		// MasterCostCenterId: d.MasterCostCenterId,
		// MasterCostProfitId: d.MasterCostProfitId,
		// MasterCashFlowId:   d.MasterCashFlowId,
	}

	newdataSppbIsi, errr := s.CreateSppbIsi(newDetail)
	if errr != nil {
		fmt.Println(errr)
	}
	for _, d := range request.RequestSppbIsiDanUraianKariyawan {

		newDetailUraian := RequestSppbIsiDanUraianKariyawan{
			SppbIsiId:         newdataSppbIsi.SppIsiId,
			SppbUraianUraian:  d.SppbUraianUraian,
			SppbUraianNominal: d.SppbUraianNominal,
		}
		newdataSppbUraian, errr := s.CreateSppbUraian(newDetailUraian)
		fmt.Println("uraian", newdataSppbUraian)
		if errr != nil {
			fmt.Println(errr)
		}
		// err = s.SppSupermanRepository.UpdateStatusBPDPegawai(StatusUpdateBpdSuperman{
		// 	ID:     d.IdBpdPegawai,
		// 	IsSppb: true,
		// })

	}
	newDetailRekening := RequestSppbIsiDanUraianKariyawan{
		SppbId: &dataSppnSppb.DataSppb.SppbID,
		// KaryawanNama:     d.KaryawanNama,
		// KaryawanNamaBank: d.KaryawanNamaBank,
		// KaryawanNoRek:    d.KaryawanNoRek,
		// KaryawanAlamat:   d.KaryawanAlamat,
	}
	newdatarekening, errr := s.CreateRekeningKariyawan(newDetailRekening)
	fmt.Println("rekening", newdatarekening)
	if errr != nil {
		fmt.Println(errr)
	}

	t := time.Now()
	inputReqRekamJejak := RequestSppb{
		SppId: &dataSppnSppb.DataSppb.SppbID,
		// MasterUserId:    newData.MasterUserId,
		RekamJejakWaktu: t.String(),
	}

	dataRekamJejak, err := s.CreateRekamJejak(inputReqRekamJejak)
	fmt.Println("datarekamjejak", dataRekamJejak)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func (s *SppSupermanServiceImpl) CreateSPPB(request RequestSppb) (data DataSppb, err error) {
	url := s.Config.App.APIExternal.IpSupermanDev + "/api/createSppb"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("master_user_id", "2")
	_ = writer.WriteField("master_bagian_id", *request.MasterBagianId)
	_ = writer.WriteField("sppb_jenis", "karyawan")
	_ = writer.WriteField("sppb_no", *request.SppbNo)
	_ = writer.WriteField("sppb_urutan", *request.SppbUrutan)
	_ = writer.WriteField("sppb_bulan", *request.SppbBulan)
	_ = writer.WriteField("sppb_tahun", *request.SppbTahun)
	_ = writer.WriteField("sppb_kwitansi", *request.SppbKwitansi)
	_ = writer.WriteField("sppb_referensi", *request.SppbReferensi)
	_ = writer.WriteField("sppb_au_53", *request.SppbAu53)
	_ = writer.WriteField("sppb_berita_acara", *request.SppbBeritaAcara)
	_ = writer.WriteField("sppb_sp_opl", *request.SppbSpOpl)
	_ = writer.WriteField("sppb_faktur_pajak", *request.SppbFakturPajak)
	_ = writer.WriteField("sppb_tanggal", *request.SppbTanggal)
	_ = writer.WriteField("sppb_metode_pembayaran", "bank")
	_ = writer.WriteField("sppb_data_metpen", *request.SppbDateMetpen)
	_ = writer.WriteField("sppb_catatan", *request.SppbCatatan)
	_ = writer.WriteField("sppb_total", *request.SppbTotal)
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

	fmt.Println("SPPBHASIL", string(responseData))

	rest := ResponseDataSppb{}
	json.Unmarshal([]byte(string(responseData)), &rest)
	data = rest.Data

	return
}
func (s *SppSupermanServiceImpl) CreateSpp(request RequestSppb) (data DataSpp, err error) {
	url := s.Config.App.APIExternal.IpSupermanDev + "/api/createSpp"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	sppbId := strconv.Itoa(*request.SppbId)
	_ = writer.WriteField("sppb_id", sppbId)
	_ = writer.WriteField("master_bagian_id", "1")
	_ = writer.WriteField("spp_tanggal", *request.SppTanggal)
	_ = writer.WriteField("spp_jenis_sumber_dana", *request.SppJenisSumberdana)
	_ = writer.WriteField("spp_status_posisi", "1")
	_ = writer.WriteField("spp_alur", "1")
	_ = writer.WriteField("spp_jenis", "1")
	_ = writer.WriteField("spp_apk_bpd", "1")
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

	rest := ResponseDataSpp{}
	json.Unmarshal([]byte(string(responseData)), &rest)
	data = rest.Data

	return
}

func (s *SppSupermanServiceImpl) CreateSppbIsi(request RequestSppbIsiDanUraianKariyawan) (data DataSppbIsi, err error) {
	url := s.Config.App.APIExternal.IpSupermanDev + "/api/createSppbIsi"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	sppbId := strconv.Itoa(*request.SppbId)
	_ = writer.WriteField("sppb_id", sppbId)
	_ = writer.WriteField("master_kode_vendor_id", *request.MasterKodeVendorId)
	_ = writer.WriteField("master_gl_id", *request.MasterGlId)
	_ = writer.WriteField("master_customer_id", *request.MasterCustomerId)
	_ = writer.WriteField("master_cost_center_id", *request.MasterCostCenterId)
	_ = writer.WriteField("master_profit_center_id", *request.MasterCostProfitId)
	_ = writer.WriteField("master_cash_flow_id", *request.MasterCashFlowId)
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

	rest := ResponseDataSppbIsi{}
	json.Unmarshal([]byte(string(responseData)), &rest)
	data = rest.Data

	return
}
func (s *SppSupermanServiceImpl) CreateSppbUraian(request RequestSppbIsiDanUraianKariyawan) (data DataSppbUraian, err error) {
	url := s.Config.App.APIExternal.IpSupermanDev + "/api/createUraian"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	sppbIsiId := strconv.Itoa(*request.SppbIsiId)
	_ = writer.WriteField("sppb_isi_id", sppbIsiId)
	_ = writer.WriteField("sppb_uraian_uraian", *request.SppbUraianUraian)
	_ = writer.WriteField("sppb_uraian_nominal", *request.SppbUraianNominal)
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

	rest := ResponseDataSppbUraian{}
	json.Unmarshal([]byte(string(responseData)), &rest)
	data = rest.Data

	return
}

func (s *SppSupermanServiceImpl) CreateRekeningKariyawan(request RequestSppbIsiDanUraianKariyawan) (data DataRekeningKariyawan, err error) {
	url := s.Config.App.APIExternal.IpSupermanDev + "/api/insertRekeningKaryawan"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	sppbId := strconv.Itoa(*request.SppbId)
	_ = writer.WriteField("sppb_id", sppbId)
	_ = writer.WriteField("karyawan_nama", *request.KaryawanNama)
	_ = writer.WriteField("karyawan_nama_bank", *request.KaryawanNamaBank)
	_ = writer.WriteField("karyawan_no_rek", *request.KaryawanNoRek)
	_ = writer.WriteField("karyawan_alamat", *request.KaryawanAlamat)
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

	rest := ResponseDataRekeningKariyawan{}
	json.Unmarshal([]byte(string(responseData)), &rest)
	data = rest.Data

	return
}

func (s *SppSupermanServiceImpl) CreateRekamJejak(request RequestSppb) (data DataRekamJejak, err error) {
	url := s.Config.App.APIExternal.IpSupermanDev + "/api/insertRekamJejak"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	sppId := strconv.Itoa(*request.SppId)
	_ = writer.WriteField("spp_id", sppId)
	_ = writer.WriteField("master_user_id", *request.MasterUserId)
	_ = writer.WriteField("rekam_jejak_waktu", request.RekamJejakWaktu)
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

	rest := ResponseDataRekamJejak{}
	json.Unmarshal([]byte(string(responseData)), &rest)
	data = rest.Data

	return
}

func (s *SppSupermanServiceImpl) GetDataListSpp() (data []ListDataSpp, err error) {
	req, _ := http.NewRequest("GET", s.Config.App.APIExternal.IpSupermanDev+"/api/getSPP", nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return data, errors.New("Gagal get data !")
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest ResponseListDataSpp
	json.Unmarshal([]byte(string(responseData)), &rest)

	for _, item := range rest.Data {
		var obj ListDataSpp = item
		data = append(data, obj)
	}

	return
}

func (s *SppSupermanServiceImpl) CreateSPPBnSpp(request RequestSppbdanSpp) (data ResponseDataSppbdanSpp, err error) {
	url := s.Config.App.APIExternal.IpSupermanDev + "/api/createSPPnSPPB"
	method := "POST"

	// // Variabel untuk data yang akan diisi dalam payload
	masterUserID := 239
	masterBagianID := 138
	sppbJenis := "karyawan"
	sppbNo := "BPTI3/SPPb/16/2024"
	sppbUrutan := "16"
	sppbBulan := "VIII"
	sppbTahun := "2024"
	sppbKwitansi := "YUDI WAHONO"
	sppbReferensi := "RK3F/MO/2024.08.21-4"
	sppbSpOpl := "RK3F/MO/2024.08.21-4"
	sppbTanggal := "2024-08-30"
	sppbMetodePembayaran := "bank"
	sppbDataMetpen := "input_data"
	sppbCatatan := "Terlampir di Uraian"
	sppbTotal := 1000000
	flowID := 24
	companyID := 9
	sppTanggal := "2024-01-01"
	sppJenisSumberDana := 1

	payload := fmt.Sprintf(`{
		"master_user_id": %d,
		"master_bagian_id": %d,
		"sppb_jenis": "%s",
		"sppb_no": "%s",
		"sppb_urutan": "%s",
		"sppb_bulan": "%s",
		"sppb_tahun": "%s",
		"sppb_kwitansi": "%s",
		"sppb_referensi": "%s",
		"sppb_sp_opl": "%s",
		"sppb_tanggal": "%s",
		"sppb_metode_pembayaran": "%s",
		"sppb_data_metpen": "%s",
		"sppb_catatan": "%s",
		"sppb_total": "%s",
		"flow_id": %d,
		"company_id": %d,
		"spp_tanggal": "%s",
		"spp_jenis_sumber_dana": %d
	}`,
		masterUserID, masterBagianID, sppbJenis, sppbNo, sppbUrutan, sppbBulan, sppbTahun, sppbKwitansi, sppbReferensi, sppbSpOpl, sppbTanggal, sppbMetodePembayaran, sppbDataMetpen, sppbCatatan, sppbTotal, flowID, companyID, sppTanggal, sppJenisSumberDana)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

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

	fmt.Println("SPPBdanSPPHASIL", string(responseData))

	rest := ResponseDataSppbdanSpp{}
	json.Unmarshal([]byte(string(responseData)), &rest)
	data = rest

	fmt.Println("dataspp", data.DataSpp)
	fmt.Println("datasppb", data.DataSppb)

	return
}

func (s *SppSupermanServiceImpl) CreateSppbIsiNew(request RequestSppbdanSpp) (data ResponseDataSppbdanSpp, err error) {
	url := s.Config.App.APIExternal.IpSupermanDev + "/api/createSppbIsi"
	method := "POST"

	// // Variabel untuk data yang akan diisi dalam payload
	sppbId := 1372
	masterGlID := "5785"
	masterCostCenterID := "27"

	payload := fmt.Sprintf(`{
		"sppb_id": %d,
		"master_gl_id": %s,
		"master_cost_center_id": "%s",
	}`,
		sppbId, masterGlID, masterCostCenterID)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

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

	fmt.Println("SPPBdanSPPHASIL", string(responseData))

	rest := ResponseDataSppbdanSpp{}
	json.Unmarshal([]byte(string(responseData)), &rest)
	data = rest

	fmt.Println("dataspp", data.DataSpp)
	fmt.Println("datasppb", data.DataSppb)

	return
}

func (s *SppSupermanServiceImpl) CreateUraian(request RequestSppbdanSpp) (data ResponseDataSppbdanSpp, err error) {
	url := s.Config.App.APIExternal.IpSupermanDev + "/api/createUraian"
	method := "POST"
	// // Variabel untuk data yang akan diisi dalam payload

	sppbIsiId := 3718
	sppbUraianUraian := "Deskripsi bpd"
	sppbUraianNominal := 50000
	sppbUraianTotal := 100000

	payload := fmt.Sprintf(`{
		"sppb_isi_id": %d,
		"sppb_uraian_uraian": %s,
		"sppb_uraian_nominal": "%d",
		"sppb_uraian_total": "%d",
	}`,
		sppbIsiId, sppbUraianUraian, sppbUraianNominal, sppbUraianTotal)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

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

	fmt.Println("SPPBdanSPPHASIL", string(responseData))

	rest := ResponseDataSppbdanSpp{}
	json.Unmarshal([]byte(string(responseData)), &rest)
	data = rest

	fmt.Println("dataspp", data.DataSpp)
	fmt.Println("datasppb", data.DataSppb)

	return
}

func (s *SppSupermanServiceImpl) InsertRekeningKaryawan(request RequestSppbdanSpp) (data ResponseDataSppbdanSpp, err error) {
	url := s.Config.App.APIExternal.IpSupermanDev + "/api/insertRekeningKaryawan"
	method := "POST"
	// // Variabel untuk data yang akan diisi dalam payload

	sppbId := 137
	karyawanNama := "Yudi Wahono"
	karyawanNamabank := "Bank XYZ"
	karyawanNorek := "1234567890"
	karyawanAlamat := "Jalan Raya No. 123"

	payload := fmt.Sprintf(`{
		"sppb_id": %d,
		"karyawan_nama": %s,
		"karyawan_nama_bank": %s,
		"karyawan_no_rek": %s,
		"karyawan_alamat": %s,
	}`,
		sppbId, karyawanNama, karyawanNamabank, karyawanNorek, karyawanAlamat)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

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

	fmt.Println("SPPBdanSPPHASIL", string(responseData))

	rest := ResponseDataSppbdanSpp{}
	json.Unmarshal([]byte(string(responseData)), &rest)
	data = rest

	fmt.Println("dataspp", data.DataSpp)
	fmt.Println("datasppb", data.DataSppb)

	return
}

func (s *SppSupermanServiceImpl) InsertRekamJejak(request RequestSppbdanSpp) (data ResponseDataSppbdanSpp, err error) {
	url := s.Config.App.APIExternal.IpSupermanDev + "/api/insertRekamJejak"
	method := "POST"
	// // Variabel untuk data yang akan diisi dalam payload

	sppId := 1440
	masterUserIdAsal := 193
	masterUserId := 34
	masterUserIdTujuan := 0
	rekamJejakStatus := 0
	rekamJejakWaktu := "2024-08-30 10:00:00"

	payload := fmt.Sprintf(`{
		"spp_id": %d,
		"master_user_id_asal": %d,
		"master_user_id": %d,
		"master_user_id_tujuan": %d,
		"rekam_jejak_status": %d,
		"rekam_jejak_waktu": %s,
	}`,
		sppId, masterUserIdAsal, masterUserId, masterUserIdTujuan, rekamJejakStatus, rekamJejakWaktu)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

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

	fmt.Println("SPPBdanSPPHASIL", string(responseData))

	rest := ResponseDataSppbdanSpp{}
	json.Unmarshal([]byte(string(responseData)), &rest)
	data = rest

	fmt.Println("dataspp", data.DataSpp)
	fmt.Println("datasppb", data.DataSppb)

	return
}

func (s *SppSupermanServiceImpl) CreateNew(request RequestPayloadUraian) (data ResponseDataSuperman, err error) {

	url := s.Config.App.APIExternal.IpSupermanDev + "/api/createSPPnSPPB"
	method := "POST"

	// pengajuan detail
	details := make([]SPPBUraian, 0)
	for _, d := range request.SPPBUraian {
		newDetail := SPPBUraian{
			SPPBUraianUraian:  d.SPPBUraianUraian,
			SPPBUraianNominal: d.SPPBUraianNominal,
			SPPBUraianTotal:   d.SPPBUraianTotal,
		}

		details = append(details, newDetail)
	}
	// Create payload
	payload := RequestPayloadSuperman{
		MasterUserID:         190,
		MasterBagianID:       130,
		SPPBJenis:            "karyawan",
		SPPBKwitansi:         *request.SPPBKwitansi,
		SPPBReferensi:        "RK3F/MO/2024.08.21-4",
		SPPBSpOPL:            "RK3F/MO/2024.08.21-4",
		SPPBMetodePembayaran: "bank",
		SPPBDataMetpen:       "input_data",
		SPPBCatatan:          "Terlampir di Uraian",
		SPPBTotal:            *request.SppbTotal,
		FlowID:               24,
		CompanyID:            9,
		SPPJenisSumberDana:   1,
		KaryawanNama:         "Terlampir di uraian",
		KaryawanNamaBank:     "Terlampir di uraian",
		KaryawanNoRek:        "Terlampir di uraian",
		KaryawanAlamat:       "Terlampir di uraian",
		ISISPPB: []ISISPPB{
			{
				MasterKodeVendorID:   "12",
				MasterGLID:           "12",
				MasterCustomerID:     "12",
				MasterCostCenterID:   "12",
				MasterProfitCenterID: "12",
				MasterCashFlowID:     "12",
				SPPBUraian:           details,
			},
		},
	}

	// Convert struct to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return
	}
	// Create new HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	// Read the response
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Print the response

	fmt.Println("SPP", string(responseData))
	var rest ResponseDataSuperman

	// Unmarshal the JSON response into the Go struct
	err = json.Unmarshal(responseData, &rest)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	data = rest

	// rest := ResponseDataSuperman
	// json.Unmarshal([]byte(string(responseData)), &rest)

	fmt.Println("data res", data)

	return
}

func (s *SppSupermanServiceImpl) GetDataRekamJejak() (data GetRekamJejak, err error) {

	url := s.Config.App.APIExternal.IpSupermanDev + "/spp/rekam_jejak_bpd/MzQwMTY=/cGV0dWdhcy1icGQtdW11bQ=="
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest GetRekamJejak
	json.Unmarshal([]byte(string(responseData)), &rest)

	data = rest

	return
}

func (s *SppSupermanServiceImpl) GetDetailSpp() (data GetDetailSpp, err error) {

	url := s.Config.App.APIExternal.IpSupermanDev + "/spp/detail_bpd/MzQwMTY=/cGV0dWdhcy1icGQtdW11bQ=="
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rest GetDetailSpp
	json.Unmarshal([]byte(string(responseData)), &rest)

	data = rest

	return
}

func (s *SppSupermanServiceImpl) GetSppId(sppId int) (data []GetSpp, err error) {
	url := s.Config.App.APIExternal.IpSupermanDev + "api/getStatusSPP"
	method := "POST"

	SppId := sppId
	payload := strings.NewReader(fmt.Sprintf(`{
		"spp_id":%d
	}`, SppId))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

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
	var rest ResponseDataSppId
	json.Unmarshal(responseData, &rest)

	for _, item := range rest.Data {
		var obj GetSpp = item
		data = append(data, obj)
	}

	return
}
