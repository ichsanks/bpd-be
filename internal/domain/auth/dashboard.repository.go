package auth

import (
	"database/sql"
	"fmt"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
)

var (
	dashboardQuery = struct {
		SelectJmlPegawai,
		SelectDataBpd,
		SelectDataSppd,
		SelectDataBpdNew,
		SelectDashboardSppd,
		SelectDashboardBpdNew,
		SelectJumlahSppd,
		SelectJumlahBpd,
		SelectDashboard string
	}{
		SelectJmlPegawai: `select count(mp.id) jml_pegawai from m_pegawai mp `,
		SelectDashboard: `select * from fn_dashboard_bpd($1, $2, $3) as (
			belum_proses_pengajuan bigint, pengajuan bigint, sedang_dinas bigint, dalam_penyelesaian bigint) `,
		SelectDataBpd: `select id, id_perjalanan_dinas, nomor, nama_bpd, tujuan, keperluan, TO_CHAR(tgl_berangkat, 'YYYY-MM-DD') tgl_berangkat, TO_CHAR(tgl_kembali, 'YYYY-MM-DD') tgl_kembali, status, status_bpd, jenis_perjalanan_dinas, jenis_kendaraan 
			from fn_list_penyelesaian_bpd($1, $2, 'PENYELESAIAN') as (
				id varchar, id_perjalanan_dinas varchar, nomor varchar, nama_bpd varchar, tujuan varchar, keperluan varchar, tgl_berangkat date, tgl_kembali date, id_jenis_perjalanan_dinas varchar, id_jenis_kendaraan varchar, is_rombongan boolean, 
				status varchar, created_at timestamp, created_by varchar, is_deleted boolean, jenis_perjalanan_dinas varchar, jenis_kendaraan varchar, id_pegawai varchar, nip varchar, nama_pegawai varchar, unor varchar, bidang varchar, jabatan varchar, kode_golongan varchar, 
				golongan varchar, id_jenis_approval varchar, id_pegawai_pengaju varchar, id_pengajuan_bpd_histori varchar, file_bpd varchar, file_penyelesaian_bpd varchar, bpd_penyesuaian varchar, tujuan_penyesuaian varchar, keperluan_penyesuaian varchar, 
				tgl_berangkat_penyesuaian date, tgl_kembali_penyesuaian date, total double precision, is_sppb boolean, is_revisi boolean, status_bpd varchar, type_approval varchar, id_pegawai_approval varchar, esign boolean
			)
			where status in ('0','1')
		`,
		SelectDataSppd: `
		SELECT * FROM fn_list_sppd($1, $2, 'PENGAJUANA_SPPD') as spd  (
			id varchar(36), tgl_surat date, nomor_surat varchar(100), id_pegawai varchar(36), jenis_tujuan varchar(36), 
			tujuan_dinas text, keperluan_dinas text, tenant_id varchar(36), jenis_sppd varchar(2), id_pegawai_pengaju varchar(36),
			tgl_berangkat date, tgl_kembali date, id_fasilitas_transport varchar(36), is_rombongan boolean, 
			id_branch varchar(36), created_at timestamp, is_antar boolean, is_jemput boolean, operasional_hari_dinas boolean,
			id_rule_approval varchar(36), id_jenis_approval varchar(36), is_pengajuan boolean, status varchar(2),
			nip varchar,  nama_pegawai varchar, id_jabatan varchar,id_person_grade varchar, jabatan varchar, id_level_bod varchar, file text, link_file text,
			kode varchar, person_grade varchar,nama_jenis_tujuan varchar, keterangan text,  nama_fasilitas_transport varchar,
		   kode_branch varchar,nama_branch varchar, nama_approval varchar,id_pengajuan_bpd_histori varchar, is_deleted boolean,
		  status_bpd varchar, type_approval varchar, id_pegawai_approval varchar, esign boolean, is_max_pengajuan boolean, is_max_penyelesaian boolean, nama_fungsionalitas varchar
	  )
	  where spd.status in ('0','1','2','3')
	   and spd.tgl_berangkat::date>= $3::date
	   and spd.tgl_berangkat::date<= $4::date
		`,
		SelectDataBpdNew: `
		SELECT * FROM fn_list_bpd_selesai($1, $2, 'PENGAJUANA_BPD') as pd (
			id varchar, nomor varchar, nama varchar, tujuan varchar, keperluan varchar, tgl_berangkat date, tgl_kembali date, id_jenis_perjalanan_dinas varchar, id_jenis_kendaraan varchar, is_rombongan boolean, 
			status varchar, id_rule_approval varchar, created_at timestamp, created_by varchar, updated_by varchar, updated_at timestamp, is_deleted boolean, nama_jenis_perjalanan_dinas varchar, nama_jenis_kendaraan varchar, 
			nama_kendaraan varchar, nama_approval varchar, id_bpd_pegawai varchar, id_pegawai varchar, nip varchar, nama_pegawai varchar, id_bidang varchar, id_jenis_approval varchar, id_pegawai_pengaju varchar, id_pengajuan_bpd_histori varchar, 
			file varchar, id_sppd varchar, spp_id integer, status_bpd varchar, type_approval varchar, id_pegawai_approval varchar, esign boolean, is_max_pengajuan boolean, is_max_penyelesaian boolean, keterangan_tujuan text, nama_fungsionalitas varchar, id_jenis_tujuan varchar, jenis_sppd varchar
		)
	  where pd.status in ('0','1','2','3')
	  and pd.tgl_berangkat::date>= $3::date
	  and pd.tgl_berangkat::date<= $4::date
		`,
		SelectDashboardSppd: `select * from fn_dashboard_sppd($1, $2, $3, $4, $5) as (
			belum_proses_pengajuan bigint, pengajuan bigint, pengajuan_disetujui bigint, revisi bigint)
			`,
		SelectDashboardBpdNew: `select * from fn_dashboard_bpd_new($1, $2, $3, $4, $5) as (
				belum_proses_pengajuan bigint, pengajuan bigint, pengajuan_disetujui bigint, revisi bigint)
				`,
		SelectJumlahSppd: `
		SELECT count(spd.id) as pengajuan FROM fn_list_sppd($1, $2, 'PENGAJUAN_SPPD', $3) as spd  (
			id varchar(36), tgl_surat date, nomor_surat varchar(100), id_pegawai varchar(36), jenis_tujuan varchar(36), 
			tujuan_dinas text, keperluan_dinas text, tenant_id varchar(36), jenis_sppd varchar(2), id_pegawai_pengaju varchar(36),
			tgl_berangkat date, tgl_kembali date, id_fasilitas_transport varchar(36), is_rombongan boolean, 
			id_branch varchar(36), created_at timestamp, is_antar boolean, is_jemput boolean, operasional_hari_dinas boolean,
			id_rule_approval varchar(36), id_jenis_approval varchar(36), is_pengajuan boolean, status varchar(2),
			nip varchar,  nama_pegawai varchar, id_jabatan varchar,id_person_grade varchar, jabatan varchar,  id_level_bod varchar, file text, link_file text,
			kode varchar, person_grade varchar,nama_jenis_tujuan varchar, keterangan text,  nama_fasilitas_transport varchar,
			kode_branch varchar,nama_branch varchar, nama_approval varchar,id_pengajuan_bpd_histori varchar, is_deleted boolean, bidang varchar,
			status_bpd varchar, type_approval varchar, id_pegawai_approval varchar, esign boolean, is_max_pengajuan boolean, is_max_penyelesaian boolean, nama_fungsionalitas varchar
			)
			where coalesce(spd.is_deleted, false) = false
			and spd.id_pengajuan_bpd_histori is not null
			and (case when spd.status='5' then spd.status else spd.status_bpd end) ='1'			  
	`,
		SelectJumlahBpd: `
	SELECT count(pd.id) as pengajuan FROM fn_list_bpd_selesai($1, $2, 'PENGAJUAN_BPD') as pd (
		id varchar, nomor varchar, nama varchar, tujuan varchar, keperluan varchar, tgl_berangkat date, tgl_kembali date, id_jenis_perjalanan_dinas varchar, id_jenis_kendaraan varchar, is_rombongan boolean, 
		status varchar, id_rule_approval varchar, created_at timestamp, created_by varchar, updated_by varchar, updated_at timestamp, is_deleted boolean, nama_jenis_perjalanan_dinas varchar, nama_jenis_kendaraan varchar, 
		nama_kendaraan varchar, nama_approval varchar, id_bpd_pegawai varchar, id_pegawai varchar, nip varchar, nama_pegawai varchar, id_bidang varchar, id_jenis_approval varchar, id_pegawai_pengaju varchar, id_pengajuan_bpd_histori varchar, 
		file varchar, id_sppd varchar, spp_id integer, status_bpd varchar, type_approval varchar, id_pegawai_approval varchar, esign boolean, is_max_pengajuan boolean, is_max_penyelesaian boolean, keterangan_tujuan text, nama_fungsionalitas varchar, id_jenis_tujuan varchar, jenis_sppd varchar
	)
	`,
	}
)

