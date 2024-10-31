package bpd

import "time"

type SapVendor struct {
	MasterRekeningId             *int    `json:"master_rekening_id"`
	MasterRekeningKodeKbb        *string `json:"master_rekening_kode_kbb"`
	MasterRekeningKodeSap        *string `json:"master_rekening_kode_sap"`
	MasterRekeningKodeKeterangan *string `json:"master_rekening_keterangan"`
}

type SapGL struct {
	MasterGlId         *int    `json:"master_gl_id"`
	MasterGlKode       *string `json:"master_gl_kode"`
	MasterGlKeterangan *string `json:"master_gl_keterangan"`
}

type SapCustomer struct {
	MasterCustomerId      *int    `json:"master_customer_id"`
	MasterCustomerKodeKbb *string `json:"master_customer_kode_kbb"`
	MasterCustomerKodeSap *string `json:"master_customer_kode_sap"`
	MasterCustomerNama    *string `json:"master_customer_nama"`
}

type CostCenter struct {
	MasterCostCenterId         *int    `json:"master_cost_center_id"`
	MasterCostCenterKode       *string `json:"master_cost_center_kode"`
	MasterCostCenterKeterangan *string `json:"master_cost_center_keterangan"`
	MasterCostBudget           *string `json:"master_cost_budget"`
}

type ProfitCenter struct {
	MasterProfitCenterId   *int    `json:"master_profit_center_id"`
	MasterProfitUnit       *string `json:"master_profit_unit"`
	MasterCenterProfitKode *string `json:"master_profit_center_kode"`
}
type CashFlow struct {
	MasterCashFlowId         *int    `json:"master_cash_flow_id"`
	MasterCashFlowKode       *string `json:"master_cash_flow_kode"`
	MasterCashFlowKey        *string `json:"master_cash_flow_key"`
	MasterCashFlowKeterangan *string `json:"master_cash_flow_keterangan"`
}
type SumberDana struct {
	SumberDanaId   *int    `json:"sumber_dana_id"`
	NamaSumberDana *string `json:"nama_sumber_dana"`
}
type Bagian struct {
	MasterBagianId           *int    `json:"master_bagian_id"`
	MasterBagianKode         *string `json:"master_bagian_kode"`
	MasterBagianNama         *string `json:"master_bagian_nama"`
	MasterBagianKepalaBagian *string `json:"master_bagian_kepala_bagian"`
	MasterBagianKeterangan   *string `json:"master_bagian_keterangan"`
}
type InoKaryawan struct {
	KaryawanId         *string `json:"karyawan_id"`
	KaryawanNik        *string `json:"karyawan_nik"`
	KaryawanNama       *string `json:"karyawan_nama"`
	KaryawanNoVendor   *string `json:"karyawan_no_vendor"`
	KaryawanNoRekening *string `json:"karyawan_no_rekening"`
	KaryawanNamaBank   *string `json:"karyawan_nama_bank"`
}

type GetSpp struct {
	SppId              *int    `json:"spp_id"`
	MasterHakAksesNama *string `json:"master_hak_akses_nama"`
	SppbNo             *string `json:"sppb_no"`
	SppnNo             *string `json:"sppn_no"`
	StatusBayar        *string `json:"status_bayar"`
	StatusTerima       *string `json:"status_terima"`
}

type GetNomor struct {
	Data int `json:"data"`
}

type GetRekamJejak struct {
	Data string `json:"data"`
}

type GetDetailSpp struct {
	Data string `json:"data"`
}

type ResponseDataVendor struct {
	Data []SapVendor `json:"data"`
}

type ResponseDataGL struct {
	Data []SapGL `json:"data"`
}

type ResponseDataCustomer struct {
	Data []SapCustomer `json:"data"`
}

type ResponseDataCostCenter struct {
	Data []CostCenter `json:"data"`
}

type ResponseDataProfitCenter struct {
	Data []ProfitCenter `json:"data"`
}

