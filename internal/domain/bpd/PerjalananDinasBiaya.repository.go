package bpd

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
)

var (
	perjalananDinasBiayaQuery = struct {
		Select                       string
		SelectJenisBiaya             string
		SelectKomponenBiaya          string
		SelectKomponenBiayaDTO       string
		SelectKomponenBiayaUmDTO     string
		InsertBulk                   string
		InsertBulkPlaceholder        string
		InsertBulkUm                 string
		InsertBulkPlaceholderUm      string
		UpdateFileBpd                string
		UpdateUangMuka               string
		Insert                       string
		InsertAk                     string
		SelectDto                    string
		SelectDtoHistori             string
		SelectDtoHistoriIdPerjalanan string
		DeleteBiayaById              string
		Update                       string
		UpdateAk                     string
		ExistAk                      string
	}{
		Select: `SELECT id, id_bpd_pegawai, id_jenis_biaya, id_komponen_biaya, keterangan, to_char(tanggal,'YYYY-MM-DD') tanggal, qty, nominal, created_at, created_by, updated_at, updated_by, is_deleted
		FROM public.perjalanan_dinas_biaya`,
		SelectJenisBiaya: `select b.id id_jenis_biaya, b.nama nama_jenis_biaya, b.is_multiple from m_jenis_biaya b
			where coalesce(b.is_deleted, false)=false
			order by b.urut asc
		`,
		SelectKomponenBiayaDTO: `select p.id, p.id_bpd_pegawai, to_char(p.tanggal,'YYYY-MM-DD') tanggal, b.id id_komponen_biaya, b.id_jenis_biaya, b.nama, b.is_harian, p.keterangan, p.qty, p.nominal 
			from perjalanan_dinas_biaya p
			left join m_komponen_biaya b on b.id = p.id_komponen_biaya
			where coalesce(p.is_deleted, false)=false
			and p.id_bpd_pegawai = $1
			and b.id_jenis_biaya = $2
			order by b.urut asc `,
		SelectKomponenBiayaUmDTO: `select p.id, p.id_bpd_pegawai, b.id id_komponen_biaya, b.id_jenis_biaya, b.nama, b.is_harian, p.keterangan, p.qty, p.nominal 
			from perjalanan_dinas_um p
			left join m_komponen_biaya b on b.id = p.id_komponen_biaya
			where coalesce(p.is_deleted, false)=false
			and p.id_bpd_pegawai = $1
			and b.id_jenis_biaya = $2
			order by b.urut asc `,
		SelectKomponenBiaya: `select b.id id_komponen_biaya, b.id_jenis_biaya, b.nama, b.is_harian from m_komponen_biaya b
			where coalesce(b.is_deleted, false)=false
			and b.id_jenis_biaya = $1
			order by b.urut asc `,
		InsertBulk:              `INSERT INTO public.perjalanan_dinas_biaya(id, id_bpd_pegawai, id_jenis_biaya, id_komponen_biaya, keterangan, qty, nominal, created_at, created_by, id_pegawai, tanggal) values `,
		InsertBulkPlaceholder:   ` (:id, :id_bpd_pegawai, :id_jenis_biaya, :id_komponen_biaya, :keterangan, :qty, :nominal, :created_at, :created_by, :id_pegawai, :tanggal) `,
		InsertBulkUm:            `INSERT INTO public.perjalanan_dinas_um(id, id_bpd_pegawai, id_jenis_biaya, id_komponen_biaya, keterangan, qty, nominal, created_at, created_by) values `,
		InsertBulkPlaceholderUm: ` (:id, :id_bpd_pegawai, :id_jenis_biaya, :id_komponen_biaya, :keterangan, :qty, :nominal, :created_at, :created_by) `,
		UpdateFileBpd: `update perjalanan_dinas_pegawai set 
						file=:file,
						is_revisi=:is_revisi `,
		UpdateUangMuka: `update perjalanan_dinas_pegawai set 
						is_um=:is_um,
						persentase_um=:persentase_um,
						persentase_sisa=:persentase_sisa,
						show_um=:show_um,
						show_sisa=:show_sisa,
						total_um=:total_um,
						sisa_um=:sisa_um `,
		Insert: `insert into perjalanan_dinas_biaya
					(id, id_bpd_pegawai, id_jenis_biaya, id_komponen_biaya, nominal, 
					file, created_at, created_by, is_reimbursement, keterangan, id_pegawai, tanggal, aksi) 
				values 
					(:id, :id_bpd_pegawai, :id_jenis_biaya, :id_komponen_biaya, :nominal,
				 	:file, :created_at, :created_by, :is_reimbursement, :keterangan, :id_pegawai, :tanggal, :aksi) `,
		InsertAk: `insert into perjalanan_dinas_biaya
					 (id, id_bpd_pegawai, id_jenis_biaya, id_komponen_biaya, nominal, 
					 file, created_at, created_by, is_reimbursement, keterangan, id_pegawai, aksi) 
				 values 
					 (:id, :id_bpd_pegawai, :id_jenis_biaya, :id_komponen_biaya, :nominal,
					  :file, :created_at, :created_by, :is_reimbursement, :keterangan, :id_pegawai, :aksi) `,
		SelectDto: `select 
						pdb.id, pdb.id_bpd_pegawai, pdb.id_jenis_biaya, pdb.id_komponen_biaya, pdb.nominal, pdb.id_pegawai, pdb.keterangan, to_char(pdb.tanggal,'YYYY-MM-DD') tanggal,
						pdb.file, mjb.nama
					from 
						perjalanan_dinas_biaya pdb 
					left join 
						m_jenis_biaya mjb on mjb.id = pdb.id_jenis_biaya `,
		SelectDtoHistori: `select 
						pdb.id, pdb.id_bpd_pegawai, pdb.id_jenis_biaya, pdb.id_komponen_biaya, pdb.nominal, pdb.id_pegawai, pdb.keterangan, to_char(pdb.tanggal,'YYYY-MM-DD') tanggal,
						pdb.file, mjb.nama, pdb.created_at, pdb.created_by, pdb.updated_at, pdb.updated_by, pdb.deleted_at, pdb.deleted_by, pdb.aksi, pdb.is_reimbursement,
						(SELECT c.nama FROM auth_user c WHERE c.id = pdb.created_by) AS user_create,
						(SELECT u.nama FROM auth_user u WHERE u.id = pdb.updated_by) AS user_update,
						(SELECT d.nama FROM auth_user d WHERE d.id = pdb.deleted_by) AS user_delete
					from 
						log_biaya pdb 
					left join 
						m_jenis_biaya mjb on mjb.id = pdb.id_jenis_biaya `,
		SelectDtoHistoriIdPerjalanan: `select 
						pdb.id, pdb.id_bpd_pegawai, pdb.id_jenis_biaya, pdb.id_komponen_biaya, pdb.nominal, pdb.id_pegawai, pdb.keterangan, to_char(pdb.tanggal,'YYYY-MM-DD') tanggal,
						pdb.file, mjb.nama, pdb.created_at, pdb.created_by, pdb.updated_at, pdb.updated_by, pdb.deleted_at, pdb.deleted_by, pdb.aksi, pdb.is_reimbursement,
						(SELECT c.nama FROM auth_user c WHERE c.id = pdb.created_by) AS user_create,
						(SELECT u.nama FROM auth_user u WHERE u.id = pdb.updated_by) AS user_update,
						(SELECT d.nama FROM auth_user d WHERE d.id = pdb.deleted_by) AS user_delete
					from log_biaya pdb 
					left join perjalanan_dinas_pegawai pdp on pdp.id = pdb.id_bpd_pegawai
					left join m_jenis_biaya mjb on mjb.id = pdb.id_jenis_biaya `,
		DeleteBiayaById: `delete from perjalanan_dinas_biaya `,
		Update: `UPDATE perjalanan_dinas_biaya SET
					id=:id,
					id_bpd_pegawai=:id_bpd_pegawai,
					id_jenis_biaya=:id_jenis_biaya,
					nominal=:nominal,
					id_pegawai=:id_pegawai,
					id_komponen_biaya=:id_komponen_biaya,  
					file=:file,
					keterangan=:keterangan,
					tanggal=:tanggal,
					is_reimbursement=:is_reimbursement,
					updated_at=:updated_at,
					updated_by=:updated_by,
					aksi=:aksi,
					is_deleted=:is_deleted,
					deleted_at=:deleted_at,
					deleted_by=:deleted_by `,
		UpdateAk: `UPDATE perjalanan_dinas_biaya SET
					id=:id,
					id_bpd_pegawai=:id_bpd_pegawai,
					id_jenis_biaya=:id_jenis_biaya,
					nominal=:nominal,
					id_pegawai=:id_pegawai,
					id_komponen_biaya=:id_komponen_biaya,  
					file=:file,
					keterangan=:keterangan,
					is_reimbursement=:is_reimbursement,
					aksi=:aksi,
					updated_at=:updated_at,
					updated_by=:updated_by `,
		ExistAk: `select count(id)>0 from perjalanan_dinas_biaya`,
	}
)

