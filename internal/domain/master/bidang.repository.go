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
	bidangQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id, kode, nama, tenant_id, id_branch, created_at, created_by, updated_at, updated_by, is_deleted from m_bidang `,
		Insert: `insert into m_bidang
				(id, kode, nama, tenant_id, id_branch, created_at, created_by) 
				values 
				(:id, :kode, :nama, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_bidang set
				id=:id,
				kode=:kode,
				nama=:nama,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_bidang `,
		Count: `select count (id)
				from m_bidang `,
		Exist: `select count(id)>0 from m_bidang `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_unit_organisasi_kerja pd 
			where id_bidang = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x `,
	}
)

type BidangRepository interface {
	Create(data Bidang) error
	GetAll(req model.StandardRequest) (data []Bidang, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Bidang, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data Bidang) error
	ExistByNama(nama string, id string, idBranch string) (bool, error)
	ExistByKode(kode string, id string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type BidangRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideBidangRepositoryPostgreSQL(db *infras.PostgresqlConn) *BidangRepositoryPostgreSQL {
	s := new(BidangRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *BidangRepositoryPostgreSQL) Create(data Bidang) error {
	stmt, err := r.DB.Write.PrepareNamed(bidangQuery.Insert)
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

func (r *BidangRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(is_deleted, false) = false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(kode, nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(bidangQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappBidang[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchbidangQuery := searchRoleBuff.String()
	searchbidangQuery = r.DB.Read.Rebind(bidangQuery.Select + searchbidangQuery)
	fmt.Println("query", searchbidangQuery)
	rows, err := r.DB.Read.Queryx(searchbidangQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var bidang Bidang
		err = rows.StructScan(&bidang)
		if err != nil {
			return
		}

		data.Items = append(data.Items, bidang)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *BidangRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []Bidang, err error) {

	where := " where coalesce(is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and id_branch='%s' ", req.IdBranch)
	}
	where += "order by kode "

	rows, err := r.DB.Read.Queryx(bidangQuery.Select + where)
	// rows, err := r.DB.Read.Queryx(bidangQuery.Select + " where coalesce(is_deleted, false) = false order by kode")
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Bidang NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList Bidang
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *BidangRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data Bidang, err error) {
	err = r.DB.Read.Get(&data, bidangQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *BidangRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(bidangQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *BidangRepositoryPostgreSQL) Update(data Bidang) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *BidangRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data Bidang) (err error) {
	stmt, err := tx.PrepareNamed(bidangQuery.Update + " WHERE id=:id")
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

func (r *BidangRepositoryPostgreSQL) ExistByNama(nama string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, bidangQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *BidangRepositoryPostgreSQL) ExistByKode(kode string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(kode)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, bidangQuery.Exist+criteria, kode)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *BidangRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, bidangQuery.ExistRelasi, id)

	return
}
