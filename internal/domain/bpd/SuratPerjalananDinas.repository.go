package bpd

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

var (
	suratPerjalananDinasQuery = struct {
		Select             string
		Insert             string
		Update             string
		SelectNoSurat      string
		SelectDto          string
		SelectDto2         string
		CountDto           string
		SelectListDto      string
		UpdateFileSppd     string
		UpdateLinkFileSppd string
	}{
		Select: `select 
					spd.id, spd.tgl_surat, spd.nomor_surat, spd.id_pegawai, spd.jenis_tujuan, spd.tujuan_dinas,
					spd.keperluan_dinas, spd.tgl_berangkat, spd.tgl_kembali, spd.id_fasilitas_transport, spd.is_rombongan,
					spd.created_at, spd.created_by, spd.updated_at, spd.updated_by, spd.is_deleted, spd.id_branch, 
					spd.is_antar, spd.is_jemput, spd.operasional_hari_dinas, spd.id_pegawai_pengaju, 
					spd.jenis_sppd, spd.tenant_id, spd.id_rule_approval, spd.type_approval, spd.id_jenis_approval, 
					spd.is_pengajuan, spd.status, spd.file
				from 
					public.surat_perjalanan_dinas spd `,
		Insert: `INSERT INTO public.surat_perjalanan_dinas 
					(id, tgl_surat, nomor_surat, id_pegawai, jenis_tujuan, tujuan_dinas, keperluan_dinas, 
					tgl_berangkat, tgl_kembali, id_fasilitas_transport, is_rombongan, is_jemput, is_antar, 
					operasional_hari_dinas, id_branch, created_at, created_by, tenant_id, jenis_sppd, id_pegawai_pengaju,
					id_rule_approval, type_approval, id_jenis_approval, is_pengajuan, status) 
				values
					(:id, :tgl_surat, :nomor_surat, :id_pegawai, :jenis_tujuan, :tujuan_dinas, :keperluan_dinas, 
					:tgl_berangkat, :tgl_kembali, :id_fasilitas_transport, :is_rombongan, :is_jemput, :is_antar, 
					:operasional_hari_dinas, :id_branch, :created_at, :created_by, :tenant_id, :jenis_sppd, 
					:id_pegawai_pengaju, :id_rule_approval, :type_approval, :id_jenis_approval, :is_pengajuan, :status) `,
		Update: `UPDATE public.surat_perjalanan_dinas SET 
		        id=:id, 
				tgl_surat=:tgl_surat,
				nomor_surat=:nomor_surat, 
				id_pegawai=:id_pegawai,
				jenis_tujuan=:jenis_tujuan,
				tujuan_dinas=:tujuan_dinas,
				keperluan_dinas=:keperluan_dinas, 
				tgl_berangkat=:tgl_berangkat,
				tgl_kembali=:tgl_kembali,
				id_fasilitas_transport=:id_fasilitas_transport, 
				is_rombongan=:is_rombongan,
				id_branch=:id_branch,
				is_antar=:is_antar,
				is_jemput=:is_jemput,
				operasional_hari_dinas=:operasional_hari_dinas,
				id_rule_approval=:id_rule_approval, 
				type_approval=:type_approval, 
				id_jenis_approval=:id_jenis_approval, 
				is_pengajuan=:is_pengajuan,
				tenant_id=:tenant_id,
				jenis_sppd=:jenis_sppd,
				id_pegawai_pengaju=:id_pegawai_pengaju,
				status=:status,
				file=:file,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		SelectNoSurat: `SELECT fn_get_no_surat() `,
		SelectDto: `select 
						spd.id, spd.tgl_surat, spd.nomor_surat, spd.id_pegawai, spd.jenis_tujuan, 
						spd.tujuan_dinas, spd.keperluan_dinas, spd.tenant_id, spd.jenis_sppd, spd.id_pegawai_pengaju,
						spd.tgl_berangkat, spd.tgl_kembali, spd.id_fasilitas_transport, spd.is_rombongan, 
						spd.id_branch, spd.created_at, spd.is_antar, spd.is_jemput, spd.operasional_hari_dinas,
						spd.id_rule_approval, spd.type_approval, spd.id_jenis_approval, spd.is_pengajuan, spd.status,
						p.nip, p.nama as nama_pegawai, p.id_jabatan, p.id_person_grade, id_level_bod,
						j.nama as jabatan,  
						pg.kode, pg.nama as person_grade,
						jt.nama as nama_jenis_tujuan, jt.keterangan, ft.nama as nama_fasilitas_transport,
						b.kode as kode_branch, b.nama as nama_branch, mb.nama bidang
					from 
						public.surat_perjalanan_dinas spd
					left join 
						m_pegawai p on p.id = spd.id_pegawai
						left join 
						m_bidang mb on mb.id = p.id_bidang
					left join 
						m_jabatan j on j.id = p.id_jabatan
					left join 
						m_person_grade pg on pg.id = p.id_person_grade
					left join 
						m_jenis_tujuan jt on jt.id = spd.jenis_tujuan
					left join 
						m_fasilitas_transport ft on ft.id = spd.id_fasilitas_transport
					left join 
						m_branch b on b.id = spd.id_branch `,
		CountDto: `select 
						count(spd.id)
					from 
						public.surat_perjalanan_dinas spd
					left join 
						m_pegawai p on p.id = spd.id_pegawai
					left join 
						m_jabatan j on j.id = p.id_jabatan
					left join 
						m_person_grade pg on pg.id = p.id_person_grade
					left join 
						m_jenis_tujuan jt on jt.id = spd.jenis_tujuan
					left join 
						m_fasilitas_transport ft on ft.id = spd.id_fasilitas_transport
					left join 
						m_branch b on b.id = spd.id_branch `,
		SelectListDto: `SELECT * FROM fn_list_sppd($1, $2, $3, $4) as spd  (
			id varchar(36), tgl_surat date, nomor_surat varchar(100), id_pegawai varchar(36), jenis_tujuan varchar(36), 
			tujuan_dinas text, keperluan_dinas text, tenant_id varchar(36), jenis_sppd varchar(2), id_pegawai_pengaju varchar(36),
			tgl_berangkat date, tgl_kembali date, id_fasilitas_transport varchar(36), is_rombongan boolean, 
			id_branch varchar(36), created_at timestamp, is_antar boolean, is_jemput boolean, operasional_hari_dinas boolean,
			id_rule_approval varchar(36), id_jenis_approval varchar(36), is_pengajuan boolean, status varchar(2),
			nip varchar,  nama_pegawai varchar, id_jabatan varchar,id_person_grade varchar, jabatan varchar,  id_level_bod varchar, file text, link_file text,
			kode varchar, person_grade varchar,nama_jenis_tujuan varchar, keterangan text,  nama_fasilitas_transport varchar,
		   kode_branch varchar,nama_branch varchar, nama_approval varchar,id_pengajuan_bpd_histori varchar, is_deleted boolean, bidang varchar,
		  status_bpd varchar, type_approval varchar, id_pegawai_approval varchar, esign boolean, is_max_pengajuan boolean, is_max_penyelesaian boolean, nama_fungsionalitas varchar
	  )`,
		SelectDto2: `SELECT * FROM fn_list_sppd(?, ?, ?, ?) as spd  (
			id varchar(36), tgl_surat date, nomor_surat varchar(100), id_pegawai varchar(36), jenis_tujuan varchar(36), 
			tujuan_dinas text, keperluan_dinas text, tenant_id varchar(36), jenis_sppd varchar(2), id_pegawai_pengaju varchar(36),
			tgl_berangkat date, tgl_kembali date, id_fasilitas_transport varchar(36), is_rombongan boolean, 
			id_branch varchar(36), created_at timestamp, is_antar boolean, is_jemput boolean, operasional_hari_dinas boolean,
			id_rule_approval varchar(36), id_jenis_approval varchar(36), is_pengajuan boolean, status varchar(2),
			nip varchar,  nama_pegawai varchar, id_jabatan varchar,id_person_grade varchar, jabatan varchar, id_level_bod varchar, file text, link_file text,
			kode varchar, person_grade varchar,nama_jenis_tujuan varchar, keterangan text,  nama_fasilitas_transport varchar,
		   kode_branch varchar,nama_branch varchar, nama_approval varchar,id_pengajuan_bpd_histori varchar, is_deleted boolean, bidang varchar,
		  status_bpd varchar, type_approval varchar, id_pegawai_approval varchar, esign boolean, is_max_pengajuan boolean, is_max_penyelesaian boolean, nama_fungsionalitas varchar
	  )`,
		UpdateFileSppd:     `update surat_perjalanan_dinas set file=:file `,
		UpdateLinkFileSppd: `update surat_perjalanan_dinas set link_file=:link_file `,
	}
)

var suratPerjalananDinasPegawailQuery = struct {
	InsertBulk            string
	InsertBulkPlaceholder string
	SelectDto             string
}{
	InsertBulk: `INSERT INTO 
					surat_perjalanan_dinas_pegawai
						(id, id_surat_perjalanan_dinas, id_pegawai, status, created_at, created_by) 
					values `,
	InsertBulkPlaceholder: ` (:id, :id_surat_perjalanan_dinas, :id_pegawai, :status, :created_at, :created_by) `,
	SelectDto: `select 
					spdp.id, spdp.id_surat_perjalanan_dinas, spdp.id_pegawai,
					p.nip, p.nama as nama_pegawai, p.id_jabatan, p.id_person_grade, j.nama as jabatan, p.id_level_bod, mb.nama bidang
				from 
					surat_perjalanan_dinas_pegawai spdp
				left join 
					m_pegawai p on p.id = spdp.id_pegawai
				left join 
					m_jabatan j on j.id = p.id_jabatan
				left join 
					m_person_grade pg on pg.id = p.id_person_grade
					left join 
					m_bidang mb on mb.id = p.id_bidang `,
}

var sppdDokumensQuery = struct {
	InsertBulk            string
	InsertBulkPlaceholder string
	SelectDto             string
}{
	InsertBulk: `INSERT INTO 
						sppd_dokumen
							(id, id_sppd, file, id_dokumen, created_at, created_by) 
						values `,
	InsertBulkPlaceholder: ` (:id, :id_sppd, :file, :id_dokumen, :created_at, :created_by) `,
	SelectDto: `select 
					sd.id, sd.id_sppd, sd.file, sd.keterangan, sd.id_dokumen, md.nama, syd.is_mandatory
				from 
					sppd_Dokumen sd
				left join 
					m_syarat_dokumen syd on syd.id = sd.id_dokumen
					left join m_dokumen md on md.id = syd.id_dokumen
					 `,
}

var (
	detailBiayaQuery = struct {
		Select string
	}{
		Select: `select * from fn_get_generate_bpd_v2($1, $2) as (
					id_bpd varchar, id_sppd varchar, id_pegawai varchar, nama_pegawai varchar, tgl date, tgl_berangkat date,
					tgl_kembali date, id_komponen_biaya varchar, id_jenis_biaya varchar, biaya numeric, urut text,jenis_biaya text) `,
	}
)

type SuratPerjalananDinasRepository interface {
	Create(data SuratPerjalananDinas) error
	Update(data SuratPerjalananDinas) error
	UpdateSuratPerjalananDinas(data SuratPerjalananDinas) error
	GetNoSurat() (nomor string, err error)
	GetDetailRombonganByIdSurat(id string) (data []SuratPerjalananDinasPegawaiDto, err error)
	ResolveAllDto(req model.StandardRequest) (data pagination.Response, err error)
	ResolveDtoByID(id string) (data SuratPerjalananDinasDto, err error)
	ResolveByID(id string) (data SuratPerjalananDinas, err error)
	ResolveByDetailID(req FilterDetailSPPD) (data SuratPerjalananDinasListDto, err error)
	ResolveAllApproval(req model.StandardRequest) (data pagination.Response, err error)
	UpdateFileSuratPerjalananDinas(data FilesSuratPerjalananDinas) error
	GetDetailBiaya(tglAwal string, tglAkhir string, idSppd string, idBodLevel string) (data []DetailBiaya, err error)
	UpdateLinkFileSuratPerjalananDinas(data LinkFilesSuratPerjalananDinas) error
	GetSppdDokumenDto(id string) (data []SppdDokumenDto, err error)
}

type SuratPerjalananDinasRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideSuratPerjalananDinasRepositoryPostgreSQL(db *infras.PostgresqlConn) *SuratPerjalananDinasRepositoryPostgreSQL {
	s := new(SuratPerjalananDinasRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) Create(data SuratPerjalananDinas) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		// Function create table perjalanan_dinas
		if err := r.CreateTxSPDinas(tx, data); err != nil {
			e <- err
			return
		}

		// Function Insert Bulk table perjalanan_dinas_pegawai
		if err := txCreateSuratPerjalananDinasPegawai(tx, data.Detail); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) CreateTxSPDinas(tx *sqlx.Tx, data SuratPerjalananDinas) error {
	stmt, err := tx.PrepareNamed(suratPerjalananDinasQuery.Insert)
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

func txCreateSuratPerjalananDinasPegawai(tx *sqlx.Tx, details []SuratPerjalananDinasPegawai) (err error) {
	if len(details) == 0 {
		return
	}
	query, args, err := composeBulkUpsertSuratPerjalananDinasPegawaiDetailQuery(details)
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

func composeBulkUpsertSuratPerjalananDinasPegawaiDetailQuery(details []SuratPerjalananDinasPegawai) (qResult string, params []interface{}, err error) {
	values := []string{}
	for _, d := range details {
		param := map[string]interface{}{
			"id":                        d.ID,
			"id_surat_perjalanan_dinas": d.IdSuratPerjalananDinas,
			"id_pegawai":                d.IdPegawai,
			"status":                    d.Status,
			"created_at":                d.CreatedAt,
			"created_by":                d.CreatedBy,
		}
		q, args, err := sqlx.Named(suratPerjalananDinasPegawailQuery.InsertBulkPlaceholder, param)
		if err != nil {
			return qResult, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}
	qResult = fmt.Sprintf(`%v %v 
						ON CONFLICT (id) 
						DO UPDATE SET 
						id_pegawai=EXCLUDED.id_pegawai, 
						status=EXCLUDED.status, 
						created_at=EXCLUDED.created_at,
						created_by=EXCLUDED.created_by
						 `, suratPerjalananDinasPegawailQuery.InsertBulk, strings.Join(values, ","))
	return
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) UpdateSuratPerjalananDinas(data SuratPerjalananDinas) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		// Function update table perjalanan_dinas
		if err := r.UpdateTxSPDinas(tx, data); err != nil {
			e <- err
			return
		}

		// fmt.Print("leng detail : ", len(data.Detail))
		if len(data.Detail) > 0 {
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
			if err := txCreateSuratPerjalananDinasPegawai(tx, data.Detail); err != nil {
				e <- err
				return
			}

			// Function Insert Bulk table dokumen sppd
			if err := txCreateSppdDokumen(tx, data.DetailDokumen); err != nil {
				e <- err
				return
			}
		}
		e <- nil
	})
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) UpdateTxSPDinas(tx *sqlx.Tx, data SuratPerjalananDinas) error {
	stmt, err := tx.PrepareNamed(suratPerjalananDinasQuery.Update + " WHERE id=:id")
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

func (r *SuratPerjalananDinasRepositoryPostgreSQL) txDeleteDetailNotIn(tx *sqlx.Tx, idPDinas string, ids []string) (err error) {
	query, args, err := sqlx.In("update surat_perjalanan_dinas_pegawai set is_deleted=true where id_surat_perjalanan_dinas = ? AND id NOT IN (?)", idPDinas, ids)
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

func (r *SuratPerjalananDinasRepositoryPostgreSQL) GetNoSurat() (nomor string, err error) {
	err = r.DB.Read.Get(&nomor, suratPerjalananDinasQuery.SelectNoSurat)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) ResolveAllDto(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchParams = append(searchParams, req.IdPegawai)
	searchParams = append(searchParams, req.IdPegawaiApproval)
	searchParams = append(searchParams, req.TypeApproval)
	searchParams = append(searchParams, req.IdBidang)
	searchRoleBuff.WriteString(" WHERE coalesce(spd.is_deleted, false) = false ")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(spd.nomor_surat, spd.tujuan_dinas, spd.keperluan_dinas, spd.nama_pegawai, spd.nip, spd.jabatan, spd.person_grade, spd.nama_jenis_tujuan, spd.nama_fasilitas_transport, spd.kode_branch, spd.nama_branch) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" spd.id_branch = ?  ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.IdPegawaiApproval != "" {
		searchRoleBuff.WriteString(" and spd.id_pengajuan_bpd_histori is not null and spd.id_pegawai <> ? ")
		searchParams = append(searchParams, req.IdPegawaiApproval)
	}

	if req.Status != "" {
		searchRoleBuff.WriteString(" and spd.status = ? ")
		searchParams = append(searchParams, req.Status)
	}

	query := r.DB.Read.Rebind("select count(*) from (" + suratPerjalananDinasQuery.SelectDto2 + searchRoleBuff.String() + ")s")
	// query := r.DB.Read.Rebind(suratPerjalananDinasQuery.CountDto + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappSuratPerjalananDinasDto[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchSuratPerjalananDinasQuery := searchRoleBuff.String()
	searchSuratPerjalananDinasQuery = r.DB.Read.Rebind(suratPerjalananDinasQuery.SelectDto2 + searchSuratPerjalananDinasQuery)
	// fmt.Println("query", searchSuratPerjalananDinasQuery)
	rows, err := r.DB.Read.Queryx(searchSuratPerjalananDinasQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var item SuratPerjalananDinasDto
		err = rows.StructScan(&item)
		if err != nil {
			return
		}

		// detailBiaya, x := r.GetDetailBiaya(item.TglBerangkat, item.TglKembali, item.ID.String(), *item.IdLevelBod)
		// if x != nil {
		// 	return
		// }
		// item.DetailBiaya = detailBiaya

		if item.IsRombongan == true {
			detail, errs := r.GetDetailRombonganByIdSurat(item.ID.String())
			if errs != nil {
				return
			}
			item.Detail = detail
		}

		detailDokumen, errs := r.GetSppdDokumenDto(item.ID.String())
		if errs != nil {
			return
		}
		item.DetailDokumen = detailDokumen

		data.Items = append(data.Items, item)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) GetDetailRombonganByIdSurat(id string) (data []SuratPerjalananDinasPegawaiDto, err error) {
	err = r.DB.Read.Select(&data, suratPerjalananDinasPegawailQuery.SelectDto+" where spdp.id_surat_perjalanan_dinas = $1 and spdp.is_deleted = false ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) ResolveDtoByID(id string) (data SuratPerjalananDinasDto, err error) {
	err = r.DB.Read.Get(&data, suratPerjalananDinasQuery.SelectDto+" WHERE spd.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	if data.IsRombongan == true {
		detail, errs := r.GetDetailRombonganByIdSurat(data.ID.String())
		if errs != nil {
			return
		}
		data.Detail = detail
	}

	detailDokumen, errs := r.GetSppdDokumenDto(data.ID.String())
	if errs != nil {
		return
	}
	data.DetailDokumen = detailDokumen

	detailBiaya, x := r.GetDetailBiaya(data.TglBerangkat, data.TglKembali, data.ID.String(), *data.IdLevelBod)
	if x != nil {
		return
	}
	data.DetailBiaya = detailBiaya

	return
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) ResolveByID(id string) (data SuratPerjalananDinas, err error) {
	err = r.DB.Read.Get(&data, suratPerjalananDinasQuery.Select+" WHERE spd.id=$1 ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) Update(data SuratPerjalananDinas) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		// Function update table perjalanan_dinas
		if err := r.UpdateTxSPDinas(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txCreateSppdDokumen(tx *sqlx.Tx, details []SppdDokumen) (err error) {
	if len(details) == 0 {
		return
	}
	query, args, err := composeBulkUpsertSppdDokumenDetailQuery(details)
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

func composeBulkUpsertSppdDokumenDetailQuery(details []SppdDokumen) (qResult string, params []interface{}, err error) {
	values := []string{}
	for _, d := range details {
		param := map[string]interface{}{
			"id":         d.ID,
			"id_sppd":    d.IdSppd,
			"id_dokumen": d.IdDokumen,
			"created_at": d.CreatedAt,
			"created_by": d.CreatedBy,
		}
		q, args, err := sqlx.Named(sppdDokumensQuery.InsertBulkPlaceholder, param)
		if err != nil {
			return qResult, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}
	qResult = fmt.Sprintf(`%v %v 
						ON CONFLICT (id) 
						DO UPDATE SET 
						id_sppd=EXCLUDED.id_sppd, 
						created_at=EXCLUDED.created_at,
						created_by=EXCLUDED.created_by
						 `, sppdDokumensQuery.InsertBulk, strings.Join(values, ","))
	return
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) GetSppdDokumenDto(id string) (data []SppdDokumenDto, err error) {
	err = r.DB.Read.Select(&data, sppdDokumensQuery.SelectDto+" where sd.id_sppd = $1 ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) ResolveByDetailID(req FilterDetailSPPD) (data SuratPerjalananDinasListDto, err error) {
	err = r.DB.Read.Get(&data, suratPerjalananDinasQuery.SelectListDto+" WHERE spd.id=$5 ", req.IdPegawai, req.IdPegawaiApproval, req.TypeApproval, req.IdBidang, req.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	PDinasDetail, err := r.GetDetailRombonganByIdSurat(req.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	data.Detail = PDinasDetail

	detailDokumen, errs := r.GetSppdDokumenDto(req.ID)
	if errs != nil {
		return
	}
	data.DetailDokumen = detailDokumen
	return
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) ResolveAllApproval(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchParams = append(searchParams, req.IdPegawai)
	searchParams = append(searchParams, req.IdPegawaiApproval)
	searchParams = append(searchParams, req.TypeApproval)
	searchParams = append(searchParams, req.IdBidang)
	searchRoleBuff.WriteString(" WHERE coalesce(spd.is_deleted, false) = false ")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(spd.nomor_surat, spd.tujuan_dinas, spd.keperluan_dinas, spd.nama_pegawai, spd.nip, spd.jabatan, spd.person_grade, spd.nama_jenis_tujuan, spd.nama_fasilitas_transport, spd.kode_branch, spd.nama_branch) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" spd.id_branch = ?  ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.IdPegawaiApproval != "" {
		searchRoleBuff.WriteString(" and spd.id_pengajuan_bpd_histori is not null ")
	}

	if req.Status != "" {
		searchRoleBuff.WriteString(" and (case when spd.status='5' then spd.status else spd.status_bpd end) = ? ")
		searchParams = append(searchParams, req.Status)
	}

	if req.StartDate != "" {
		searchRoleBuff.WriteString(" and spd.tgl_berangkat= ? ")
		searchParams = append(searchParams, req.StartDate)
	}

	if req.EndDate != "" {
		searchRoleBuff.WriteString(" and spd.tgl_kembali= ? ")
		searchParams = append(searchParams, req.EndDate)
	}

	if req.IdTransaksi != "" {
		searchRoleBuff.WriteString(" and spd.jenis_sppd= ? ")
		searchParams = append(searchParams, req.IdTransaksi)
	}

	query := r.DB.Read.Rebind("select count(*) from (" + suratPerjalananDinasQuery.SelectDto2 + searchRoleBuff.String() + ")s")

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

	orders := ` ORDER BY (case when spd.status_bpd = '5' and spd.status = '5' then 4
			when spd.status_bpd = '1' then 1
			when spd.status_bpd = '3' then 2
			when spd.status_bpd = '2' then 3
			else 5 end
		),  `
	searchRoleBuff.WriteString(orders + ColumnMappSuratPerjalananDinasDto[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchBpdQuery := searchRoleBuff.String()
	searchBpdQuery = r.DB.Read.Rebind(suratPerjalananDinasQuery.SelectDto2 + searchBpdQuery)
	rows, err := r.DB.Read.Queryx(searchBpdQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var items SuratPerjalananDinasListDto
		err = rows.StructScan(&items)
		if err != nil {
			return
		}

		if items.IsRombongan == true {
			detail, errs := r.GetDetailRombonganByIdSurat(items.ID.String())
			if errs != nil {
				return
			}
			items.Detail = detail
		}

		data.Items = append(data.Items, items)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) GetDetailBiaya(tglAwal string, tglAkhir string, idSppd string, idBodLevel string) (data []DetailBiaya, err error) {
	rows, err := r.DB.Read.Queryx(detailBiayaQuery.Select, tglAwal, tglAkhir, idSppd, idBodLevel)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Detail Biaya NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items DetailBiaya
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	return
}

func (r *SuratPerjalananDinasRepositoryPostgreSQL) UpdateFileSuratPerjalananDinas(data FilesSuratPerjalananDinas) error {
	stmt, err := r.DB.Write.PrepareNamed(suratPerjalananDinasQuery.UpdateFileSppd + " WHERE id=:id ")
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

func (r *SuratPerjalananDinasRepositoryPostgreSQL) UpdateLinkFileSuratPerjalananDinas(data LinkFilesSuratPerjalananDinas) error {
	stmt, err := r.DB.Write.PrepareNamed(suratPerjalananDinasQuery.UpdateLinkFileSppd + " WHERE id=:id ")
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
