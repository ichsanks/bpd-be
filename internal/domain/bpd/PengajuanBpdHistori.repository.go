package bpd

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
)

var (
	pengajuanBpdHistoriQuery = struct {
		Select,
		SelectDTO,
		SelectUnor,
		SelectTimeline,
		SelectTimelineTtd,
		Insert,
		UpdateApprove,
		UpdateRevisi,
		UpdatePerjalananDinas string
		UpdateBpdPegawai string
	}{
		Select: `SELECT id, tanggal, id_perjalanan_dinas, id_pegawai, id_fungsionalitas, id_unor, id_rule_approval_detail, catatan, keterangan, status, type_approval, 
		created_at, created_by, approved_at, approved_by, id_bpd_histori_revisi, id_approval_line, group_approval, id_bpd_pegawai 
		FROM public.pengajuan_bpd_histori`,
		SelectDTO: `SELECT h.id, h.tanggal, h.id_perjalanan_dinas, h.id_pegawai, h.id_fungsionalitas, h.id_unor, h.id_rule_approval_detail, h.catatan, 
		h.keterangan, h.status, h.type_approval, h.created_at, h.created_by, h.approved_at, h.approved_by, p.nip, p.nama nama_pegawai,
		f.nama nama_fungsionalitas, u.nama nama_unor, u.kode kode_unor, h.id_bpd_histori_revisi
		FROM public.pengajuan_bpd_histori h
		left join m_pegawai p on p.id = h.id_pegawai
		left join m_fungsionalitas f on f.id = h.id_fungsionalitas
		left join m_unit_organisasi_kerja u on u.id = h.id_unor
		`,
		SelectUnor: `SELECT u.id_unor, u.kode_unor, u.nama_unor from (
				SELECT x.id_unor, x.kode_unor, x.nama_unor, count(x.kode_unor) jml FROM (
					select u.id id_unor, u.kode kode_unor, u.nama nama_unor from m_pegawai p
					left join m_unit_organisasi_kerja u on u.id = p.id_unor
					where coalesce(p.is_deleted,false)=false
					and p.id_fungsionalitas = $1
					and (u.kode ilike ''||(
						SELECT SUBSTRING ($2, 1, char_length($2)-2)
					)||'%')
					UNION ALL
					select u.id id_unor, u.kode kode_unor, u.nama nama_unor from m_pegawai p
					left join m_unit_organisasi_kerja u on u.id = p.id_unor
					where coalesce(p.is_deleted,false)=false
					and p.id_fungsionalitas = $1
					and u.kode ilike ''||$2||'%' 
				)x
				group by x.id_unor, x.kode_unor, x.nama_unor
			)u order by jml desc
			limit 1
		`,
		SelectTimeline: `SELECT * FROM fn_get_timeline($1, $2) as (
				id varchar, tanggal text, id_perjalanan_dinas varchar, id_pegawai varchar, id_fungsionalitas varchar, id_unor varchar, id_rule_approval_detail varchar, catatan text, 
				keterangan varchar, status varchar, type_approval varchar, nip varchar, nama_pegawai varchar, nama_fungsionalitas varchar, nama_unor varchar, kode_unor varchar, 
				nama_bidang varchar, nama_jabatan varchar, created_at text, approved_at text, nama_pengaju varchar, is_pengaju boolean, ket_status text
			)`,
		SelectTimelineTtd: `SELECT * FROM fn_get_timeline_ttd($1, $2) as (
				id varchar, tanggal text, id_perjalanan_dinas varchar, id_pegawai varchar, id_fungsionalitas varchar, id_unor varchar, id_rule_approval_detail varchar, catatan text, 
				keterangan varchar, status varchar, type_approval varchar, nip varchar, nama_pegawai varchar, nama_fungsionalitas varchar, nama_unor varchar, kode_unor varchar, 
				nama_bidang varchar, nama_jabatan varchar, created_at text, approved_at text, nama_pengaju varchar, is_pengaju boolean, ket_status text, ket_ttd text 
			)`,
		Insert: `INSERT INTO public.pengajuan_bpd_histori(
			id, tanggal, id_perjalanan_dinas, id_pegawai, id_fungsionalitas, id_unor, id_rule_approval_detail, catatan, keterangan, status, type_approval, created_at, created_by, id_bpd_histori_revisi, id_approval_line, group_approval, id_bpd_pegawai)
			VALUES (:id, :tanggal, :id_perjalanan_dinas, :id_pegawai, :id_fungsionalitas, :id_unor, :id_rule_approval_detail, :catatan, :keterangan, :status, :type_approval, :created_at, :created_by, :id_bpd_histori_revisi, :id_approval_line, :group_approval, :id_bpd_pegawai)`,
		UpdateApprove: `update pengajuan_bpd_histori set
				id=:id,
				id_pegawai=:id_pegawai,
				catatan=:catatan,
				keterangan=:keterangan,
				status=:status,
				approved_at=:approved_at,
				approved_by=:approved_by `,
		UpdatePerjalananDinas: `UPDATE perjalanan_dinas SET 
		        id=:id, 
				status=:status `,
		UpdateBpdPegawai: `UPDATE perjalanan_dinas_pegawai SET 
		        id=:id, 
				status=:status `,
		UpdateRevisi: `update pengajuan_bpd_histori set
				is_revisi=:is_revisi`,
	}
)