type ResponseDataCashFlow struct {
	Data []CashFlow `json:"data"`
}
type ResponseDataSumberDana struct {
	Data []SumberDana `json:"data"`
}
type ResponseDataBagian struct {
	Data []Bagian `json:"data"`
}

type ResponseDataInoKaryawan struct {
	Data []InoKaryawan `json:"data"`
}

type ResponseDataSppId struct {
	Data []GetSpp `json:"data"`
}

// Sppb

type ResponseDataSppb struct {
	Data struct {
		MasterUserId         *string `json:"master_user_id"`
		MasterBagianId       *string `json:"master_bagian_id"`
		SppbJenis            *string `json:"sppb_jenis"`
		SppbNo               *string `json:"sppb_no"`
		SppbUrutan           *string `json:"sppb_urutan"`
		SppbBulan            *string `json:"sppb_bulan"`
		SppbTahun            *string `json:"sppb_tahun"`
		SppbReferensi        *string `json:"sppb_referensi"`
		SppbAu53             *string `json:"sppb_au_53"`
		SppbBeritaAcara      *string `json:"sppb_berita_acara"`
		SppbSpOpl            *string `json:"sppb_sp_opl"`
		SppbFakturPajak      *string `json:"sppb_faktur_pajak"`
		SppbTanggal          *string `json:"sppb_tanggal"`
		SppbMetodePembayaran *string `json:"sppb_metode_pembayaran"`
		SppbDateMetpen       *string `json:"sppb_data_metpen"`
		SppbCatatan          *string `json:"sppb_catatan"`
		SppbTotal            *string `json:"sppb_total"`
		UpdatedAt            *string `json:"updated_at"`
		CreatedAt            *string `json:"created_at"`
		SppbId               *int    `json:"sppb_id"`
	} `json:"data"`
}
type DataSppb struct {
	MasterUserId         *string `json:"master_user_id"`
	MasterBagianId       *string `json:"master_bagian_id"`
	SppbJenis            *string `json:"sppb_jenis"`
	SppbNo               *string `json:"sppb_no"`
	SppbUrutan           *string `json:"sppb_urutan"`
	SppbBulan            *string `json:"sppb_bulan"`
	SppbTahun            *string `json:"sppb_tahun"`
	SppbReferensi        *string `json:"sppb_referensi"`
	SppbAu53             *string `json:"sppb_au_53"`
	SppbBeritaAcara      *string `json:"sppb_berita_acara"`
	SppbSpOpl            *string `json:"sppb_sp_opl"`
	SppbFakturPajak      *string `json:"sppb_faktur_pajak"`
	SppbTanggal          *string `json:"sppb_tanggal"`
	SppbMetodePembayaran *string `json:"sppb_metode_pembayaran"`
	SppbDateMetpen       *string `json:"sppb_data_metpen"`
	SppbCatatan          *string `json:"sppb_catatan"`
	SppbTotal            *string `json:"sppb_total"`
	UpdatedAt            *string `json:"updated_at"`
	CreatedAt            *string `json:"created_at"`
	SppbId               *int    `json:"sppb_id"`
}

