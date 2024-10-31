package master

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
	jenisBiayaQuery = struct {
		Select                string
		Insert                string
		Update                string
		Delete                string
		Count                 string
		Exist                 string
		ExistRelasi           string
		SelectJumlahBiaya     string
		SelectBiayaPenginapan string
		SelectAllDto          string
		SelectHeader          string
	}{
		Select: `select mjb.id, mjb.nama, mjb.created_at, mjb.created_by, mjb.updated_at, mjb.updated_by, mjb.is_deleted, mjb.is_multiple, mjb.urut, mjb.tenant_id, mjb.id_branch, mjb.id_kategori_biaya, mjb.kelompok_biaya from m_jenis_biaya mjb`,
		Insert: `insert into m_jenis_biaya (id, nama, created_at, created_by, is_multiple, urut, tenant_id, id_branch, id_kategori_biaya, kelompok_biaya)
				values (:id, :nama, :created_at, :created_by, :is_multiple, :urut, :tenant_id, :id_branch, :id_kategori_biaya, :kelompok_biaya) `,
		Update: `update m_jenis_biaya set
				id=:id,
				nama=:nama,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted,
				urut=:urut,
				is_multiple=:is_multiple,
				id_kategori_biaya=:id_kategori_biaya,
				kelompok_biaya=:kelompok_biaya,
				tenant_id=:tenant_id,
				id_branch=:id_branch
				 `,
		Delete: `delete from m_jenis_biaya mjb `,
		Count:  `select count (mjb.id) from m_jenis_biaya mjb `,
		Exist:  `select count(id)>0 from m_jenis_biaya `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from perjalanan_dinas_biaya pd 
			where id_jenis_biaya = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
		SelectJumlahBiaya: `SELECT kb.jumlah_biaya 
								FROM public.m_jenis_biaya jb
								LEFT JOIN public.m_komponen_biaya kb on jb.id = kb.id_jenis_biaya
								`,
		SelectBiayaPenginapan: `SELECT kb.jumlah_biaya 
								FROM public.m_jenis_biaya jb
								LEFT JOIN public.m_komponen_biaya kb on jb.id = kb.id_jenis_biaya
								`,
		SelectAllDto: `select
							mjb.id, mjb.nama, mkb.id as id_komponen_biaya, mkb.jumlah_biaya, mkb.jumlah_hari, mkb.is_max 
						from 
							m_jenis_biaya mjb
						left join 
							m_komponen_biaya mkb on mkb.id_jenis_biaya = mjb.id `,
		SelectHeader: `select id, nama, coalesce(id_kategori_biaya, 'REIMBURSEMENT') kategori from m_jenis_biaya
			where is_deleted=false
			and (id_kategori_biaya in ('JEMPUT', 'PERDIEM', 'TRANSPORT_LOKAL', 'ANTAR') or kelompok_biaya ilike '%REIMBURSEMENT%')
			order by urut, kategori `,
	}
)
var (
	komponenBiayaQuery = struct {
		Select                string
		InsertBulk            string
		InsertBulkPlaceholder string
	}{
		Select:                `select id, nama, id_jenis_biaya, urut, is_harian, created_at, created_by, updated_at, updated_by, is_deleted, tenant_id, id_branch from m_komponen_biaya `,
		InsertBulk:            `INSERT INTO public.m_komponen_biaya(id, nama, id_jenis_biaya, is_harian, urut, created_by, created_at, tenant_id, id_branch) values `,
		InsertBulkPlaceholder: ` (:id, :nama, :id_jenis_biaya, :is_harian, :urut, :created_by, :created_at, :tenant_id, :id_branch) `,
	}
)

