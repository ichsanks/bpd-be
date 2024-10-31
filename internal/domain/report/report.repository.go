package report

import (
	"database/sql"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
)

var (
	reportQuery = struct {
		SelectRekap,
		SelectRekapBiayaReimbusment,
		SelectRekapBiayaAkomodasi,
		SelectRekapBagian,
		SelectRekapTotal string
	}{
		SelectRekap: `SELECT * FROM fn_rpt_rekap_biaya_bpd($1, $2, $3, $4) as (
			id varchar, id_perjalanan_dinas varchar, nomor varchar, nip varchar, nama varchar, tgl_berangkat date, tgl_kembali date, nama_unor varchar, nama_bidang varchar, nama_jabatan varchar, 
			kode_golongan varchar, nama_golongan varchar, nama_dinas varchar, tujuan varchar, keperluan varchar, biaya_dinas double precision, reimbursement double precision, akomodasi double precision) `,
		SelectRekapBagian: `SELECT * FROM fn_rpt_rekap_bpd_bagian($1, $2, $3) as (
			id varchar, nomor varchar, nama_bpd varchar, tujuan varchar, keperluan varchar, tgl_berangkat date, tgl_kembali date, 
			is_rombongan boolean, jenis_bpd varchar, jenis_kendaraan varchar, nama_bidang varchar, total double precision
		)`,
		SelectRekapTotal: `SELECT * FROM fn_rpt_rekap_total_bpd($1, $2, $3) as (
			id_bidang varchar, bidang varchar, id_unor varchar, unor varchar, jml_pengajuan integer, 
			jml_sdh_penyelesaian integer, jml_blm_penyelesaian integer, total double precision
		)`,
		SelectRekapBiayaReimbusment: `SELECT * FROM fn_rpt_rekap_biaya_reimbusment($1, $2, $3) as (
			id_perjalanan_dinas varchar, id_pegawai varchar, nip varchar, nama varchar, nomor varchar, nama_bpd varchar, tujuan varchar, keperluan varchar,
			tgl_berangkat date, tgl_kembali date,  is_rombongan boolean, jml_hari integer,
			id_bidang varchar, nama_bidang varchar, nama_jabatan varchar, status varchar, total double precision 
		)`,
		SelectRekapBiayaAkomodasi: `SELECT * FROM fn_rpt_rekap_biaya_akomodasi($1, $2, $3, $4) as (
			id varchar, id_perjalanan_dinas varchar, nomor varchar, nip varchar, nama varchar, tgl_berangkat date, tgl_kembali date, 
			nama_bidang varchar, nama_jabatan varchar, kode_golongan varchar, nama_golongan varchar, nama_dinas varchar, tujuan varchar, 
			keperluan varchar, total_akomodasi double precision, details json) `,
	}
)

type ReportRepository interface {
	RptRekapBpd(req FilterReport) (data []ReportRekapBpd, err error)
	RptRekapBpdBagian(req FilterReport) (data []ReportRekapBpdBagian, err error)
	RptRekapTotalBpd(req FilterReport) (data []ReportRekapTotalBpd, err error)
	RptRekapAkomodasi(req FilterReport) (data []RekapBiayaAkomodasi, err error)
	RptRekapReimbusment(req FilterReport) (data []ReportRekapAkReim, err error)
}

type ReportRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideReportRepositoryPostgreSQL(db *infras.PostgresqlConn) *ReportRepositoryPostgreSQL {
	s := new(ReportRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *ReportRepositoryPostgreSQL) RptRekapBpd(req FilterReport) (data []ReportRekapBpd, err error) {
	rows, err := r.DB.Read.Queryx(reportQuery.SelectRekap, req.IdUnor, req.IdBidang, req.TglAwal, req.TglAkhir)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Report NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items ReportRekapBpd
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	return
}

func (r *ReportRepositoryPostgreSQL) RptRekapBpdBagian(req FilterReport) (data []ReportRekapBpdBagian, err error) {
	rows, err := r.DB.Read.Queryx(reportQuery.SelectRekapBagian, req.IdBidang, req.TglAwal, req.TglAkhir)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Report NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items ReportRekapBpdBagian
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	return
}

func (r *ReportRepositoryPostgreSQL) RptRekapTotalBpd(req FilterReport) (data []ReportRekapTotalBpd, err error) {
	rows, err := r.DB.Read.Queryx(reportQuery.SelectRekapTotal, req.IdBidang, req.TglAwal, req.TglAkhir)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Report NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items ReportRekapTotalBpd
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	return
}

func (r *ReportRepositoryPostgreSQL) RptRekapReimbusment(req FilterReport) (data []ReportRekapAkReim, err error) {
	rows, err := r.DB.Read.Queryx(reportQuery.SelectRekapBiayaReimbusment, req.IdBidang, req.TglAwal, req.TglAkhir)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Report NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items ReportRekapAkReim
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	return
}

func (r *ReportRepositoryPostgreSQL) RptRekapAkomodasi(req FilterReport) (data []RekapBiayaAkomodasi, err error) {
	rows, err := r.DB.Read.Queryx(reportQuery.SelectRekapBiayaAkomodasi, req.IdUnor, req.IdBidang, req.TglAwal, req.TglAkhir)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Report NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items RekapBiayaAkomodasi
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	return
}