type PerjalananDinasBiayaRepository interface {
	CreateBulk(data []PerjalananDinasBiaya) error
	CreateBulkUm(data []PerjalananDinasBiaya) error
	GetAllData(idBpdPegawai string) (data []PerjalananDinasBiayaDTO, err error)
	GetAllDataUm(idBpdPegawai string) (data []PerjalananDinasBiayaDTO, err error)
	UploadDocPenyelesaianBpd(data DocPenyelesaianBpdPegawai) error
	UpdateUangMuka(data UangMukaBpd) error
	Create(data PerjalananDinasBiaya, jenis string) error
	GetBiayaDto(idBpdPegawai string, idPegawai string, isReimbursement string) (data []BiayaPerjalananDinasDto, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data PerjalananDinasBiaya, jenis string) error
	ResolveByID(id uuid.UUID) (data BiayaPerjalananDinasDto, err error)
	ExistAkomodasi(idJenisBiaya string, idBpdPegawai string, id string) (bool, error)
	GetHistoriBiaya(idBpdPegawai string, idPegawai string, idJenisBiaya string, isReimbursement string) (data []HistoriPerjalananDinas, err error)
	ResolveByIdBiaya(id uuid.UUID) (data PerjalananDinasBiaya, err error)
	GetHistoriBiayaIdPerjalanan(idPerjalananDinas string) (data []HistoriPerjalananDinas, err error)
}

type PerjalananDinasBiayaRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvidePerjalananDinasBiayaRepositoryPostgreSQL(db *infras.PostgresqlConn) *PerjalananDinasBiayaRepositoryPostgreSQL {
	s := new(PerjalananDinasBiayaRepositoryPostgreSQL)
	s.DB = db
	return s
}

// Function digunakan untuk create with transaction
func (r *PerjalananDinasBiayaRepositoryPostgreSQL) CreateBulk(data []PerjalananDinasBiaya) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		ids := make([]string, 0)
		for _, d := range data {
			ids = append(ids, d.ID.String())
		}

		// if len(ids) > 0 {
		// if err := r.txDeleteDetailNotIn(tx, data[0].IDBpdPegawai, ids); err != nil {
		// 	if err := r.txDeleteDetailNotIn(tx, ids); err != nil {
		// 		e <- err
		// 		return
		// 	}
		// }

		if err := txCreatePerjalananDinasBiaya(tx, data, "BIAYA"); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

// Function digunakan untuk create with transaction
func (r *PerjalananDinasBiayaRepositoryPostgreSQL) CreateBulkUm(data []PerjalananDinasBiaya) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		ids := make([]string, 0)
		for _, d := range data {
			ids = append(ids, d.ID.String())
		}

		if len(ids) > 0 {
			if err := r.txDeleteDetailNotInUm(tx, data[0].IDBpdPegawai, ids); err != nil {
				e <- err
				return
			}
		}

		if err := txCreatePerjalananDinasBiaya(tx, data, "UM"); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txCreatePerjalananDinasBiaya(tx *sqlx.Tx, details []PerjalananDinasBiaya, jenis string) (err error) {
	if len(details) == 0 {
		return
	}
	query, args, err := composeBulkUpsertPerjalananDinasBiayaQuery(details, jenis)
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

func composeBulkUpsertPerjalananDinasBiayaQuery(details []PerjalananDinasBiaya, jenis string) (qResult string, params []interface{}, err error) {
	values := []string{}
	var qInsert, qInsertPl string
	if jenis == "BIAYA" {
		qInsert = perjalananDinasBiayaQuery.InsertBulk
		qInsertPl = perjalananDinasBiayaQuery.InsertBulkPlaceholder
	} else if jenis == "UM" {
		qInsert = perjalananDinasBiayaQuery.InsertBulkUm
		qInsertPl = perjalananDinasBiayaQuery.InsertBulkPlaceholderUm
	}

	for _, d := range details {
		param := map[string]interface{}{
			"id":                d.ID,
			"id_bpd_pegawai":    d.IDBpdPegawai,
			"id_jenis_biaya":    d.IDJenisBiaya,
			"id_komponen_biaya": d.IDKomponenBiaya,
			"keterangan":        d.Keterangan,
			"id_pegawai":        d.IdPegawai,
			"tanggal":           d.Tanggal,
			"qty":               d.Qty,
			"nominal":           d.Nominal,
			"created_at":        d.CreatedAt,
			"created_by":        d.CreatedBy,
		}

		q, args, err := sqlx.Named(qInsertPl, param)
		if err != nil {
			return qResult, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}
	qResult = fmt.Sprintf(`%v %v 
						ON CONFLICT (id) 
						DO UPDATE SET 
						id_jenis_biaya=EXCLUDED.id_jenis_biaya, 
						id_komponen_biaya=EXCLUDED.id_komponen_biaya, 
						keterangan=EXCLUDED.keterangan, 
						id_pegawai=EXCLUDED.id_pegawai, 
						qty=EXCLUDED.qty, 
						nominal=EXCLUDED.nominal `, qInsert, strings.Join(values, ","))
	return
}

// func (r *PerjalananDinasBiayaRepositoryPostgreSQL) txDeleteDetailNotIn(tx *sqlx.Tx, idBpdPegawai string, ids []string) (err error) {
func (r *PerjalananDinasBiayaRepositoryPostgreSQL) txDeleteDetailNotIn(tx *sqlx.Tx, ids []string) (err error) {
	query, args, err := sqlx.In("update perjalanan_dinas_biaya set is_deleted=true where  id NOT IN (?)", ids)
	// query, args, err := sqlx.In("update perjalanan_dinas_biaya set is_deleted=true where id_bpd_pegawai = ? AND id NOT IN (?)", idBpdPegawai, ids)
	query = tx.Rebind(query)

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	res, err := r.DB.Write.Exec(query, args...)
	if err != nil {
		return
	}
	_, err = res.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) txDeleteDetailNotInUm(tx *sqlx.Tx, idBpdPegawai string, ids []string) (err error) {
	query, args, err := sqlx.In("update perjalanan_dinas_um set is_deleted=true where id_bpd_pegawai = ? AND id NOT IN (?)", idBpdPegawai, ids)
	query = tx.Rebind(query)

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	res, err := r.DB.Write.Exec(query, args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = res.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) GetAllData(idBpdPegawai string) (data []PerjalananDinasBiayaDTO, err error) {
	rows, err := r.DB.Read.Queryx(perjalananDinasBiayaQuery.SelectJenisBiaya)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items PerjalananDinasBiayaDTO
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		listKomponenBiaya, errr := r.GetAllKomponenBiaya(items.IDJenisBiaya)
		if errr != nil {
			fmt.Println(errr)
		}
		items.ListKomponenBiaya = listKomponenBiaya

		detailBiaya, errr := r.GetAllKomponenBiayaDTO(idBpdPegawai, items.IDJenisBiaya, "BIAYA")
		if errr != nil {
			fmt.Println(errr)
		}
		items.DetailBiaya = detailBiaya

		data = append(data, items)
	}
	return
}

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) GetAllDataUm(idBpdPegawai string) (data []PerjalananDinasBiayaDTO, err error) {
	rows, err := r.DB.Read.Queryx(perjalananDinasBiayaQuery.SelectJenisBiaya)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items PerjalananDinasBiayaDTO
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		listKomponenBiaya, errr := r.GetAllKomponenBiaya(items.IDJenisBiaya)
		if errr != nil {
			fmt.Println(errr)
		}
		items.ListKomponenBiaya = listKomponenBiaya

		detailBiaya, errr := r.GetAllKomponenBiayaDTO(idBpdPegawai, items.IDJenisBiaya, "UM")
		if errr != nil {
			fmt.Println(errr)
		}
		items.DetailBiaya = detailBiaya

		data = append(data, items)
	}
	return
}

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) GetAllKomponenBiaya(idJenisBiaya string) (data []KomponenBiayaDetailDTO, err error) {
	rows, err := r.DB.Read.Queryx(perjalananDinasBiayaQuery.SelectKomponenBiaya, idJenisBiaya)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var items KomponenBiayaDetailDTO
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}
	return
}

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) GetAllKomponenBiayaDTO(idBpdPegawai string, idJenisBiaya string, jenis string) (data []KomponenBiayaDetailDTO, err error) {
	var rawQuery string
	if jenis == "BIAYA" {
		rawQuery = perjalananDinasBiayaQuery.SelectKomponenBiayaDTO
	} else if jenis == "UM" {
		rawQuery = perjalananDinasBiayaQuery.SelectKomponenBiayaUmDTO
	}
	rows, err := r.DB.Read.Queryx(rawQuery, idBpdPegawai, idJenisBiaya)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var items KomponenBiayaDetailDTO
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}
	return
}

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) UploadDocPenyelesaianBpd(data DocPenyelesaianBpdPegawai) error {
	stmt, err := r.DB.Write.PrepareNamed(perjalananDinasBiayaQuery.UpdateFileBpd + " WHERE id=:id ")
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

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) UpdateUangMuka(data UangMukaBpd) error {
	stmt, err := r.DB.Write.PrepareNamed(perjalananDinasBiayaQuery.UpdateUangMuka + " WHERE id=:id ")
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

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) Create(data PerjalananDinasBiaya, jenis string) error {
	var query string
	if jenis == "REM" {
		query = perjalananDinasBiayaQuery.Insert
	} else if jenis == "AK" {
		query = perjalananDinasBiayaQuery.InsertAk
	}

	stmt, err := r.DB.Write.PrepareNamed(query)
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

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) GetBiayaDto(idBpdPegawai string, idPegawai string, isReimbursement string) (data []BiayaPerjalananDinasDto, err error) {
	rows, err := r.DB.Read.Queryx(perjalananDinasBiayaQuery.SelectDto+"where pdb.id_bpd_pegawai = $1 and pdb.id_pegawai = $2 and pdb.is_reimbursement = $3 ", idBpdPegawai, idPegawai, isReimbursement)
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

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(perjalananDinasBiayaQuery.DeleteBiayaById+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) Update(data PerjalananDinasBiaya, jenis string) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		// Function update table perjalanan_dinas
		if err := r.UpdateTxDinasBiaya(tx, data, jenis); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) UpdateTxDinasBiaya(tx *sqlx.Tx, data PerjalananDinasBiaya, jenis string) error {
	var query string
	if jenis == "REM" {
		query = perjalananDinasBiayaQuery.Update
	} else if jenis == "AK" {
		query = perjalananDinasBiayaQuery.UpdateAk
	}
	stmt, err := tx.PrepareNamed(query + " WHERE id=:id")
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

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data BiayaPerjalananDinasDto, err error) {
	err = r.DB.Read.Get(&data, perjalananDinasBiayaQuery.SelectDto+" WHERE pdb.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) ExistAkomodasi(IDJenisBiaya string, idBpdPegawai string, id string) (bool, error) {
	var exist bool

	criteria := ` where id_jenis_biaya=$1 and id_bpd_pegawai=$2  and coalesce(is_deleted, false)=false and coalesce(is_reimbursement, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	err := r.DB.Read.Get(&exist, perjalananDinasBiayaQuery.ExistAk+criteria, IDJenisBiaya, idBpdPegawai)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) GetHistoriBiaya(idBpdPegawai string, idPegawai string, idJenisBiaya string, isReimbursement string) (data []HistoriPerjalananDinas, err error) {
	rows, err := r.DB.Read.Queryx(perjalananDinasBiayaQuery.SelectDtoHistori+" where pdb.id_bpd_pegawai = $1 and pdb.id_pegawai = $2 and pdb.id_jenis_biaya = $3 and pdb.is_reimbursement = $4 ORDER BY pdb.is_reimbursement ASC", idBpdPegawai, idPegawai, idJenisBiaya, isReimbursement)
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

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) ResolveByIdBiaya(id uuid.UUID) (data PerjalananDinasBiaya, err error) {
	err = r.DB.Read.Get(&data, perjalananDinasBiayaQuery.Select+" WHERE id=$1 and is_deleted=false  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *PerjalananDinasBiayaRepositoryPostgreSQL) GetHistoriBiayaIdPerjalanan(idPerjalananDinas string) (data []HistoriPerjalananDinas, err error) {
	rows, err := r.DB.Read.Queryx(perjalananDinasBiayaQuery.SelectDtoHistoriIdPerjalanan+" where pdp.id_perjalanan_dinas = $1 ORDER BY pdb.is_reimbursement ASC ", idPerjalananDinas)
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
