package master

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

var (
	settingBiayaQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `SELECT a.id, b.nama, c.nama as level_bod, d.nama as jenis_tujuan, a.id_jenis_biaya, a.created_at, a.created_by, a.updated_at, a.updated_by, a.is_deleted, a.is_harian,  a.tenant_id, a.id_branch, a.id_bod_level, a.jumlah_biaya, a.id_jenis_tujuan, e.id id_kategori_biaya, e.nama nama_kategori_biaya, a.is_max, a.jumlah_hari
		FROM m_komponen_biaya a
		left join m_jenis_biaya b on b.id = a.id_jenis_biaya
		left join m_kategori_biaya e on e.id = b.id_kategori_biaya
		left join m_level_bod c on c.id = a.id_bod_level
		left join m_jenis_tujuan d on d.id = a.id_jenis_tujuan
		`,
		Insert: `insert into m_komponen_biaya (id,  id_jenis_biaya, created_at, created_by, is_harian,  tenant_id, id_branch, id_bod_level, jumlah_biaya, id_jenis_tujuan, is_max, jumlah_hari)
				values (:id,  :id_jenis_biaya, :created_at, :created_by, :is_harian,  :tenant_id, :id_branch, :id_bod_level, :jumlah_biaya, :id_jenis_tujuan, :is_max, :jumlah_hari) `,
		Update: `update m_komponen_biaya set
					id=:id,
					id_jenis_biaya=:id_jenis_biaya,
					updated_at=:updated_at,
					updated_by=:updated_by,
					is_deleted=:is_deleted,
					tenant_id=:tenant_id,
					id_branch=:id_branch,
					is_harian=:is_harian,
					id_bod_level=:id_bod_level,
					jumlah_biaya=:jumlah_biaya,
					id_jenis_tujuan=:id_jenis_tujuan,
					is_max=:is_max,
					jumlah_hari=:jumlah_hari
				 `,
		Delete: `delete from m_komponen_biaya `,
		Count: `select count (a.id) from m_komponen_biaya a
		left join m_jenis_biaya b on b.id = a.id_jenis_biaya 
		left join m_kategori_biaya e on e.id = b.id_kategori_biaya
		left join m_level_bod c on c.id = a.id_bod_level
		left join m_jenis_tujuan d on d.id = a.id_jenis_tujuan
		`,
		Exist: `select count(id)>0 from m_komponen_biaya `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from perjalanan_dinas_biaya pd 
			where id_jenis_biaya = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type SettingBiayaRepository interface {
	Create(data SettingBiaya) error
	GetAll(req model.StandardRequest) (data []SettingBiaya, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data SettingBiaya, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data SettingBiaya) error
	ExistByNama(id string, idBranch string, idLevelBod string, idJenisTujuan string, idJenisBiaya string) (bool, error)
	GetAllKomponenBiaya(idSettingBiaya string) (data []KomponenBiaya, err error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type SettingBiayaRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideSettingBiayaRepositoryPostgreSQL(db *infras.PostgresqlConn) *SettingBiayaRepositoryPostgreSQL {
	s := new(SettingBiayaRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *SettingBiayaRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE a.is_deleted is false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND a.id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.IdJenisTujuan != "" {
		searchRoleBuff.WriteString(" AND a.id_jenis_tujuan = ? ")
		searchParams = append(searchParams, req.IdJenisTujuan)
	}

	if req.IdBodLevel != "" {
		searchRoleBuff.WriteString(" AND a.id_bod_level = ? ")
		searchParams = append(searchParams, req.IdBodLevel)
	}

	if req.IdTransaksi != "" {
		searchRoleBuff.WriteString(" AND b.id_kategori_biaya = ? ")
		searchParams = append(searchParams, req.IdTransaksi)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(b.nama, c.nama, d.nama, e.nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(settingBiayaQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappSettingBiaya[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchsettingBiayaQuery := searchRoleBuff.String()
	searchsettingBiayaQuery = r.DB.Read.Rebind(settingBiayaQuery.Select + searchsettingBiayaQuery)
	rows, err := r.DB.Read.Queryx(searchsettingBiayaQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var SettingBiaya SettingBiaya
		err = rows.StructScan(&SettingBiaya)
		if err != nil {
			return
		}

		data.Items = append(data.Items, SettingBiaya)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *SettingBiayaRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []SettingBiaya, err error) {
	where := " where coalesce(a.is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and a.id_branch='%s' ", req.IdBranch)
	}

	rows, err := r.DB.Read.Queryx(settingBiayaQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("SettingBiaya NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList SettingBiaya
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *SettingBiayaRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (SettingBiaya SettingBiaya, err error) {
	err = r.DB.Read.Get(&SettingBiaya, settingBiayaQuery.Select+" WHERE a.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *SettingBiayaRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(settingBiayaQuery.Delete+" WHERE a.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

// Function digunakan untuk create with transaction
func (r *SettingBiayaRepositoryPostgreSQL) Create(data SettingBiaya) error {
	stmt, err := r.DB.Write.PrepareNamed(settingBiayaQuery.Insert)
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

func (r *SettingBiayaRepositoryPostgreSQL) Update(data SettingBiaya) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdateSettingBiaya(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txUpdateSettingBiaya(tx *sqlx.Tx, data SettingBiaya) (err error) {
	stmt, err := tx.PrepareNamed(settingBiayaQuery.Update + " WHERE id=:id")
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

func (r *SettingBiayaRepositoryPostgreSQL) GetAllKomponenBiaya(idSettingBiaya string) (data []KomponenBiaya, err error) {
	where := " where a.is_deleted=false "
	if idSettingBiaya != "" {
		where += fmt.Sprintf(" and a.id_jenis_biaya = '%v' ", idSettingBiaya)
	}
	where += " order by a.created_at desc "

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

func (r *SettingBiayaRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, settingBiayaQuery.ExistRelasi, id)

	return
}

func (r *SettingBiayaRepositoryPostgreSQL) ExistByNama(id string, idBranch string, idLevelBod, idJenisTujuan string, idJenisBiaya string) (bool, error) {
	var exist bool

	criteria := ` where  coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	if idLevelBod != "" {
		criteria += fmt.Sprintf(" AND id_bod_level = '%s' ", idLevelBod)
	}

	if idJenisTujuan != "" {
		criteria += fmt.Sprintf(" and id_jenis_tujuan= '%s' ", idJenisTujuan)
	}

	if idJenisBiaya != "" {
		criteria += fmt.Sprintf(" and id_jenis_biaya= '%s' ", idJenisBiaya)
	}

	fmt.Println("query", settingBiayaQuery.Exist+criteria)
	err := r.DB.Read.Get(&exist, settingBiayaQuery.Exist+criteria)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}