type RequestSppb struct {
	SppbJenis            *string `json:"sppb_jenis"`
	SppbNo               *string `json:"sppb_no"`
	SppbUrutan           *string `json:"sppb_urutan"`
	SppbBulan            *string `json:"sppb_bulan"`
	SppbTahun            *string `json:"sppb_tahun"`
	SppbKwitansi         *string `json:"sppb_kwitansi"`
	SppbReferensi        *string `json:"sppb_referensi"`
	SppbAu53             *string `json:"sppb_au_53"`
	SppbBeritaAcara      *string `json:"sppb_berita_acara"`
	SppbSpOpl            *string `json:"sppb_sp_opl"`
	SppbFakturPajak      *string `json:"sppb_faktur_pajak"`
	SppbTanggal          *string `json:"sppb_tanggal"`
	SppbMetodePembayaran *string `json:"sppb_metode_pembayaran"`
	SppbDateMetpen       *string `json:"sppb_data_metpen"`
	SppbCatatan          *string `json:"sppb_catatan"`
	SppbTotal            *string `json:"sppb_total"`
	// req SPP
	SppbId             *int    `json:"sppb_id"`
	MasterBagianId     *string `json:"master_bagian_id"`
	SppTanggal         *string `json:"spp_tanggal"`
	SppJenisSumberdana *string `json:"spp_jenis_sumber_dana"`
	SppStatusPosisi    *string `json:"spp_status_posisi"`
	SppAlur            *string `json:"spp_alur"`
	SppJenis           *string `json:"spp_jenis"`
	SppApkBpd          *string `json:"spp_apk_bpd"`
	// req SPPB ISI dan kariyawan
	RequestSppbIsiDanUraianKariyawan []RequestSppbIsiDanUraianKariyawan `json:"request_sppb_isi_uraian_kariyawan"`
	// req rekamjejak
	SppId           *int    `json:"spp_id"`
	MasterUserId    *string `json:"master_user_id"`
	RekamJejakWaktu string  `json:"rekam_jejak_waktu"`
}

type RequestSppbIsiDanUraianKariyawan struct {
	IdBpdPegawai       string  `json:"idBpdPegawai"`
	SppbId             *int    `json:"sppb_id"`
	MasterKodeVendorId *string `json:"master_kode_vendor_id"`
	MasterGlId         *string `json:"master_gl_id"`
	MasterCustomerId   *string `json:"master_customer_id"`
	MasterCostCenterId *string `json:"master_cost_center_id"`
	MasterCostProfitId *string `json:"master_cost_profit_id"`
	MasterCashFlowId   *string `json:"master_cash_flow_id"`
	// req uraian
	SppbIsiId         *int    `json:"sppb_isi_id"`
	SppbUraianUraian  *string `json:"sppb_uraian_uraian"`
	SppbUraianNominal *string `json:"sppb_uraian_nominal"`
	// req kariyawan
	KaryawanNama     *string `json:"karyawan_nama"`
	KaryawanNamaBank *string `json:"karyawan_nama_bank"`
	KaryawanNoRek    *string `json:"karyawan_no_rek"`
	KaryawanAlamat   *string `json:"karyawan_alamat"`
}

type ResponseDataSpp struct {
	Data struct {
		SppbId             *string `json:"sppb_id"`
		MasterBagianId     *string `json:"master_bagian_id"`
		SppTanggal         *string `json:"sppb_tanggal"`
		SppJenisSumberdana *string `json:"spp_jenis_sumber_dana"`
		SppStatusPosisi    *string `json:"spp_status_posisi"`
		SppAlur            *string `json:"spp_alur"`
		SppJenis           *string `json:"spp_jenis"`
		UpdatedAt          *string `json:"updated_at"`
		CreatedAt          *string `json:"created_at"`
		SppId              *int    `json:"spp_id"`
	} `json:"data"`
}

type DataSpp struct {
	SppbId             *string `json:"sppb_id"`
	MasterBagianId     *string `json:"master_bagian_id"`
	SppTanggal         *string `json:"sppb_tanggal"`
	SppJenisSumberdana *string `json:"spp_jenis_sumber_dana"`
	SppStatusPosisi    *string `json:"spp_status_posisi"`
	SppAlur            *string `json:"spp_alur"`
	SppJenis           *string `json:"spp_jenis"`
	UpdatedAt          *string `json:"updated_at"`
	CreatedAt          *string `json:"created_at"`
	SppId              *int    `json:"spp_id"`
}

type ResponseDataSppbIsi struct {
	Data struct {
		SppbId             *string `json:"sppb_id"`
		MasterKodeVendorId *string `json:"master_kode_vendor_id"`
		MasterGlId         *string `json:"master_gl_id"`
		MasterCustomerId   *string `json:"master_customer_id"`
		MasterCostCenterId *string `json:"master_cost_center_id"`
		MasterCostProfitId *string `json:"master_cost_profit_id"`
		MasterCashFlowId   *string `json:"master_cash_flow_id"`
		UpdatedAt          *string `json:"updated_at"`
		CreatedAt          *string `json:"created_at"`
		SppIsiId           *int    `json:"sppb_isi_id"`
	} `json:"data"`
}

