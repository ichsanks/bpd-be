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
	syaratSyaratDokumenQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select a.id, a.id_transaksi, a.id_dokumen, a.is_mandatory, a.tenant_id, a.id_branch, a.created_at, a.created_by, a.updated_at, a.updated_by, a.is_deleted, b.nama,
		(CASE WHEN a.id_transaksi='1' THEN 'SPPD'
		WHEN a.id_transaksi='2' THEN 'BPD'
		ELSE '' END) AS jenis_transaksi,
		(CASE WHEN a.is_mandatory=true THEN 'YA'
		WHEN a.is_mandatory=false THEN 'TIDAK'
		ELSE '' END) AS mandatori
		from m_syarat_dokumen a
		        left  join m_dokumen b on b.id = a.id_dokumen`,
		Insert: `insert into m_syarat_dokumen
				(id, id_transaksi, id_dokumen, is_mandatory, tenant_id, id_branch, created_at, created_by)
				values
				(:id, :id_transaksi, :id_dokumen, :is_mandatory, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_syarat_dokumen set
				id=:id,
				id_transaksi=:id_transaksi,
				id_dokumen=:id_dokumen,
				is_mandatory=:is_mandatory,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_syarat_dokumen `,
		Count: `select count (a.id)
				from m_syarat_dokumen a
				left  join m_dokumen b on b.id = a.id_dokumen `,
		Exist: `select count(id)>0 from m_syarat_dokumen `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_pegawai pd 
			where id = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type SyaratDokumenRepository interface {
	Create(data SyaratDokumen) error
	GetAll(req model.StandardRequest) (data []SyaratDokumen, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data SyaratDokumen, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data SyaratDokumen) error
	ExistByNama(nama string, id string, idBranch string) (bool, error)
	ExistByKode(kode string, id string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type SyaratDokumenRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideSyaratDokumenRepositoryPostgreSQL(db *infras.PostgresqlConn) *SyaratDokumenRepositoryPostgreSQL {
	s := new(SyaratDokumenRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *SyaratDokumenRepositoryPostgreSQL) Create(data SyaratDokumen) error {
	stmt, err := r.DB.Write.PrepareNamed(syaratSyaratDokumenQuery.Insert)
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

func (r *SyaratDokumenRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(a.is_deleted, false) = false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND a.id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(b.nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(syaratSyaratDokumenQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappSyaratDokumen[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchsyaratSyaratDokumenQuery := searchRoleBuff.String()
	searchsyaratSyaratDokumenQuery = r.DB.Read.Rebind(syaratSyaratDokumenQuery.Select + searchsyaratSyaratDokumenQuery)
	fmt.Println("query", searchsyaratSyaratDokumenQuery)
	rows, err := r.DB.Read.Queryx(searchsyaratSyaratDokumenQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var SyaratDokumen SyaratDokumen
		err = rows.StructScan(&SyaratDokumen)
		if err != nil {
			return
		}

		data.Items = append(data.Items, SyaratDokumen)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *SyaratDokumenRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []SyaratDokumen, err error) {
	where := " where coalesce(a.is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and a.id_branch='%s' ", req.IdBranch)
	}

	if req.IdTransaksi != "" {
		where += fmt.Sprintf(" and a.id_transaksi='%s' ", req.IdTransaksi)
	}

	rows, err := r.DB.Read.Queryx(syaratSyaratDokumenQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("SyaratDokumen NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList SyaratDokumen
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *SyaratDokumenRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data SyaratDokumen, err error) {
	err = r.DB.Read.Get(&data, syaratSyaratDokumenQuery.Select+" WHERE a.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *SyaratDokumenRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(syaratSyaratDokumenQuery.Delete+" WHERE a.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *SyaratDokumenRepositoryPostgreSQL) Update(data SyaratDokumen) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *SyaratDokumenRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data SyaratDokumen) (err error) {
	stmt, err := tx.PrepareNamed(syaratSyaratDokumenQuery.Update + " WHERE id=:id")
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

func (r *SyaratDokumenRepositoryPostgreSQL) ExistByNama(nama string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, syaratSyaratDokumenQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *SyaratDokumenRepositoryPostgreSQL) ExistByKode(kode string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(kode)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, syaratSyaratDokumenQuery.Exist+criteria, kode)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *SyaratDokumenRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, syaratSyaratDokumenQuery.ExistRelasi, id)

	return
}
