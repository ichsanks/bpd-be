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
	kendaraanQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select mk.id, mk.nopol, mk.nama, mk.id_jenis_kendaraan, mjk.nama nama_jenis_kendaraan, mk.created_at, mk.created_by, mk.updated_at, mk.updated_by, mk.is_deleted
							from m_kendaraan mk
							left join m_jenis_kendaraan mjk on ( mk.id_jenis_kendaraan = mjk.id and mjk.is_deleted = false )`,
		Insert: `insert into m_kendaraan
				(id, nopol, nama, id_jenis_kendaraan, created_at, created_by)
				values
				(:id, :nopol, :nama, :id_jenis_kendaraan, :created_at, :created_by) `,
		Update: `update m_kendaraan set
				id=:id,
				nopol=:nopol,
				nama=:nama,
				id_jenis_kendaraan=:id_jenis_kendaraan,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_kendaraan `,
		Count: `select count (id)
				from m_kendaraan `,
		Exist: `select count(id)>0 from m_kendaraan`,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from perjalanan_dinas_kendaraan pd 
			where id_kendaraan = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type KendaraanRepository interface {
	Create(data Kendaraan) error
	GetAll() (data []Kendaraan, err error)
	ResolveAll(req model.StandardRequestKendaraan) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Kendaraan, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data Kendaraan) error
	ExistByNama(nama string, id string) (bool, error)
	ExistByNopol(nopol string, id string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type KendaraanRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideKendaraanRepositoryPostgreSQL(db *infras.PostgresqlConn) *KendaraanRepositoryPostgreSQL {
	s := new(KendaraanRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *KendaraanRepositoryPostgreSQL) Create(data Kendaraan) error {
	stmt, err := r.DB.Write.PrepareNamed(kendaraanQuery.Insert)
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

func (r *KendaraanRepositoryPostgreSQL) ResolveAll(req model.StandardRequestKendaraan) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(mk.is_deleted, false) = false ")

	if req.IdJenisKendaraan != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" mk.id_jenis_kendaraan = ? ")
		searchParams = append(searchParams, req.IdJenisKendaraan)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(mk.nopol, mk.nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind("select count(*) from (" + kendaraanQuery.Select + searchRoleBuff.String() + ")s")

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

	searchRoleBuff.WriteString("order by " + ColumnMappKendaraan[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchkendaraanQuery := searchRoleBuff.String()
	searchkendaraanQuery = r.DB.Read.Rebind(kendaraanQuery.Select + searchkendaraanQuery)
	fmt.Println("query", searchkendaraanQuery)
	rows, err := r.DB.Read.Queryx(searchkendaraanQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var kendaraan Kendaraan
		err = rows.StructScan(&kendaraan)
		if err != nil {
			return
		}

		data.Items = append(data.Items, kendaraan)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *KendaraanRepositoryPostgreSQL) GetAll() (data []Kendaraan, err error) {
	rows, err := r.DB.Read.Queryx(kendaraanQuery.Select + " where coalesce(mk.is_deleted, false) = false")
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Kendaraan NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList Kendaraan
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *KendaraanRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data Kendaraan, err error) {
	err = r.DB.Read.Get(&data, kendaraanQuery.Select+" WHERE mk.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *KendaraanRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(kendaraanQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *KendaraanRepositoryPostgreSQL) Update(data Kendaraan) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *KendaraanRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data Kendaraan) (err error) {
	stmt, err := tx.PrepareNamed(kendaraanQuery.Update + " WHERE id=:id")
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

func (r *KendaraanRepositoryPostgreSQL) ExistByNama(nama string, id string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	err := r.DB.Read.Get(&exist, kendaraanQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err

}

func (r *KendaraanRepositoryPostgreSQL) ExistByNopol(nopol string, id string) (bool, error) {
	var exist bool

	criteria := ` where upper(nopol)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	err := r.DB.Read.Get(&exist, kendaraanQuery.Exist+criteria, nopol)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *KendaraanRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, kendaraanQuery.ExistRelasi, id)

	return
}