type DataSppbIsi struct {
	SppbId             *string `json:"sppb_id"`
	MasterKodeVendorId *string `json:"master_kode_vendor_id"`
	MasterGlId         *string `json:"master_gl_id"`
	MasterCustomerId   *string `json:"master_customer_id"`
	MasterCostCenterId *string `json:"master_cost_center_id"`
	MasterCostProfitId *string `json:"master_cost_profit_id"`
	MasterCashFlowId   *string `json:"master_cash_flow_id"`
	UpdatedAt          *string `json:"updated_at"`
	CreatedAt          *string `json:"created_at"`
	SppIsiId           *int    `json:"sppb_isi_id"`
}
type ResponseDataSppbUraian struct {
	Data struct {
		SppbIsiId         *string `json:"sppb_isi_id"`
		SppbUraianUraian  *string `json:"sppb_uraian_uraian"`
		SppbUraianNominal *string `json:"sppb_uraian_nominal"`
		UpdatedAt         *string `json:"updated_at"`
		CreatedAt         *string `json:"created_at"`
		SppbUraianId      *int    `json:"sppb_uraian_id"`
	} `json:"data"`
}
type DataSppbUraian struct {
	SppbIsiId         *string `json:"sppb_isi_id"`
	SppbUraianUraian  *string `json:"sppb_uraian_uraian"`
	SppbUraianNominal *string `json:"sppb_uraian_nominal"`
	UpdatedAt         *string `json:"updated_at"`
	CreatedAt         *string `json:"created_at"`
	SppbUraianId      *int    `json:"sppb_uraian_id"`
}

type ResponseDataRekeningKariyawan struct {
	Data struct {
		SppbId            *string `json:"sppb_id"`
		KariyawanNama     *string `json:"karyawan_nama"`
		KariyawanNamaBank *string `json:"karyawan_nama_bank"`
		KariyawanNoRek    *string `json:"karyawan_no_rek"`
		KariyawanAlamat   *string `json:"karyawan_alamat"`
		UpdatedAt         *string `json:"updated_at"`
		CreatedAt         *string `json:"created_at"`
		KaryawanId        *int    `json:"karyawan_id"`
	} `json:"data"`
}

type DataRekeningKariyawan struct {
	SppbId            *string `json:"sppb_id"`
	KariyawanNama     *string `json:"karyawan_nama"`
	KariyawanNamaBank *string `json:"karyawan_nama_bank"`
	KariyawanNoRek    *string `json:"karyawan_no_rek"`
	KariyawanAlamat   *string `json:"karyawan_alamat"`
	UpdatedAt         *string `json:"updated_at"`
	CreatedAt         *string `json:"created_at"`
	KaryawanId        *int    `json:"karyawan_id"`
}

type ResponseDataRekamJejak struct {
	Data struct {
		SppId           *string `json:"spp_id"`
		MasterUserId    *string `json:"master_user_id"`
		RekamJejakWaktu *string `json:"rekam_jejak_waktu"`
		UpdatedAt       *string `json:"updated_at"`
		CreatedAt       *string `json:"created_at"`
		RekamJejakId    *int    `json:"rekam_jejak_id"`
	} `json:"data"`
}

type DataRekamJejak struct {
	SppId           *string `json:"spp_id"`
	MasterUserId    *string `json:"master_user_id"`
	RekamJejakWaktu *string `json:"rekam_jejak_waktu"`
	UpdatedAt       *string `json:"updated_at"`
	CreatedAt       *string `json:"created_at"`
	RekamJejakId    *int    `json:"rekam_jejak_id"`
}

type ResponseListDataSpp struct {
	Data []ListDataSpp `json:"data"`
}

