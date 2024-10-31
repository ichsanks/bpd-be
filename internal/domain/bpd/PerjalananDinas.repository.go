package bpd

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

var (
	perjalananDinasQuery = struct {
		Select                   string
		SelectDTO                string
		SelectDTO2               string
		SelectDetail             string
		SelectPenyelesaian       string
		SelectPenyelesaian2      string
		Insert                   string
		Update                   string
		SelectNoBpd              string
		UpdateFileBpd            string
		CountId                  string
		SelectIdBySppd           string
		SelectJmlReimbursement   string
		SelectCountReimbursement string
		SelectListBPDNew         string
		SelectListBPDNew2        string
		UpdateSppBpd             string
	}{
		SelectListBPDNew: `SELECT * FROM fn_list_bpd_new($1) as pd (
			id varchar, nomor varchar, nama varchar, tujuan varchar, keperluan varchar, tgl_berangkat date, tgl_kembali date, id_jenis_perjalanan_dinas varchar, id_jenis_kendaraan varchar, is_rombongan boolean, 
			status varchar, id_rule_approval varchar, created_at timestamp, created_by varchar, updated_by varchar, updated_at timestamp, is_deleted boolean, 
			id_bpd_pegawai varchar, id_pegawai_detail varchar, id_pegawai varchar, nip varchar, nama_pegawai varchar, id_jenis_approval varchar, id_pegawai_pengaju varchar, 
			file varchar, id_sppd varchar, keterangan_tujuan text, jenis_sppd varchar
		)`,
		SelectListBPDNew2: `SELECT * FROM fn_list_bpd_new(?) as pd (
			id varchar, nomor varchar, nama varchar, tujuan varchar, keperluan varchar, tgl_berangkat date, tgl_kembali date, id_jenis_perjalanan_dinas varchar, id_jenis_kendaraan varchar, is_rombongan boolean, 
			status varchar, id_rule_approval varchar, created_at timestamp, created_by varchar, updated_by varchar, updated_at timestamp, is_deleted boolean, 
			id_bpd_pegawai varchar, id_pegawai_detail varchar, id_pegawai varchar, nip varchar, nama_pegawai varchar, id_jenis_approval varchar, id_pegawai_pengaju varchar, 
			file varchar, id_sppd varchar, keterangan_tujuan text, jenis_sppd varchar
		)`,
		Select: `select id, nomor, nama, tujuan, keperluan, tgl_berangkat, tgl_kembali, id_jenis_perjalanan_dinas, id_jenis_kendaraan, is_rombongan, status, created_at, created_by, updated_at, updated_by, 
		is_deleted, id_rule_approval, id_jenis_approval, id_pegawai, id_pegawai_pengaju, file from perjalanan_dinas`,
		SelectDTO: `SELECT * FROM fn_list_bpd_selesai($1, $2, $3) as pd (
			id varchar, nomor varchar, nama varchar, tujuan varchar, keperluan varchar, tgl_berangkat date, tgl_kembali date, id_jenis_perjalanan_dinas varchar, id_jenis_kendaraan varchar, is_rombongan boolean, 
			status varchar, id_rule_approval varchar, created_at timestamp, created_by varchar, updated_by varchar, updated_at timestamp, is_deleted boolean, nama_jenis_perjalanan_dinas varchar, nama_jenis_kendaraan varchar, 
			nama_kendaraan varchar, nama_approval varchar, id_bpd_pegawai varchar, id_pegawai varchar, nip varchar, nama_pegawai varchar, id_bidang varchar, id_jenis_approval varchar, id_pegawai_pengaju varchar, id_pengajuan_bpd_histori varchar, 
			file varchar, id_sppd varchar, spp_id integer, status_bpd varchar, type_approval varchar, id_pegawai_approval varchar, esign boolean, is_max_pengajuan boolean, is_max_penyelesaian boolean, keterangan_tujuan text, nama_fungsionalitas varchar, id_jenis_tujuan varchar, jenis_sppd varchar
		) `,
		SelectDTO2: `SELECT * FROM fn_list_bpd_selesai(?, ?, ?) as pd (
			id varchar, nomor varchar, nama varchar, tujuan varchar, keperluan varchar, tgl_berangkat date, tgl_kembali date, id_jenis_perjalanan_dinas varchar, id_jenis_kendaraan varchar, is_rombongan boolean, 
			status varchar, id_rule_approval varchar, created_at timestamp, created_by varchar, updated_by varchar, updated_at timestamp, is_deleted boolean, nama_jenis_perjalanan_dinas varchar, nama_jenis_kendaraan varchar, 
			nama_kendaraan varchar, nama_approval varchar, id_bpd_pegawai varchar, id_pegawai varchar, nip varchar, nama_pegawai varchar, id_bidang varchar, id_jenis_approval varchar, id_pegawai_pengaju varchar, id_pengajuan_bpd_histori varchar, 
			file varchar, id_sppd varchar, spp_id integer, status_bpd varchar, type_approval varchar, id_pegawai_approval varchar, esign boolean, is_max_pengajuan boolean, is_max_penyelesaian boolean, keterangan_tujuan text, nama_fungsionalitas varchar, id_jenis_tujuan varchar, jenis_sppd varchar
		) `,
		SelectDetail: `SELECT * FROM fn_get_bpd($1, $2, $3) as pd (
			id varchar, nomor varchar, nama varchar, tujuan varchar, keperluan varchar, tgl_berangkat date, tgl_kembali date, id_jenis_perjalanan_dinas varchar, id_jenis_kendaraan varchar, is_rombongan boolean, 
			status varchar, id_rule_approval varchar, created_at timestamp, created_by varchar, updated_by varchar, updated_at timestamp, is_deleted boolean, nama_jenis_perjalanan_dinas varchar, nama_jenis_kendaraan varchar, 
			pilih_kendaraan boolean, nama_approval varchar, id_bpd_pegawai varchar, id_pegawai varchar, nip varchar, nama_pegawai varchar, id_jenis_approval varchar, id_pegawai_pengaju varchar, id_pengajuan_bpd_histori varchar, 
			file varchar, id_jenis_tujuan varchar, ket_tujuan text, id_sppd varchar, file_sppd text,  jenis_sppd varchar, status_bpd varchar, type_approval varchar, id_pegawai_approval varchar, esign boolean
		) `,
		SelectPenyelesaian: `SELECT * FROM fn_list_penyelesaian_bpd(?, ?, ?) as pp (
			id varchar, id_perjalanan_dinas varchar, nomor varchar, nama_bpd varchar, tujuan varchar, keperluan varchar, tgl_berangkat date, tgl_kembali date, id_jenis_perjalanan_dinas varchar, id_jenis_kendaraan varchar, is_rombongan boolean, 
			status varchar, created_at timestamp, created_by varchar, is_deleted boolean, jenis_perjalanan_dinas varchar, jenis_kendaraan varchar, id_pegawai varchar, nip varchar, nama_pegawai varchar, unor varchar, bidang varchar, jabatan varchar, kode_golongan varchar, 
			golongan varchar, id_jenis_approval varchar, id_pegawai_pengaju varchar, id_pengajuan_bpd_histori varchar, file_bpd varchar, file_penyelesaian_bpd varchar, bpd_penyesuaian varchar, tujuan_penyesuaian varchar, keperluan_penyesuaian varchar, 
			tgl_berangkat_penyesuaian date, tgl_kembali_penyesuaian date, total double precision, is_sppb boolean, is_revisi boolean, status_bpd varchar, type_approval varchar, id_pegawai_approval varchar, esign boolean
		) `,
		SelectPenyelesaian2: `SELECT * FROM fn_list_penyelesaian_bpd2($1, $2, $3) as pp (
			id varchar, id_perjalanan_dinas varchar, nomor varchar, nama_bpd varchar, tujuan varchar, keperluan varchar, tgl_berangkat date, tgl_kembali date, id_jenis_perjalanan_dinas varchar, id_jenis_kendaraan varchar, is_rombongan boolean, 
			status varchar, created_at timestamp, created_by varchar, is_deleted boolean, jenis_perjalanan_dinas varchar, jenis_kendaraan varchar, id_pegawai varchar, nip varchar, nama_pegawai varchar, unor varchar, bidang varchar, jabatan varchar, kode_golongan varchar, 
			golongan varchar, id_jenis_approval varchar, id_pegawai_pengaju varchar, id_pengajuan_bpd_histori varchar, file_bpd varchar, file_penyelesaian_bpd varchar, bpd_penyesuaian varchar, tujuan_penyesuaian varchar, keperluan_penyesuaian varchar, 
			tgl_berangkat_penyesuaian date, tgl_kembali_penyesuaian date, total double precision, is_sppb boolean, is_revisi boolean, is_um boolean, persentase_um double precision, persentase_sisa double precision, show_um boolean, show_sisa boolean, total_um double precision, sisa_um double precision, 
			status_bpd varchar, type_approval varchar, id_pegawai_approval varchar, esign boolean
		) `,
		Insert: `INSERT INTO perjalanan_dinas (id, nomor, nama, tujuan, keperluan, tgl_berangkat, tgl_kembali, id_jenis_perjalanan_dinas, id_jenis_kendaraan, is_rombongan, status, created_at, created_by, id_rule_approval, id_jenis_approval, id_pegawai, id_pegawai_pengaju, type_approval, is_pengajuan, id_sppd, tenant_id, id_branch) 
			values(:id, :nomor, :nama, :tujuan, :keperluan, :tgl_berangkat, :tgl_kembali, :id_jenis_perjalanan_dinas, :id_jenis_kendaraan, :is_rombongan, :status, :created_at, :created_by, :id_rule_approval, :id_jenis_approval, :id_pegawai, :id_pegawai_pengaju, :type_approval, :is_pengajuan, :id_sppd, :tenant_id, :id_branch) `,
		Update: `UPDATE perjalanan_dinas SET 
		        id=:id, 
				nomor=:nomor,
				nama=:nama, 
				tujuan=:tujuan,
				keperluan=:keperluan,
				tgl_berangkat=:tgl_berangkat,
				tgl_kembali=:tgl_kembali, 
				id_jenis_perjalanan_dinas=:id_jenis_perjalanan_dinas,
				id_jenis_kendaraan=:id_jenis_kendaraan, 
				is_rombongan=:is_rombongan,
				status=:status,
				id_rule_approval=:id_rule_approval,
				id_jenis_approval=:id_jenis_approval,
				id_pegawai=:id_pegawai,
				type_approval=:type_approval,
				is_pengajuan=:is_pengajuan,
				id_sppd=:id_sppd,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by, 
				file=:file, 
				is_deleted=:is_deleted `,
		SelectNoBpd: `select public.fn_get_no_bpd() `,
		UpdateFileBpd: `update perjalanan_dinas set id=:id, file=:file,
		               id_jenis_perjalanan_dinas=:id_jenis_perjalanan_dinas `,
		CountId: `select 
							count(pd.id) 
						from 
							public.perjalanan_dinas pd `,
		SelectIdBySppd: `select 
							pd.id
						from 
							public.perjalanan_dinas pd`,
		SelectJmlReimbursement: `select 
									sum(pdb.nominal) as nominal
								from 
									perjalanan_dinas_biaya pdb 
								left join 
									m_jenis_biaya mjb on mjb.id = pdb.id_jenis_biaya `,
		SelectCountReimbursement: `select 
										count(pdb.id)
									from 
										perjalanan_dinas_biaya pdb 
									left join 
										m_jenis_biaya mjb on mjb.id = pdb.id_jenis_biaya `,
		UpdateSppBpd: `update perjalanan_dinas set id=:id, spp_id=:spp_id `,
	}
)

