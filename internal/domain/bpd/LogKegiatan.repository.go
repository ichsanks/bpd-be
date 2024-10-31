package bpd

import (
	"database/sql"
	"fmt"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
)

var (
	logKegiatanQuery = struct {
		Select,
		Insert,
		Delete string
	}{
		Select: `SELECT id, tanggal, id_perjalanan_dinas, id_bpd_pegawai, foto, keterangan, lat, long, address, created_at, created_by, updated_at, updated_by, is_deleted FROM public.log_kegiatan`,
		Insert: `INSERT INTO public.log_kegiatan(
			id, tanggal, id_perjalanan_dinas, id_bpd_pegawai, foto, keterangan, lat, long, address, created_at, created_by)
			VALUES (:id, :tanggal, :id_perjalanan_dinas, :id_bpd_pegawai, :foto, :keterangan, :lat, :long, :address, :created_at, :created_by)`,
		Delete: `DELETE from public.log_kegiatan`,
	}
)

var (
	bpdDokumenQuery = struct {
		Select,
		Insert,
		Delete string
	}{
		Select: `SELECT id, id_bpd_pegawai, file, keterangan, created_at, created_by, updated_at, updated_by, is_deleted FROM public.perjalanan_dinas_dokumen`,
		Insert: `INSERT INTO public.perjalanan_dinas_dokumen(
			id, id_bpd_pegawai, file, keterangan, created_at, created_by)
			VALUES (:id, :id_bpd_pegawai, :file, :keterangan, :created_at, :created_by)`,
		Delete: `DELETE from public.perjalanan_dinas_dokumen`,
	}
)

type LogKegiatanRepository interface {
	Create(data LogKegiatan) error
	GetAll(idPerjalananDinas string, idBpdPegawai string) (data []LogKegiatan, err error)
	DeleteByID(id string) (err error)
	ResolveByID(id string) (data LogKegiatan, err error)
	// Dokumen BPD
	CreateDokumen(data PerjalananDinasDokumen) error
	GetAllDokumen(idBpdPegawai string) (data []PerjalananDinasDokumen, err error)
	DeleteDokumen(id string) (err error)
	ResolveByIDDokumen(id string) (data PerjalananDinasDokumen, err error)
}

type LogKegiatanRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideLogKegiatanRepositoryPostgreSQL(db *infras.PostgresqlConn) *LogKegiatanRepositoryPostgreSQL {
	s := new(LogKegiatanRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *LogKegiatanRepositoryPostgreSQL) Create(data LogKegiatan) error {
	stmt, err := r.DB.Write.PrepareNamed(logKegiatanQuery.Insert)
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

func (r *LogKegiatanRepositoryPostgreSQL) GetAll(idPerjalananDinas string, idBpdPegawai string) (data []LogKegiatan, err error) {
	where := " where coalesce(is_deleted,false)=false "
	if idPerjalananDinas != "" {
		where += fmt.Sprintf(" and id_perjalanan_dinas='%v' ", idPerjalananDinas)
	}

	if idBpdPegawai != "" {
		where += fmt.Sprintf(" and id_bpd_pegawai='%v' ", idBpdPegawai)
	}
	where += " order by created_at asc "

	rows, err := r.DB.Read.Queryx(logKegiatanQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Log Kegiatan NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items LogKegiatan
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	return
}

func (r *LogKegiatanRepositoryPostgreSQL) ResolveByID(id string) (data LogKegiatan, err error) {
	err = r.DB.Read.Get(&data, logKegiatanQuery.Select+" WHERE id=$1 ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *LogKegiatanRepositoryPostgreSQL) DeleteByID(id string) (err error) {
	stmt, err := r.DB.Read.PrepareNamed(logKegiatanQuery.Delete + " where id=:id")
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	arg := map[string]interface{}{
		"id": id,
	}

	defer stmt.Close()
	_, err = stmt.Exec(arg)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

// ==================================================== BPD Dokumen ==============================================================

func (r *LogKegiatanRepositoryPostgreSQL) CreateDokumen(data PerjalananDinasDokumen) error {
	stmt, err := r.DB.Write.PrepareNamed(bpdDokumenQuery.Insert)
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

func (r *LogKegiatanRepositoryPostgreSQL) GetAllDokumen(idBpdPegawai string) (data []PerjalananDinasDokumen, err error) {
	rows, err := r.DB.Read.Queryx(bpdDokumenQuery.Select+" where id_bpd_pegawai=$1 order by created_at asc", idBpdPegawai)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Perjalanan Dinas Dokumen NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items PerjalananDinasDokumen
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	return
}

func (r *LogKegiatanRepositoryPostgreSQL) ResolveByIDDokumen(id string) (data PerjalananDinasDokumen, err error) {
	err = r.DB.Read.Get(&data, bpdDokumenQuery.Select+" WHERE id=$1 ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *LogKegiatanRepositoryPostgreSQL) DeleteDokumen(id string) error {
	stmt, err := r.DB.Read.PrepareNamed(bpdDokumenQuery.Delete + " where id=:id")
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	arg := map[string]interface{}{
		"id": id,
	}

	defer stmt.Close()
	_, err = stmt.Exec(arg)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}