type ListDataSpp struct {
	SppId                  *int     `json:"spp_id"`
	SppbId                 *int     `json:"sppb_id"`
	SppnId                 *int     `json:"sppn_id"`
	SppStatusProses        *int     `json:"spp_status_proses"`
	Tanggal                *string  `json:"tanggal"`
	SppbNo                 *string  `json:"sppb_no"`
	SppbTanggal            *string  `json:"sppb_tanggal"`
	SppbTotal              *float64 `json:"sppb_total"`
	SppnNo                 *string  `json:"sppn_no"`
	SppnTanggal            *string  `json:"sppn_tanggal"`
	SppnJumlah             *float64 `json:"sppn_jumlah"`
	SppStatusOb            *int     `json:"spp_status_ob"`
	SppbUraian2            *string  `json:"sppb_uraian2"`
	SppnUraian2            *string  `json:"sppn_uraian2"`
	StatusSppKeterangan    *string  `json:"status_spp_keterangan"`
	StatusProsesKeterangan *string  `json:"status_proses_keterangan"`
}

type StatusUpdateBpdSuperman struct {
	ID     string `db:"id" json:"id"`
	IsSppb bool   `db:"is_sppb" json:"isSppb"`
}

type ResponseDataSppNew struct {
	MasterBagianID     int       `json:"master_bagian_id"`
	FlowID             int       `json:"flow_id"`
	CompanyID          int       `json:"company_id"`
	SppTanggal         string    `json:"spp_tanggal"`
	SppJenisSumberDana int       `json:"spp_jenis_sumber_dana"`
	SppdProses         int       `json:"sppd_proses"`
	SppdPosisi         int       `json:"sppd_posisi"`
	SppBuat            int       `json:"spp_buat"`
	SppApkBpd          int       `json:"spp_apk_bpd"`
	SppbID             int       `json:"sppb_id"`
	UpdatedAt          time.Time `json:"updated_at"`
	CreatedAt          time.Time `json:"created_at"`
	SppID              int       `json:"spp_id"`
}

// Struct untuk dataSppb
type ResponseDataSppbNew struct {
	MasterUserID         int       `json:"master_user_id"`
	MasterBagianID       int       `json:"master_bagian_id"`
	SppbJenis            string    `json:"sppb_jenis"`
	SppbNo               string    `json:"sppb_no"`
	SppbUrutan           string    `json:"sppb_urutan"`
	SppbBulan            string    `json:"sppb_bulan"`
	SppbTahun            string    `json:"sppb_tahun"`
	SppbKwitansi         string    `json:"sppb_kwitansi"`
	SppbReferensi        string    `json:"sppb_referensi"`
	SppbSpOpl            string    `json:"sppb_sp_opl"`
	SppbTanggal          string    `json:"sppb_tanggal"`
	SppbMetodePembayaran string    `json:"sppb_metode_pembayaran"`
	SppbDataMetpen       string    `json:"sppb_data_metpen"`
	SppbTotal            string    `json:"sppb_total"`
	UpdatedAt            time.Time `json:"updated_at"`
	CreatedAt            time.Time `json:"created_at"`
	SppbID               int       `json:"sppb_id"`
}

// Struct utama untuk menggabungkan dataSpp dan dataSppb
type ResponseDataSppbdanSpp struct {
	DataSpp  ResponseDataSppNew  `json:"dataSpp"`
	DataSppb ResponseDataSppbNew `json:"dataSppb"`
}