var (
	perjalananDinasPegawaiDetailQuery = struct {
		Select                string
		SelectDTO             string
		Update                string
		InsertBulk            string
		InsertBulkPlaceholder string
		InsertPegawai         string
	}{
		Select: `select id, nomor, id_perjalanan_dinas, id_pegawai, is_pic, status, file, created_by, created_at, updated_at, updated_by, 
		is_deleted, file, tgl_berangkat, tgl_kembali, tujuan, keperluan, bpd_penyesuaian, is_revisi from perjalanan_dinas_pegawai `,
		SelectDTO: `select pdp.id, pdp.id_perjalanan_dinas, pd.id_sppd, pdp.id_pegawai, p.id_branch, p.nip, p.nama nama_pegawai, p.id_fungsionalitas, p.id_unor, p.id_level_bod, pd.tujuan, pd.keperluan, mu.id_bidang, pdp.is_pic, pdp.status, 
			mu.nama nama_unor, mb.nama nama_bidang, mj.nama nama_jabatan, mg.kode kode_golongan, mg.nama nama_golongan,
			TO_CHAR(pdp.tgl_berangkat, 'YYYY-MM-DD') tgl_berangkat, TO_CHAR(pdp.tgl_kembali, 'YYYY-MM-DD') tgl_kembali
			from perjalanan_dinas_pegawai pdp
			left join m_pegawai p on pdp.id_pegawai = p.id
			left join m_unit_organisasi_kerja mu on mu.id = p.id_unor
			left join m_bidang mb on mb.id = mu.id_bidang
			left join m_jabatan mj on mj.id = p.id_jabatan
			left join m_golongan mg on mg.id = p.id_golongan
			left join perjalanan_dinas pd on pd.id = pdp.id_perjalanan_dinas
		`,
		Update: `UPDATE perjalanan_dinas_pegawai SET
			id=:id,
			tujuan=:tujuan,
			keperluan=:keperluan,
			tgl_berangkat=:tgl_berangkat,
			tgl_kembali=:tgl_kembali,  
			is_revisi=:is_revisi,  
			file=:file,
			updated_at=:updated_at,
			updated_by=:updated_by `,
		InsertBulk:            `INSERT INTO perjalanan_dinas_pegawai(id, nomor, id_perjalanan_dinas, id_pegawai, is_pic, status, tgl_berangkat, tgl_kembali, file, created_at, created_by) values `,
		InsertBulkPlaceholder: ` (:id, :nomor, :id_perjalanan_dinas, :id_pegawai, :is_pic, :status, :tgl_berangkat, :tgl_kembali, :file, :created_at, :created_by) `,
		InsertPegawai: `INSERT INTO perjalanan_dinas_pegawai
							(id, nomor, id_perjalanan_dinas, id_pegawai, is_pic, status, tgl_berangkat, tgl_kembali, file, created_at, created_by) 
						values
							(:id, :nomor, :id_perjalanan_dinas, :id_pegawai, :is_pic, :status, :tgl_berangkat, :tgl_kembali, :file, :created_at, :created_by) `,
	}
)

