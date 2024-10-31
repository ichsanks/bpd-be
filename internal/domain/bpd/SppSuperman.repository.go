package bpd

import (
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
)

var (
	supermanQuery = struct {
		SelectNoSPPB           string
		SelectNoRomawi         string
		UpdateStatusBpdPegawai string
	}{
		SelectNoSPPB:   `select public.fn_get_nomor_sppb($1,$2);`,
		SelectNoRomawi: `SELECT * FROM romawi($1)`,
		UpdateStatusBpdPegawai: `UPDATE perjalanan_dinas_pegawai SET 
		id=:id, 
		is_sppb=:is_sppb`,
	}
)

type SppSupermanRepository interface {
	GetNomorUrutSPPb(kode int, bagian string) (nomor string, err error)
	GetNomorRomawi(kode string) (nomor string, err error)
	UpdateStatusBPDPegawai(data StatusUpdateBpdSuperman) (err error)
}

type SppSupermanRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideSppSupermanRepositoryPostgreSQL(db *infras.PostgresqlConn) *SppSupermanRepositoryPostgreSQL {
	s := new(SppSupermanRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *SppSupermanRepositoryPostgreSQL) GetNomorUrutSPPb(kode int, bagian string) (nomor string, err error) {
	err = r.DB.Read.Get(&nomor, supermanQuery.SelectNoSPPB, kode, bagian)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}
func (r *SppSupermanRepositoryPostgreSQL) GetNomorRomawi(kode string) (nomor string, err error) {
	err = r.DB.Read.Get(&nomor, supermanQuery.SelectNoRomawi, kode)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *SppSupermanRepositoryPostgreSQL) UpdateStatusBPDPegawai(data StatusUpdateBpdSuperman) (err error) {
	stmt, err := r.DB.Write.PrepareNamed(supermanQuery.UpdateStatusBpdPegawai + " WHERE id=:id")
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return nil
}