type RequestSppbdanSpp struct {
	MasterUserID         int    `json:"master_user_id"`
	MasterBagianID       int    `json:"master_bagian_id"`
	SppbJenis            string `json:"sppb_jenis"`
	SppbNo               string `json:"sppb_no"`
	SppbUrutan           string `json:"sppb_urutan"`
	SppbBulan            string `json:"sppb_bulan"`
	SppbTahun            string `json:"sppb_tahun"`
	SppbKwitansi         string `json:"sppb_kwitansi"`
	SppbReferensi        string `json:"sppb_referensi"`
	SppbSpOpl            string `json:"sppb_sp_opl"`
	SppbTanggal          string `json:"sppb_tanggal"`
	SppbMetodePembayaran string `json:"sppb_metode_pembayaran"`
	SppbDataMetpen       string `json:"sppb_data_metpen"`
	SppbCatatan          string `json:"sppb_catatan"`
	SppbTotal            string `json:"sppb_total"`
	FlowID               int    `json:"flow_id"`
	CompanyID            int    `json:"company_id"`
	SppTanggal           string `json:"spp_tanggal"`
	SppJenisSumberDana   int    `json:"spp_jenis_sumber_dana"`
}

// New Superman Insert Bulk
// Define struct for "sppb_uraian" fields
type SPPBUraian struct {
	SPPBUraianUraian  string  `json:"sppb_uraian_uraian"`
	SPPBUraianNominal float64 `json:"sppb_uraian_nominal"`
	SPPBUraianTotal   float64 `json:"sppb_uraian_total"`
}

// Define struct for "isi_sppb" fields
type ISISPPB struct {
	MasterKodeVendorID   string       `json:"master_kode_vendor_id"`
	MasterGLID           string       `json:"master_gl_id"`
	MasterCustomerID     string       `json:"master_customer_id"`
	MasterCostCenterID   string       `json:"master_cost_center_id"`
	MasterProfitCenterID string       `json:"master_profit_center_id"`
	MasterCashFlowID     string       `json:"master_cash_flow_id"`
	SPPBUraian           []SPPBUraian `json:"sppb_uraian"`
}

type RequestPayloadSuperman struct {
	MasterUserID         int       `json:"master_user_id"`
	MasterBagianID       int       `json:"master_bagian_id"`
	SPPBJenis            string    `json:"sppb_jenis"`
	SPPBKwitansi         string    `json:"sppb_kwitansi"`
	SPPBReferensi        string    `json:"sppb_referensi"`
	SPPBSpOPL            string    `json:"sppb_sp_opl"`
	SPPBMetodePembayaran string    `json:"sppb_metode_pembayaran"`
	SPPBDataMetpen       string    `json:"sppb_data_metpen"`
	SPPBCatatan          string    `json:"sppb_catatan"`
	SPPBTotal            string    `json:"sppb_total"`
	FlowID               int       `json:"flow_id"`
	CompanyID            int       `json:"company_id"`
	SPPJenisSumberDana   int       `json:"spp_jenis_sumber_dana"`
	KaryawanNama         string    `json:"karyawan_nama"`
	KaryawanNamaBank     string    `json:"karyawan_nama_bank"`
	KaryawanNoRek        string    `json:"karyawan_no_rek"`
	KaryawanAlamat       string    `json:"karyawan_alamat"`
	ISISPPB              []ISISPPB `json:"isi_sppb"`
}

// Struct untuk dataSpp
// Custom time type to handle non-standard time format
type CustomTime time.Time

// UnmarshalJSON method to parse the time in the format "2006-01-02 15:04:05"
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Trim the quotes from the JSON string
	s := string(b[1 : len(b)-1])
	// Parse the string into a time.Time object using the custom format
	parsedTime, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	// Assign the parsed time to the receiver
	*ct = CustomTime(parsedTime)
	return nil
}

// Format method to print the time as a string
func (ct CustomTime) String() string {
	return time.Time(ct).Format("2006-01-02 15:04:05")
}

type DataSppNew struct {
	MasterBagianID     int        `json:"master_bagian_id"`
	FlowID             int        `json:"flow_id"`
	CompanyID          int        `json:"company_id"`
	SppJenisSumberDana int        `json:"spp_jenis_sumber_dana"`
	SppbID             int        `json:"sppb_id"`
	SppTanggal         string     `json:"spp_tanggal"`
	UpdatedAt          CustomTime `json:"updated_at"`
	CreatedAt          CustomTime `json:"created_at"`
	SppID              int        `json:"spp_id"`
}