var (
	biayaPegawai = struct {
		SelectBiaya           string
		InsertBulk            string
		InsertBulkPlaceholder string
	}{
		SelectBiaya: `select * from fn_get_generate_bpd_v2($1, $2) as (
			id_bpd varchar, id_sppd varchar, id_pegawai varchar, nama_pegawai varchar, tgl date, tgl_berangkat date,
			tgl_kembali date, id_komponen_biaya varchar, id_jenis_biaya varchar, biaya numeric, urut text,jenis_biaya text) `,
		InsertBulk: `INSERT INTO perjalanan_dinas_biaya
									(id, id_bpd_pegawai, id_jenis_biaya, id_komponen_biaya, nominal, 
									file, created_at, created_by, is_reimbursement, keterangan, id_pegawai) values `,
		InsertBulkPlaceholder: ` (:id, :id_bpd_pegawai, :id_jenis_biaya, :id_komponen_biaya, :nominal,
				 					:file, :created_at, :created_by, :is_reimbursement, :keterangan, :id_pegawai) `,
	}
)

var (
	pegawaiDokumenQuery = struct {
		Insert string
		Select string
	}{
		Insert: `INSERT INTO 
						perjalanan_dinas_dokumen 
							(id, id_bpd_pegawai, file, keterangan, created_at, created_by, id_syarat_dokumen) 
				values
							(:id, :id_bpd_pegawai, :file, :keterangan, :created_at, :created_by, :id_syarat_dokumen) `,
		Select: `select 
					pdd.id, pdd.id_bpd_pegawai, pdd.file, pdd.keterangan, pdd.id_syarat_dokumen,
					sd.id_dokumen, sd.is_mandatory, d.nama
				from
					perjalanan_dinas_dokumen pdd
				left join m_syarat_dokumen sd on sd.id = pdd.id_syarat_dokumen
				left join m_dokumen d on d.id = sd.id_dokumen `,
	}
)