type JenisBiayaRepository interface {
	Create(data JenisBiaya) error
	GetAll(req model.StandardRequest) (data []JenisBiaya, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data JenisBiaya, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data JenisBiaya) error
	UpdateBulk(data JenisBiaya) error
	ExistByNama(nama string, idBranch string) (bool, error)
	ExistByNamaID(id uuid.UUID, nama string, idBranch string) (bool, error)
	GetAllKomponenBiaya(idJenisBiaya string) (data []KomponenBiaya, err error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
	GetJumlahBiayaByIdBod(idBod uuid.UUID, ket string) (data JumlahBiaya, err error)
	GetAllDto(req model.StandardRequest) (data []JenisBiayaDto, err error)
	GetAllHeader() (data []JenisBiayaHeader, err error)
}

type JenisBiayaRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideJenisBiayaRepositoryPostgreSQL(db *infras.PostgresqlConn) *JenisBiayaRepositoryPostgreSQL {
	s := new(JenisBiayaRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *JenisBiayaRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE mjb.is_deleted is false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND mjb.id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.IdTransaksi != "" {
		searchRoleBuff.WriteString(" AND mjb.id_kategori_biaya = ? ")
		searchParams = append(searchParams, req.IdTransaksi)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(mjb.nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(jenisBiayaQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappJenisBiaya[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchjenisBiayaQuery := searchRoleBuff.String()
	searchjenisBiayaQuery = r.DB.Read.Rebind(jenisBiayaQuery.Select + searchjenisBiayaQuery)
	rows, err := r.DB.Read.Queryx(searchjenisBiayaQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var JenisBiaya JenisBiaya
		err = rows.StructScan(&JenisBiaya)
		if err != nil {
			return
		}

		data.Items = append(data.Items, JenisBiaya)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *JenisBiayaRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []JenisBiaya, err error) {
	where := " where coalesce(mjb.is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and mjb.id_branch='%s' ", req.IdBranch)
	}

	if req.IdTransaksi != "" {
		where += fmt.Sprintf(" and mjb.kelompok_biaya ilike '%%%s%%' ", req.IdTransaksi)
	}

	rows, err := r.DB.Read.Queryx(jenisBiayaQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("JenisBiaya NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList JenisBiaya
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *JenisBiayaRepositoryPostgreSQL) GetAllHeader() (data []JenisBiayaHeader, err error) {
	rows, err := r.DB.Read.Queryx(jenisBiayaQuery.SelectHeader)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("JenisBiaya NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList JenisBiayaHeader
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *JenisBiayaRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (JenisBiaya JenisBiaya, err error) {
	err = r.DB.Read.Get(&JenisBiaya, jenisBiayaQuery.Select+" WHERE mjb.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *JenisBiayaRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(jenisBiayaQuery.Delete+" WHERE mjb.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

// Function digunakan untuk create with transaction
func (r *JenisBiayaRepositoryPostgreSQL) Create(data JenisBiaya) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		// Function create table m_jenis_biaya
		if err := r.CreateTxJenisBiaya(tx, data); err != nil {
			e <- err
			return
		}

		// Function Insert Bulk table m_komponen_biaya
		if err := txCreateKomponenBiaya(tx, data.Detail); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *JenisBiayaRepositoryPostgreSQL) CreateTxJenisBiaya(tx *sqlx.Tx, data JenisBiaya) error {
	stmt, err := tx.PrepareNamed(jenisBiayaQuery.Insert)
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

func (r *JenisBiayaRepositoryPostgreSQL) Update(data JenisBiaya) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdateJenisBiaya(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *JenisBiayaRepositoryPostgreSQL) UpdateBulk(data JenisBiaya) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdateJenisBiaya(tx, data); err != nil {
			e <- err
			return
		}

		// Function delete not in table m_komponen_biaya
		ids := make([]string, 0)
		for _, d := range data.Detail {
			ids = append(ids, d.ID.String())
		}

		if err := r.txDeleteDetailNotIn(tx, data.ID.String(), ids); err != nil {
			e <- err
			return
		}

		// Function Insert Bulk table m_komponen_biaya
		if err := txCreateKomponenBiaya(tx, data.Detail); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

func txUpdateJenisBiaya(tx *sqlx.Tx, data JenisBiaya) (err error) {
	stmt, err := tx.PrepareNamed(jenisBiayaQuery.Update + " WHERE id=:id")
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (r *JenisBiayaRepositoryPostgreSQL) ExistByNama(nama string, idBranch string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, jenisBiayaQuery.Exist+" where upper(nama)=upper($1) and coalesce(is_deleted, false)=false and id_branch=$2 ", nama, idBranch)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *JenisBiayaRepositoryPostgreSQL) ExistByNamaID(id uuid.UUID, nama string, idBranch string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, jenisBiayaQuery.Exist+" where id <> $1 and upper(nama)=upper($2) and coalesce(is_deleted, false)=false and id_branch=$3 ", id, nama, idBranch)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func txCreateKomponenBiaya(tx *sqlx.Tx, details []KomponenBiaya) (err error) {
	if len(details) == 0 {
		return
	}
	query, args, err := composeBulkUpsertKomponenBiayaQuery(details)
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

func composeBulkUpsertKomponenBiayaQuery(details []KomponenBiaya) (qResult string, params []interface{}, err error) {
	values := []string{}
	for _, d := range details {
		param := map[string]interface{}{
			"id":             d.ID,
			"nama":           d.Nama,
			"id_jenis_biaya": d.IdJenisBiaya,
			"urut":           d.Urut,
			"is_harian":      d.IsHarian,
			"created_by":     d.CreatedBy,
			"created_at":     d.CreatedAt,
			"tenant_id":      d.TenantID,
			"id_branch":      d.IdBranch,
		}
		q, args, err := sqlx.Named(komponenBiayaQuery.InsertBulkPlaceholder, param)
		if err != nil {
			return qResult, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}
	qResult = fmt.Sprintf(`%v %v 
						ON CONFLICT (id) 
						DO UPDATE SET 
						nama=EXCLUDED.nama, 
						id_jenis_biaya=EXCLUDED.id_jenis_biaya, 
						urut=EXCLUDED.urut,
						tenant_id=EXCLUDED.tenant_id,
						id_branch=EXCLUDED.id_branch,
						is_harian=EXCLUDED.is_harian `, komponenBiayaQuery.InsertBulk, strings.Join(values, ","))
	return
}

func (r *JenisBiayaRepositoryPostgreSQL) txDeleteDetailNotIn(tx *sqlx.Tx, idJenisBiaya string, ids []string) (err error) {
	query, args, err := sqlx.In("update m_komponen_biaya set is_deleted=true where id_jenis_biaya = ? AND id NOT IN (?)", idJenisBiaya, ids)
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

func (r *JenisBiayaRepositoryPostgreSQL) GetAllKomponenBiaya(idJenisBiaya string) (data []KomponenBiaya, err error) {
	where := " where is_deleted=false "
	if idJenisBiaya != "" {
		where += fmt.Sprintf(" and id_jenis_biaya = '%v' ", idJenisBiaya)
	}
	where += " order by urut asc "

	rows, err := r.DB.Read.Queryx(komponenBiayaQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var items KomponenBiaya
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}
	return
}

func (r *JenisBiayaRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, jenisBiayaQuery.ExistRelasi, id)

	return
}

func (r *JenisBiayaRepositoryPostgreSQL) GetJumlahBiayaByIdBod(idBod uuid.UUID, ket string) (JenisBiaya JumlahBiaya, err error) {
	err = r.DB.Read.Get(&JenisBiaya, jenisBiayaQuery.SelectJumlahBiaya+" WHERE jb."+ket+"= true AND kb.id_bod_level=$1  ", idBod)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *JenisBiayaRepositoryPostgreSQL) GetAllDto(req model.StandardRequest) (data []JenisBiayaDto, err error) {
	where := " where coalesce(mjb.is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and mkb.id_branch='%s' ", req.IdBranch)
	}
	if req.IdBodLevel != "" {
		where += fmt.Sprintf(" and mkb.id_bod_level='%s' ", req.IdBodLevel)
	}
	if req.IdJenisTujuan != "" {
		where += fmt.Sprintf(" and mkb.id_jenis_tujuan='%s' ", req.IdJenisTujuan)
	}

	if req.IdTransaksi != "" {
		where += fmt.Sprintf(" and mjb.kelompok_biaya ilike '%%%s%%' ", req.IdTransaksi)
	}

	rows, err := r.DB.Read.Queryx(jenisBiayaQuery.SelectAllDto + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("JenisBiaya Dto NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList JenisBiayaDto
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}