// Struct untuk dataSppb
type DataSppbNew struct {
	MasterUserID         int       `json:"master_user_id"`
	MasterBagianID       int       `json:"master_bagian_id"`
	SppbJenis            string    `json:"sppb_jenis"`
	SppbKwitansi         string    `json:"sppb_kwitansi"`
	SppbReferensi        string    `json:"sppb_referensi"`
	SppbSpOpl            string    `json:"sppb_sp_opl"`
	SppbMetodePembayaran string    `json:"sppb_metode_pembayaran"`
	SppbDataMetpen       string    `json:"sppb_data_metpen"`
	SppbCatatan          string    `json:"sppb_catatan"`
	SppbTotal            string    `json:"sppb_total"`
	SppbNo               string    `json:"sppb_no"`
	SppbUrutan           int       `json:"sppb_urutan"`
	SppbBulan            string    `json:"sppb_bulan"`
	SppbTahun            int       `json:"sppb_tahun"`
	SppbTanggal          string    `json:"sppb_tanggal"`
	UpdatedAt            time.Time `json:"updated_at"`
	CreatedAt            time.Time `json:"created_at"`
	SppbID               int       `json:"sppb_id"`
}

// Struct untuk dataRekamJejak
type DataRekamJejakNew struct {
	SppID              int       `json:"spp_id"`
	MasterUserIDAsal   int       `json:"master_user_id_asal"`
	MasterUserID       int       `json:"master_user_id"`
	MasterUserIDTujuan int       `json:"master_user_id_tujuan"`
	RekamJejakStatus   int       `json:"rekam_jejak_status"`
	UpdatedAt          time.Time `json:"updated_at"`
	CreatedAt          time.Time `json:"created_at"`
	RekamJejakID       int       `json:"rekam_jejak_id"`
}

// Struct untuk dataRekening
type DataRekening struct {
	KaryawanNama     string    `json:"karyawan_nama"`
	KaryawanNamaBank string    `json:"karyawan_nama_bank"`
	KaryawanNoRek    string    `json:"karyawan_no_rek"`
	KaryawanAlamat   string    `json:"karyawan_alamat"`
	SppbID           int       `json:"sppb_id"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedAt        time.Time `json:"created_at"`
	KaryawanID       int       `json:"karyawan_id"`
}

// Struct untuk sppb_uraian dalam dataIsiSppb
type SppbUraian struct {
	SppbUraianUraian  string `json:"sppb_uraian_uraian"`
	SppbUraianNominal int    `json:"sppb_uraian_nominal"`
	SppbUraianTotal   int    `json:"sppb_uraian_total"`
}

// Struct untuk dataIsiSppb
type DataIsiSppb struct {
	MasterKodeVendorID   string       `json:"master_kode_vendor_id"`
	MasterGlID           string       `json:"master_gl_id"`
	MasterCustomerID     string       `json:"master_customer_id"`
	MasterCostCenterID   string       `json:"master_cost_center_id"`
	MasterProfitCenterID string       `json:"master_profit_center_id"`
	MasterCashFlowID     string       `json:"master_cash_flow_id"`
	SppbUraian           []SppbUraian `json:"sppb_uraian"`
}

// Struct utama yang mencakup semua data
type ResponseDataSuperman struct {
	DataSpp DataSppNew `json:"dataSpp"`
	// DataSppb       DataSppbNew       `json:"dataSppb"`
	// DataRekamJejak DataRekamJejakNew `json:"dataRekamJejak"`
	// DataRekening   DataRekening      `json:"dataRekening"`
	// DataIsiSppb    []DataIsiSppb     `json:"dataIsiSppb"`
}

type RequestPayloadUraian struct {
	IdBpd        *string      `json:"idBpd"`
	SPPBKwitansi *string      `json:"sppb_kwitansi"`
	SppbTotal    *string      `json:"sppb_total"`
	SPPBUraian   []SPPBUraian `json:"sppb_uraian"`
}