type DashboardRepository interface {
	GetAll() (dataDashboard []Dashboard, err error)
	GetDataDashboardBpd(req DashboardRequest) (data DashboardBpd, err error)
	GetDataBpd(req DashboardRequest) (data []DataAktifBpd, err error)
	GetDataDashboardSppd(req DashboardRequest) (data DashboardSppd, err error)
	GetDataDashboardBpdNew(req DashboardRequest) (data DashboardSppd, err error)
	GetDataSppd(req DashboardRequest) (data []DataAktifSppd, err error)
	GetDataBpdNew(req DashboardRequest) (data []DataAktifBpdNew, err error)
	GetJumlahSppd(req DashboardRequest) (data JumlahSppd, err error)
	GetJumlahBpd(req DashboardRequest) (data JumlahSppd, err error)
}

type DashboardRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideDashboardRepositoryPostgreSQL(db *infras.PostgresqlConn) *DashboardRepositoryPostgreSQL {
	s := new(DashboardRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *DashboardRepositoryPostgreSQL) GetAll() (dataDashboard []Dashboard, err error) {
	criteria := ` WHERE mp.is_deleted = false `

	err = r.DB.Read.Select(&dataDashboard, dashboardQuery.SelectJmlPegawai+criteria)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *DashboardRepositoryPostgreSQL) GetDataDashboardBpd(req DashboardRequest) (data DashboardBpd, err error) {
	err = r.DB.Read.Get(&data, dashboardQuery.SelectDashboard, req.IdPegawai, req.IdPegawaiApproval, req.IdBidang)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *DashboardRepositoryPostgreSQL) GetDataBpd(req DashboardRequest) (data []DataAktifBpd, err error) {
	rows, err := r.DB.Read.Queryx(dashboardQuery.SelectDataBpd, req.IdPegawai, req.IdPegawaiApproval)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var bpd DataAktifBpd
		err = rows.StructScan(&bpd)

		if err != nil {
			return
		}

		data = append(data, bpd)
	}
	return
}