var (
	pengajuanHistoridetailQuery = struct {
		Select,
		InsertBulk,
		InsertBulkPlaceholder string
	}{
		Select:                `SELECT id, id_pengajuan, id_pegawai, created_at, created_by FROM public.pengajuan_bpd_histori_detail`,
		InsertBulk:            `INSERT INTO public.pengajuan_bpd_histori_detail(id, id_pengajuan, id_pegawai, created_at, created_by) values `,
		InsertBulkPlaceholder: ` (:id, :id_pengajuan, :id_pegawai, :created_at, :created_by) `,
	}
)

type PengajuanBpdHistoriRepository interface {
	Create(data PengajuanBpdHistori) error
	UpdateApproval(data PengajuanBpdHistori) (err error)
	UpdateStatusBPD(data StatusBPD) (err error)
	UpdateStatusBPDPegawai(data StatusBPD) (err error)
	ResolveByID(id string) (data PengajuanBpdHistori, err error)
	ResolveUnorPegawai(idFungsionalitas string, kodeUnor string) (data UnorPegawai, err error)
	GetAll(idPerjalananDinas string, idBpdPegawai string, typeApproval string) (data []PengajuanBpdHistori, err error)
	GetTimeline(idPerjalananDinas string, idBpdPegawai string) (data []Timeline, err error)
	GetPreviousBpdHistori(idPerjalananDinas string, idBpdPegawai string, idBpdHistori string) (data PengajuanBpdHistori, err error)
	ResolveByIDDTO(id string) (data PengajuanBpdHistoriDTO, err error)
	UpdateRevisi(data StatusRevisi) (err error)
	GetTimelineTtd(idPerjalananDinas string, idBpdPegawai string) (data []Timeline, err error)
}

type PengajuanBpdHistoriRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvidePengajuanBpdHistoriRepositoryPostgreSQL(db *infras.PostgresqlConn) *PengajuanBpdHistoriRepositoryPostgreSQL {
	s := new(PengajuanBpdHistoriRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) Create(data PengajuanBpdHistori) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		// Function create table pengajuan_histori
		if err := r.CreateTx(tx, data); err != nil {
			e <- err
			return
		}

		// Function Insert Bulk table pengajuan_histori_detail
		if err := txCreatePengajuanDetail(tx, data.Detail); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) CreateTx(tx *sqlx.Tx, data PengajuanBpdHistori) error {
	stmt, err := tx.PrepareNamed(pengajuanBpdHistoriQuery.Insert)
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

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) UpdateApproval(data PengajuanBpdHistori) (err error) {
	stmt, err := r.DB.Write.PrepareNamed(pengajuanBpdHistoriQuery.UpdateApprove + " WHERE id=:id")
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

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) UpdateStatusBPD(data StatusBPD) (err error) {
	stmt, err := r.DB.Write.PrepareNamed(pengajuanBpdHistoriQuery.UpdatePerjalananDinas + " WHERE id=:id")
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

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) UpdateStatusBPDPegawai(data StatusBPD) (err error) {
	stmt, err := r.DB.Write.PrepareNamed(pengajuanBpdHistoriQuery.UpdateBpdPegawai + " WHERE id=:id")
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

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) GetAll(idPerjalananDinas string, idBpdPegawai string, typeApproval string) (data []PengajuanBpdHistori, err error) {
	criteria := ` where id_perjalanan_dinas=$1 `
	if typeApproval != "" {
		criteria += fmt.Sprintf(` and type_approval='%s' `, typeApproval)
	}

	if idBpdPegawai != "" {
		criteria += fmt.Sprintf(` and id_bpd_pegawai='%s' `, idBpdPegawai)
	}
	criteria += ` order by tanggal asc`

	rows, err := r.DB.Read.Queryx(pengajuanBpdHistoriQuery.Select+criteria, idPerjalananDinas)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Pengajuan BPD Histori NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items PengajuanBpdHistori
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	if data == nil {
		data = make([]PengajuanBpdHistori, 0)
	}

	return
}

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) GetTimeline(idPerjalananDinas string, idBpdPegawai string) (data []Timeline, err error) {
	rows, err := r.DB.Read.Queryx(pengajuanBpdHistoriQuery.SelectTimeline, idPerjalananDinas, idBpdPegawai)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Pengajuan BPD Histori NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items Timeline
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	if data == nil {
		data = make([]Timeline, 0)
	}

	return
}

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) ResolveByID(id string) (data PengajuanBpdHistori, err error) {
	err = r.DB.Read.Get(&data, pengajuanBpdHistoriQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) ResolveByIDDTO(id string) (data PengajuanBpdHistoriDTO, err error) {
	err = r.DB.Read.Get(&data, pengajuanBpdHistoriQuery.SelectDTO+" WHERE h.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) ResolveUnorPegawai(idFungsionalitas string, kodeUnor string) (data UnorPegawai, err error) {
	err = r.DB.Read.Get(&data, pengajuanBpdHistoriQuery.SelectUnor, idFungsionalitas, kodeUnor)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) GetPreviousBpdHistori(idPerjalananDinas string, idBpdPegawai string, idBpdHistori string) (data PengajuanBpdHistori, err error) {
	criteria := ` where id_perjalanan_dinas = $1 `
	if idBpdPegawai != "" {
		criteria += fmt.Sprintf(" and id_bpd_pegawai='%v' ", idBpdPegawai)
	}

	criteria += ` and tanggal < (select tanggal from pengajuan_bpd_histori where id = $2)
		order by tanggal desc
		limit 1	
	`
	err = r.DB.Read.Get(&data, pengajuanBpdHistoriQuery.Select+criteria, idPerjalananDinas, idBpdHistori)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func txCreatePengajuanDetail(tx *sqlx.Tx, details []PengajuanBpdHistoriDetail) (err error) {
	if len(details) == 0 {
		return
	}
	query, args, err := composeBulkUpsertPengajuanDetailQuery(details)
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

func composeBulkUpsertPengajuanDetailQuery(details []PengajuanBpdHistoriDetail) (qResult string, params []interface{}, err error) {
	values := []string{}
	for _, d := range details {
		param := map[string]interface{}{
			"id":           d.ID,
			"id_pengajuan": d.IdPengajuan,
			"id_pegawai":   d.IdPegawai,
			"created_at":   d.CreatedAt,
			"created_by":   d.CreatedBy,
		}
		q, args, err := sqlx.Named(pengajuanHistoridetailQuery.InsertBulkPlaceholder, param)
		if err != nil {
			return qResult, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}
	qResult = fmt.Sprintf(`%v %v`, pengajuanHistoridetailQuery.InsertBulk, strings.Join(values, ","))
	return
}

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) UpdateRevisi(data StatusRevisi) (err error) {
	stmt, err := r.DB.Write.PrepareNamed(pengajuanBpdHistoriQuery.UpdateRevisi + " WHERE id_perjalanan_dinas=:id_perjalanan_dinas")
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

func (r *PengajuanBpdHistoriRepositoryPostgreSQL) GetTimelineTtd(idPerjalananDinas string, idBpdPegawai string) (data []Timeline, err error) {
	rows, err := r.DB.Read.Queryx(pengajuanBpdHistoriQuery.SelectTimelineTtd, idPerjalananDinas, idBpdPegawai)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Pengajuan BPD Histori NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items Timeline
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	if data == nil {
		data = make([]Timeline, 0)
	}

	return
}
