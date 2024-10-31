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
	fasilitasTransportQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id, nama, tenant_id, id_branch, created_at, created_by, updated_at, updated_by, is_deleted from m_fasilitas_transport `,
		Insert: `insert into m_fasilitas_transport
				(id, nama, tenant_id, id_branch, created_at, created_by)
				values
				(:id, :nama, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_fasilitas_transport set
				id=:id,
				nama=:nama,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_fasilitas_transport `,
		Count: `select count (id)
				from m_fasilitas_transport `,
		Exist: `select count(id)>0 from m_fasilitas_transport `,
	}
)

type FasilitasTransportRepository interface {
	Create(data FasilitasTransport) error
	GetAll(req model.StandardRequest) (data []FasilitasTransport, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data FasilitasTransport, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data FasilitasTransport) error
	ExistByNama(nama string, id string, idBranch string) (bool, error)
}

type FasilitasTransportRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideFasilitasTransportRepositoryPostgreSQL(db *infras.PostgresqlConn) *FasilitasTransportRepositoryPostgreSQL {
	s := new(FasilitasTransportRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *FasilitasTransportRepositoryPostgreSQL) Create(data FasilitasTransport) error {
	stmt, err := r.DB.Write.PrepareNamed(fasilitasTransportQuery.Insert)
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

func (r *FasilitasTransportRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
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

	query := r.DB.Read.Rebind(fasilitasTransportQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappFasilitasTransport[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchfasilitasTransportQuery := searchRoleBuff.String()
	searchfasilitasTransportQuery = r.DB.Read.Rebind(fasilitasTransportQuery.Select + searchfasilitasTransportQuery)
	fmt.Println("query", searchfasilitasTransportQuery)
	rows, err := r.DB.Read.Queryx(searchfasilitasTransportQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var FasilitasTransport FasilitasTransport
		err = rows.StructScan(&FasilitasTransport)
		if err != nil {
			return
		}

		data.Items = append(data.Items, FasilitasTransport)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *FasilitasTransportRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []FasilitasTransport, err error) {
	where := " where coalesce(is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and id_branch='%s' ", req.IdBranch)
	}
	rows, err := r.DB.Read.Queryx(fasilitasTransportQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("FasilitasTransport NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList FasilitasTransport
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *FasilitasTransportRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data FasilitasTransport, err error) {
	err = r.DB.Read.Get(&data, fasilitasTransportQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *FasilitasTransportRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(fasilitasTransportQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *FasilitasTransportRepositoryPostgreSQL) Update(data FasilitasTransport) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *FasilitasTransportRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data FasilitasTransport) (err error) {
	stmt, err := tx.PrepareNamed(fasilitasTransportQuery.Update + " WHERE id=:id")
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

func (r *FasilitasTransportRepositoryPostgreSQL) ExistByNama(nama string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, fasilitasTransportQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}
