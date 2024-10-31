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
	sTtdQuery = struct {
		Select      string
		SelectDto   string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id, id_pegawai, id_jabatan, jenis,  tenant_id, id_branch, created_at, created_by, updated_at, updated_by, is_deleted from s_ttd  `,
		SelectDto: `select a.id, a.id_pegawai, a.id_jabatan, a.jenis,  a.tenant_id, a.id_branch, a.created_at, a.created_by, a.updated_at, a.updated_by, a.is_deleted, b.nama nama_pegawai, b.nip, c.nama nama_jabatan,
		(CASE WHEN a.jenis='1' THEN 'SPPD'
		WHEN a.jenis='2' THEN 'BPD'
		ELSE '' END) AS nama_jenis from s_ttd a
		left join m_pegawai b on b.id = a.id_pegawai
		left join m_jabatan c on c.id = a.id_jabatan `,
		Insert: `insert into s_ttd
				(id, id_pegawai, id_jabatan, jenis, tenant_id, id_branch, created_at, created_by) 
				values 
				(:id, :id_pegawai, :id_jabatan, :jenis, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update s_ttd set
				id=:id,
				id_pegawai=:id_pegawai,
				id_jabatan=:id_jabatan,
				jenis=:jenis,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from s_ttd `,
		Count: `select count (a.id)
	    from s_ttd a
		left join m_pegawai b on b.id = a.id_pegawai
		left join m_fungsionalitas c on c.id = a.id_jabatan `,
	}
)

type STtdRepository interface {
	Create(data STtd) error
	GetAll(req model.StandardRequest) (data []STtdDto, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data STtd, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data STtd) error
}

type STtdRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideSTtdRepositoryPostgreSQL(db *infras.PostgresqlConn) *STtdRepositoryPostgreSQL {
	s := new(STtdRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *STtdRepositoryPostgreSQL) Create(data STtd) error {
	stmt, err := r.DB.Write.PrepareNamed(sTtdQuery.Insert)
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

func (r *STtdRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(a.is_deleted, false) = false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND a.id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(b.nama, c.nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(sTtdQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappSTtd[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchsTtdQuery := searchRoleBuff.String()
	searchsTtdQuery = r.DB.Read.Rebind(sTtdQuery.SelectDto + searchsTtdQuery)
	fmt.Println("query", searchsTtdQuery)
	rows, err := r.DB.Read.Queryx(searchsTtdQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var sTtd STtdDto
		err = rows.StructScan(&sTtd)
		if err != nil {
			return
		}

		data.Items = append(data.Items, sTtd)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *STtdRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []STtdDto, err error) {

	where := " where coalesce(a.is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and a.id_branch='%s' ", req.IdBranch)
	}

	if req.IdTransaksi != "" {
		where += fmt.Sprintf(" and a.jenis='%s' ", req.IdTransaksi)
	}

	rows, err := r.DB.Read.Queryx(sTtdQuery.SelectDto + where)
	// rows, err := r.DB.Read.Queryx(sTtdQuery.Select + " where coalesce(is_deleted, false) = false order by kode")
	if err == sql.ErrNoRows {
		_ = failure.NotFound("STtd NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList STtdDto
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *STtdRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data STtd, err error) {
	err = r.DB.Read.Get(&data, sTtdQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *STtdRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(sTtdQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *STtdRepositoryPostgreSQL) Update(data STtd) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *STtdRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data STtd) (err error) {
	stmt, err := tx.PrepareNamed(sTtdQuery.Update + " WHERE id=:id")
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
