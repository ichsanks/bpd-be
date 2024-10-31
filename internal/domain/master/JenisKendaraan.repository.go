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
	jenisKendaraanQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select mjk.id, mjk.nama, mjk.pilih_kendaraan, mjk.created_at, mjk.created_by, mjk.updated_at, mjk.updated_by, mjk.is_deleted from m_jenis_kendaraan mjk `,
		Insert: `insert into m_jenis_kendaraan
				(id, nama, pilih_kendaraan, created_at, created_by)
				values
				(:id, :nama, :pilih_kendaraan, :created_at, :created_by) `,
		Update: `update m_jenis_kendaraan set
				id=:id,
				nama=:nama,
				pilih_kendaraan=:pilih_kendaraan,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_jenis_kendaraan mjk `,
		Count: `select count (mjk.id)
				from m_jenis_kendaraan mjk `,
		Exist: `select count(id)>0 from m_jenis_kendaraan `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from perjalanan_dinas pd 
			where id_jenis_kendaraan = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type JenisKendaraanRepository interface {
	Create(data JenisKendaraan) error
	GetAll() (data []JenisKendaraan, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data JenisKendaraan, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data JenisKendaraan) error
	ExistByNama(nama string) (bool, error)
	ExistByNamaID(id uuid.UUID, nama string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type JenisKendaraanRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideJenisKendaraanRepositoryPostgreSQL(db *infras.PostgresqlConn) *JenisKendaraanRepositoryPostgreSQL {
	s := new(JenisKendaraanRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *JenisKendaraanRepositoryPostgreSQL) Create(data JenisKendaraan) error {
	stmt, err := r.DB.Read.PrepareNamed(jenisKendaraanQuery.Insert)
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

func (r *JenisKendaraanRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE mjk.is_deleted is false ")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(mjk.nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(jenisKendaraanQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappJenisKendaraan[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchjenisKendaraanQuery := searchRoleBuff.String()
	searchjenisKendaraanQuery = r.DB.Read.Rebind(jenisKendaraanQuery.Select + searchjenisKendaraanQuery)
	fmt.Println("query", searchjenisKendaraanQuery)
	rows, err := r.DB.Read.Queryx(searchjenisKendaraanQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var JenisKendaraan JenisKendaraan
		err = rows.StructScan(&JenisKendaraan)
		if err != nil {
			return
		}

		data.Items = append(data.Items, JenisKendaraan)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *JenisKendaraanRepositoryPostgreSQL) GetAll() (data []JenisKendaraan, err error) {
	rows, err := r.DB.Read.Queryx(jenisKendaraanQuery.Select + " where mjk.is_deleted is false")
	if err == sql.ErrNoRows {
		_ = failure.NotFound("JenisKendaraan NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList JenisKendaraan
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *JenisKendaraanRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (JenisKendaraan JenisKendaraan, err error) {
	err = r.DB.Read.Get(&JenisKendaraan, jenisKendaraanQuery.Select+" WHERE mjk.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *JenisKendaraanRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(jenisKendaraanQuery.Delete+" WHERE mjk.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *JenisKendaraanRepositoryPostgreSQL) Update(data JenisKendaraan) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdateJenisKendaraan(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txUpdateJenisKendaraan(tx *sqlx.Tx, data JenisKendaraan) (err error) {
	stmt, err := tx.PrepareNamed(jenisKendaraanQuery.Update + " WHERE id=:id")
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

func (r *JenisKendaraanRepositoryPostgreSQL) ExistByNama(nama string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, jenisKendaraanQuery.Exist+" where upper(nama)=upper($1) and coalesce(is_deleted, false)=false ", nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *JenisKendaraanRepositoryPostgreSQL) ExistByNamaID(id uuid.UUID, nama string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, jenisKendaraanQuery.Exist+" where id <> $1 and upper(nama)=upper($2) and coalesce(is_deleted, false)=false ", id, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *JenisKendaraanRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, jenisKendaraanQuery.ExistRelasi, id)

	return
}