func (r *DashboardRepositoryPostgreSQL) GetDataDashboardSppd(req DashboardRequest) (data DashboardSppd, err error) {
	err = r.DB.Read.Get(&data, dashboardQuery.SelectDashboardSppd, req.IdPegawai, req.IdPegawaiApproval, req.IdBidang, req.StartDate, req.EndDate)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *DashboardRepositoryPostgreSQL) GetDataDashboardBpdNew(req DashboardRequest) (data DashboardSppd, err error) {
	err = r.DB.Read.Get(&data, dashboardQuery.SelectDashboardBpdNew, req.IdPegawai, req.IdPegawaiApproval, req.IdBidang, req.StartDate, req.EndDate)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *DashboardRepositoryPostgreSQL) GetDataSppd(req DashboardRequest) (data []DataAktifSppd, err error) {
	rows, err := r.DB.Read.Queryx(dashboardQuery.SelectDataSppd, req.IdPegawai, req.IdPegawaiApproval, req.StartDate, req.EndDate)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var sppd DataAktifSppd
		err = rows.StructScan(&sppd)

		if err != nil {
			return
		}

		data = append(data, sppd)
	}
	return
}

func (r *DashboardRepositoryPostgreSQL) GetDataBpdNew(req DashboardRequest) (data []DataAktifBpdNew, err error) {
	rows, err := r.DB.Read.Queryx(dashboardQuery.SelectDataBpdNew, req.IdPegawai, req.IdPegawaiApproval, req.StartDate, req.EndDate)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var bpdNew DataAktifBpdNew
		err = rows.StructScan(&bpdNew)

		if err != nil {
			return
		}

		data = append(data, bpdNew)
	}
	return
}

func (r *DashboardRepositoryPostgreSQL) GetJumlahSppd(req DashboardRequest) (data JumlahSppd, err error) {
	err = r.DB.Read.Get(&data, dashboardQuery.SelectJumlahSppd, req.IdPegawai, req.IdPegawaiApproval, req.IdBidang)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *DashboardRepositoryPostgreSQL) GetJumlahBpd(req DashboardRequest) (data JumlahSppd, err error) {
	criteria := `where coalesce(pd.is_deleted, false) = false
	and pd.id_pengajuan_bpd_histori is not null
	and (case when pd.status='5' then pd.status else pd.status_bpd end) ='1'			   `

	if req.IdBidang != "" {
		criteria += fmt.Sprintf(" and pd.id_bidang = '%s' ", req.IdBidang)
	}
	err = r.DB.Read.Get(&data, dashboardQuery.SelectJumlahBpd+criteria, req.IdPegawai, req.IdPegawaiApproval)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}