type PerjalananDinasRepository interface {
	Create(data PerjalananDinas) error
	Update(data PerjalananDinas) error
	UpdatePerjalananDinas(data PerjalananDinas) error
	ResolveByIDDTO(id string) (data ListBpdNew, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveAllApproval(req model.StandardRequest) (data pagination.Response, err error)
	ResolveAllPenyelesaian(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id string) (data PerjalananDinas, err error)
	ResolveByDetailID(req FilterDetailBPD) (data PerjalananDinasDTO, err error)
	GetNoBpd() (nomor string, err error)
	UpdateFilePerjalananDinas(data FilesPerjalananDinas) error
	ResolveBpdPegawaiByID(id string) (data PerjalananDinasPegawaiDetail, err error)
	ResolveBpdPegawaiByIDDTO(req FilterDetailBPD) (data BpdPegawaiDTO, err error)
	UpdateBpdPegawai(data PerjalananDinasPegawaiDetail) error
	ExistByIdSppd(id string) (bool, string)
	ExistReimbursement(idBpdPegawai uuid.UUID, idPegawai uuid.UUID) (int64, string)
	GetDetailBiaya(idBpd string, idPegawai string) (data []BiayaPegawai, err error)
	UpdateSppBpd(data ResponseSpp) error
	ResolveByDetailHistori(req FilterDetailBPD) (data PerjalananDinasDTO, err error)
}

type PerjalananDinasRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvidePerjalananDinasRepositoryPostgreSQL(db *infras.PostgresqlConn) *PerjalananDinasRepositoryPostgreSQL {
	s := new(PerjalananDinasRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *PerjalananDinasRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchParams = append(searchParams, req.IdPegawai)
	searchParams = append(searchParams, req.IdPegawaiApproval)
	searchParams = append(searchParams, req.TypeApproval)
	searchRoleBuff.WriteString(" WHERE coalesce(pd.is_deleted, false) = false ")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(pd.nomor, pd.nama, pd.keperluan, pd.tujuan, pd.nama_jenis_perjalanan_dinas, pd.nama_jenis_kendaraan, pd.nama_pegawai, pd.nip) ilike ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	if req.IdPegawaiApproval != "" {
		searchRoleBuff.WriteString(" and pd.id_pengajuan_bpd_histori is not null and pd.id_pegawai <> ? ")
		searchParams = append(searchParams, req.IdPegawaiApproval)
	}

	if req.Status != "" {
		searchRoleBuff.WriteString(" and pd.status = ? ")
		searchParams = append(searchParams, req.Status)
	}

	query := r.DB.Read.Rebind("select count(*) from (" + perjalananDinasQuery.SelectDTO2 + searchRoleBuff.String() + ")s")
	// query := r.DB.Read.Rebind("select count(*) from (" + perjalananDinasQuery.SelectListBPDNew2 + searchRoleBuff.String() + ")s")

	var totalData int
	err = r.DB.Read.QueryRow(query, searchParams...).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		data.Items = make([]interface{}, 0)
		return
	}

	searchRoleBuff.WriteString(" ORDER BY " + ColumnMappPerjalananDinas[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchBpdQuery := searchRoleBuff.String()
	searchBpdQuery = r.DB.Read.Rebind(perjalananDinasQuery.SelectDTO2 + searchBpdQuery)
	rows, err := r.DB.Read.Queryx(searchBpdQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var items PerjalananDinasDTO
		// var items ListBpdNew
		err = rows.StructScan(&items)
		if err != nil {
			return
		}

		data.Items = append(data.Items, items)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) ResolveAllApproval(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchParams = append(searchParams, req.IdPegawai)
	searchParams = append(searchParams, req.IdPegawaiApproval)
	searchParams = append(searchParams, req.TypeApproval)
	searchRoleBuff.WriteString(" WHERE coalesce(pd.is_deleted, false) = false ")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(pd.nomor, pd.nama, pd.keperluan, pd.tujuan, pd.nama_jenis_perjalanan_dinas, pd.nama_jenis_kendaraan, pd.nama_pegawai, pd.nip) ilike ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	if req.IdPegawaiApproval != "" {
		searchRoleBuff.WriteString(" and pd.id_pengajuan_bpd_histori is not null ")
	}

	if req.Status != "" {
		searchRoleBuff.WriteString(" and (case when pd.status='5' then pd.status else pd.status_bpd end) = ? ")
		searchParams = append(searchParams, req.Status)
	}

	if req.StartDate != "" {
		searchRoleBuff.WriteString(" and pd.tgl_berangkat= ? ")
		searchParams = append(searchParams, req.StartDate)
	}

	if req.EndDate != "" {
		searchRoleBuff.WriteString(" and pd.tgl_kembali= ? ")
		searchParams = append(searchParams, req.EndDate)
	}

	if req.IdPegawaiApproval != "" {
		searchRoleBuff.WriteString(" and pd.id_pengajuan_bpd_histori is not null ")
	}

	if req.IdBidang != "" {
		searchRoleBuff.WriteString(" and pd.id_bidang =? ")
		searchParams = append(searchParams, req.IdBidang)

	}

	// if req.IdTransaksi != "" {
	// 	searchRoleBuff.WriteString(" and pd.jenis_sppd =? ")
	// 	searchParams = append(searchParams, req.IdTransaksi)
	// }

	query := r.DB.Read.Rebind("select count(*) from (" + perjalananDinasQuery.SelectDTO2 + searchRoleBuff.String() + ")s")

	var totalData int
	err = r.DB.Read.QueryRow(query, searchParams...).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		data.Items = make([]interface{}, 0)
		return
	}

	orders := ` ORDER BY (case when pd.status_bpd = '5' and pd.status = '5' then 4
			when pd.status_bpd = '1' then 1
			when pd.status_bpd = '3' then 2
			when pd.status_bpd = '2' then 3
			else 5 end
		),  `
	searchRoleBuff.WriteString(orders + ColumnMappPerjalananDinas[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchBpdQuery := searchRoleBuff.String()
	searchBpdQuery = r.DB.Read.Rebind(perjalananDinasQuery.SelectDTO2 + searchBpdQuery)
	rows, err := r.DB.Read.Queryx(searchBpdQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var items PerjalananDinasDTO
		err = rows.StructScan(&items)
		if err != nil {
			return
		}
		data.Items = append(data.Items, items)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) ResolveAllPenyelesaian(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchParams = append(searchParams, req.IdPegawai)
	searchParams = append(searchParams, req.IdPegawaiApproval)
	searchParams = append(searchParams, req.TypeApproval)
	searchRoleBuff.WriteString(" WHERE coalesce(pp.is_deleted, false) = false ")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(pp.nomor, pp.nama_bpd, pp.keperluan, pp.tujuan, pp.nama_pegawai, pp.nip, pp.unor, pp.bidang) ilike ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	if req.IdPegawaiApproval != "" {
		searchRoleBuff.WriteString(" and pp.id_pengajuan_bpd_histori is not null ")
	}

	if req.Status != "" {
		searchRoleBuff.WriteString(" and (case when pp.status='5' or pp.status='0' then pp.status else pp.status_bpd end) = ? ")
		searchParams = append(searchParams, req.Status)
	}
	if req.IsSppb != "" {
		searchRoleBuff.WriteString(" and pp.is_sppb= ? ")
		searchParams = append(searchParams, req.IsSppb)
	}

	query := r.DB.Read.Rebind("select count(*) from (" + perjalananDinasQuery.SelectPenyelesaian + searchRoleBuff.String() + ")s")

	var totalData int
	err = r.DB.Read.QueryRow(query, searchParams...).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		data.Items = make([]interface{}, 0)
		return
	}

	orders := ` ORDER BY (case when pp.status = '0' then 2
			when pp.status = '5' then 5
			when pp.status_bpd = '1' then 1
			when pp.status_bpd = '3' then 3
			when pp.status_bpd = '2' then 4
			else 6 end
		),  `
	searchRoleBuff.WriteString(orders + ColumnMappPenyelesaianBpd[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchBpdQuery := searchRoleBuff.String()
	searchBpdQuery = r.DB.Read.Rebind(perjalananDinasQuery.SelectPenyelesaian + searchBpdQuery)
	rows, err := r.DB.Read.Queryx(searchBpdQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var items BpdPegawaiDTO
		err = rows.StructScan(&items)
		if err != nil {
			return
		}
		data.Items = append(data.Items, items)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}

// Function digunakan untuk create with transaction
func (r *PerjalananDinasRepositoryPostgreSQL) Create(data PerjalananDinas) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		// Function create table perjalanan_dinas
		if err := r.CreateTxPDinas(tx, data); err != nil {
			e <- err
			return
		}

		//Function Insert Bulk table perjalanan_dinas_pegawai
		if err := r.txCreatePerjalananDinasPegawaiDetail(tx, data.Detail); err != nil {
			e <- err
			return
		}

		for _, d := range data.Detail {
			if err := r.txCreatePerjalananDinasDokDetail(tx, d.Dokumen); err != nil {
				e <- err
				return
			}
		}

		e <- nil

	})
}

// Function digunakan untuk update with transaction
func (r *PerjalananDinasRepositoryPostgreSQL) UpdatePerjalananDinas(data PerjalananDinas) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		// Function update table perjalanan_dinas
		if err := r.UpdateTxPDinas(tx, data); err != nil {
			e <- err
			return
		}

		// Function delete not in table perjalanan_dinas_pegawai
		ids := make([]string, 0)
		for _, d := range data.Detail {
			ids = append(ids, d.ID.String())
		}

		if err := r.txDeleteDetailNotIn(tx, data.ID.String(), ids); err != nil {
			e <- err
			return
		}

		// Function Insert Bulk table perjalanan_dinas_pegawai
		if err := r.txCreatePerjalananDinasPegawaiDetail(tx, data.Detail); err != nil {
			e <- err
			return
		}
		e <- nil

	})
}

func (r *PerjalananDinasRepositoryPostgreSQL) CreateTxPDinas(tx *sqlx.Tx, data PerjalananDinas) error {
	stmt, err := tx.PrepareNamed(perjalananDinasQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

func (r *PerjalananDinasRepositoryPostgreSQL) Update(data PerjalananDinas) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		// Function update table perjalanan_dinas
		if err := r.UpdateTxPDinas(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *PerjalananDinasRepositoryPostgreSQL) UpdateTxPDinas(tx *sqlx.Tx, data PerjalananDinas) error {
	stmt, err := tx.PrepareNamed(perjalananDinasQuery.Update + " WHERE id=:id")
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return nil
}

func (r *PerjalananDinasRepositoryPostgreSQL) txCreatePerjalananDinasPegawaiDetail(tx *sqlx.Tx, details []PerjalananDinasPegawaiDetail) (err error) {
	if len(details) == 0 {
		return
	}
	query, args, err := r.composeBulkUpsertperjalananDinasPegawaiDetailQuery(details)
	if err != nil {
		return
	}

	query = tx.Rebind(query)
	stmt, err := tx.Preparex(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Stmt.Exec(args...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) composeBulkUpsertperjalananDinasPegawaiDetailQuery(details []PerjalananDinasPegawaiDetail) (qResult string, params []interface{}, err error) {
	values := []string{}
	for _, d := range details {
		param := map[string]interface{}{
			"id":                  d.ID,
			"nomor":               d.Nomor,
			"id_perjalanan_dinas": d.IdPerjalananDinas,
			"id_pegawai":          d.IdPegawai,
			"is_pic":              d.IsPic,
			"status":              d.Status,
			"tgl_berangkat":       d.TglBerangkat,
			"tgl_kembali":         d.TglKembali,
			"file":                d.File,
			"created_at":          d.CreatedAt,
			"created_by":          d.CreatedBy,
		}
		q, args, err := sqlx.Named(perjalananDinasPegawaiDetailQuery.InsertBulkPlaceholder, param)
		if err != nil {
			return qResult, params, err
		}

		values = append(values, q)
		params = append(params, args...)

	}
	qResult = fmt.Sprintf(`%v %v 
						ON CONFLICT (id) 
						DO UPDATE SET 
						created_at=EXCLUDED.created_at,
						created_by=EXCLUDED.created_by
						 `, perjalananDinasPegawaiDetailQuery.InsertBulk, strings.Join(values, ","))
	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) txCreatePerjalananDinasDokDetail(tx *sqlx.Tx, details []PerjalananDinasBiaya) (err error) {
	if len(details) == 0 {
		return
	}
	query, args, err := r.composeBulkUpsertperjalananDinasDokDetailQuery(details)
	if err != nil {
		return
	}

	query = tx.Rebind(query)
	stmt, err := tx.Preparex(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Stmt.Exec(args...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) composeBulkUpsertperjalananDinasDokDetailQuery(details []PerjalananDinasBiaya) (qResult string, params []interface{}, err error) {
	values := []string{}
	for _, d := range details {
		param := map[string]interface{}{
			"id":                d.ID,
			"id_bpd_pegawai":    d.IDBpdPegawai,
			"id_jenis_biaya":    d.IDJenisBiaya,
			"id_komponen_biaya": d.IDKomponenBiaya,
			"keterangan":        d.Keterangan,
			"nominal":           d.Nominal,
			"id_pegawai":        d.IdPegawai,
			"file":              d.File,
			"is_reimbursement":  d.IsReimbursement,
			"created_at":        d.CreatedAt,
			"created_by":        d.CreatedBy,
		}
		q, args, err := sqlx.Named(biayaPegawai.InsertBulkPlaceholder, param)
		if err != nil {
			return qResult, params, err
		}

		values = append(values, q)
		params = append(params, args...)

	}
	qResult = fmt.Sprintf(`%v %v 
						ON CONFLICT (id) 
						DO UPDATE SET
						created_at=EXCLUDED.created_at,
						created_by=EXCLUDED.created_by
						 `, biayaPegawai.InsertBulk, strings.Join(values, ","))
	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) CreateDokumenBiaya(data PerjalananDinasBiaya) error {
	stmt, err := r.DB.Write.PrepareNamed(perjalananDinasBiayaQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

func (r *PerjalananDinasRepositoryPostgreSQL) CreatePegawaiPerjalananDinas(data PerjalananDinasPegawaiDetail) error {
	stmt, err := r.DB.Write.PrepareNamed(perjalananDinasPegawaiDetailQuery.InsertPegawai)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

func (r *PerjalananDinasRepositoryPostgreSQL) CreateDokumen(data PerjalananDinasDokumenDetail) error {
	stmt, err := r.DB.Write.PrepareNamed(pegawaiDokumenQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

func (r *PerjalananDinasRepositoryPostgreSQL) txDeleteDetailNotIn(tx *sqlx.Tx, idPDinas string, ids []string) (err error) {
	query, args, err := sqlx.In("update perjalanan_dinas_pegawai set is_deleted=true where id_perjalanan_dinas = ? AND id NOT IN (?)", idPDinas, ids)
	query = tx.Rebind(query)

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	res, err := r.DB.Write.Exec(query, args...)
	_, err = res.RowsAffected()
	if err != nil {
		return
	}
	return

}

func (r *PerjalananDinasRepositoryPostgreSQL) ResolveByID(id string) (data PerjalananDinas, err error) {
	err = r.DB.Read.Get(&data, perjalananDinasQuery.Select+" WHERE id=$1 ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) GetNoBpd() (nomor string, err error) {
	err = r.DB.Read.Get(&nomor, perjalananDinasQuery.SelectNoBpd)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) ResolveByIDDTO(id string) (data ListBpdNew, err error) {
	// err = r.DB.Read.Get(&data, perjalananDinasQuery.SelectDTO+" WHERE pd.id=$4 ", "", "", "", id)
	err = r.DB.Read.Get(&data, perjalananDinasQuery.SelectListBPDNew+" WHERE pd.id=$2 ", "", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	PDinasDetail, err := r.GetAllPerjalananDinasPegawaiDetail(id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	data.Detail = PDinasDetail
	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) ResolveByDetailID(req FilterDetailBPD) (data PerjalananDinasDTO, err error) {
	err = r.DB.Read.Get(&data, perjalananDinasQuery.SelectDetail+" WHERE pd.id=$4 ", req.IdPegawai, req.IdPegawaiApproval, req.TypeApproval, req.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	PDinasDetail, err := r.GetAllPerjalananDinasPegawaiDetail(req.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	data.Detail = PDinasDetail
	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) ResolveBpdPegawaiByID(id string) (data PerjalananDinasPegawaiDetail, err error) {
	err = r.DB.Read.Get(&data, perjalananDinasPegawaiDetailQuery.Select+" WHERE id=$1 ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) ResolveBpdPegawaiByIDDTO(req FilterDetailBPD) (data BpdPegawaiDTO, err error) {
	err = r.DB.Read.Get(&data, perjalananDinasQuery.SelectPenyelesaian2+" WHERE pp.id=$4 ", req.IdPegawai, req.IdPegawaiApproval, req.TypeApproval, req.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) GetAllPerjalananDinasPegawaiDetail(idPDinas string) (data []PerjalananDinasPegawaiDetailDTO, err error) {
	where := ` where pdp.id_perjalanan_dinas=$1 and coalesce(pdp.is_deleted,false)=false
		order by (case when pdp.is_pic=true then 1 else 2 end), pdp.created_at asc `
	rows, err := r.DB.Read.Queryx(perjalananDinasPegawaiDetailQuery.SelectDTO+where, idPDinas)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var bpd PerjalananDinasPegawaiDetailDTO
		err = rows.StructScan(&bpd)

		if err != nil {
			return
		}

		detAkom, errs := r.GetBiayaDto(bpd.ID.String(), bpd.IdPegawai, "false")
		if errs != nil {
			return
		}
		bpd.DetailAkomodasi = detAkom

		bPeg, errs := r.GetDetailBiaya(bpd.IdPerjalananDinas, bpd.IdPegawai)
		if errs != nil {
			return
		}
		bpd.Detail = bPeg

		detReim, errs := r.GetBiayaDto(bpd.ID.String(), bpd.IdPegawai, "true")
		if errs != nil {
			return
		}
		bpd.DetailReimbursement = detReim

		data = append(data, bpd)
	}
	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) UpdateFilePerjalananDinas(data FilesPerjalananDinas) error {
	stmt, err := r.DB.Write.PrepareNamed(perjalananDinasQuery.UpdateFileBpd + " WHERE id=:id ")
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return nil
}

func (r *PerjalananDinasRepositoryPostgreSQL) UpdateBpdPegawai(data PerjalananDinasPegawaiDetail) error {
	stmt, err := r.DB.Write.PrepareNamed(perjalananDinasPegawaiDetailQuery.Update + " WHERE id=:id ")
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return nil
}

func (r *PerjalananDinasRepositoryPostgreSQL) ExistByIdSppd(id string) (bool, string) {
	var exist bool
	var ids string
	err := r.DB.Read.Get(&exist, perjalananDinasQuery.CountId+" where pd.id_sppd = $1", id)
	if err != nil {
		return false, ""
	}
	errs := r.DB.Read.Get(&ids, perjalananDinasQuery.SelectIdBySppd+" where pd.id_sppd = $1", id)
	if errs != nil {
		return false, ""
	}
	return exist, ids
}

func (r *PerjalananDinasRepositoryPostgreSQL) ExistReimbursement(idBpdPegawai uuid.UUID, idPegawai uuid.UUID) (int64, string) {
	var exist int64
	var ids string
	where := ` where pdb.id_bpd_pegawai = $1 and pdb.id_pegawai = $2 `
	err := r.DB.Read.Get(&exist, perjalananDinasQuery.SelectCountReimbursement+where, idBpdPegawai, idPegawai)
	if err != nil {
		logger.ErrorWithStack(err)
		return 0, ""
	}
	errs := r.DB.Read.Get(&ids, perjalananDinasQuery.SelectJmlReimbursement+where, idBpdPegawai, idPegawai)
	if errs != nil {
		return 0, ""
	}
	return exist, ids
}

func (r *PerjalananDinasRepositoryPostgreSQL) GetDetailBiaya(idBpd string, idPegawai string) (data []BiayaPegawai, err error) {
	rows, err := r.DB.Read.Queryx(biayaPegawai.SelectBiaya, idBpd, idPegawai)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Detail Biaya NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items BiayaPegawai
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) GetDokumenByIdBpdPegawai(idBpdPegawai string) (data []PerjalananDinasDokumenDto, err error) {
	where := ` where pdd.id_bpd_pegawai=$1
		order by pdd.created_at asc `
	rows, err := r.DB.Read.Queryx(pegawaiDokumenQuery.Select+where, idBpdPegawai)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var bpd PerjalananDinasDokumenDto
		err = rows.StructScan(&bpd)

		if err != nil {
			return
		}

		data = append(data, bpd)
	}
	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) GetBiayaDto(idBpdPegawai string, idPegawai string, isReimbursement string) (data []BiayaPerjalananDinasDto, err error) {
	rows, err := r.DB.Read.Queryx(perjalananDinasBiayaQuery.SelectDto+"where pdb.id_bpd_pegawai = $1 and pdb.id_pegawai = $2 and pdb.is_reimbursement = $3 and pdb.is_deleted=false ", idBpdPegawai, idPegawai, isReimbursement)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var items BiayaPerjalananDinasDto
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}
	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) UpdateSppBpd(data ResponseSpp) error {
	stmt, err := r.DB.Write.PrepareNamed(perjalananDinasQuery.UpdateSppBpd + " WHERE id=:id ")
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return nil
}

func (r *PerjalananDinasRepositoryPostgreSQL) ResolveByDetailHistori(req FilterDetailBPD) (data PerjalananDinasDTO, err error) {
	err = r.DB.Read.Get(&data, perjalananDinasQuery.SelectDetail+" WHERE pd.id=$4 ", req.IdPegawai, req.IdPegawaiApproval, req.TypeApproval, req.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	PDinasDetail, err := r.GetAllPerjalananDinasPegawaiDetailHistori(req.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	data.Detail = PDinasDetail
	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) GetAllPerjalananDinasPegawaiDetailHistori(idPDinas string) (data []PerjalananDinasPegawaiDetailDTO, err error) {
	where := ` where pdp.id_perjalanan_dinas=$1 and coalesce(pdp.is_deleted,false)=false
		order by (case when pdp.is_pic=true then 1 else 2 end), pdp.created_at asc `
	rows, err := r.DB.Read.Queryx(perjalananDinasPegawaiDetailQuery.SelectDTO+where, idPDinas)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var bpd PerjalananDinasPegawaiDetailDTO
		err = rows.StructScan(&bpd)

		if err != nil {
			return
		}

		bPeg, errs := r.GetHistoriBiayaDinas(bpd.IdPerjalananDinas, bpd.IdPegawai)
		if errs != nil {
			return
		}
		bpd.HistoriBiaya = bPeg

		data = append(data, bpd)
	}
	return
}

func (r *PerjalananDinasRepositoryPostgreSQL) GetHistoriBiayaDinas(idPerjalananDinas string, idPegawai string) (data []HistoriPerjalananDinas, err error) {
	rows, err := r.DB.Read.Queryx(perjalananDinasBiayaQuery.SelectDtoHistoriIdPerjalanan+" where pdp.id_perjalanan_dinas = $1 and pdp.id_pegawai= $2 ORDER BY pdb.is_reimbursement ASC ", idPerjalananDinas, idPegawai)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var items HistoriPerjalananDinas
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}
	return
}
