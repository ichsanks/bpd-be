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
	unorQuery = struct {
		Select      string
		SelectDTO   string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id, kode, nama, id_bidang, tenant_id, id_branch, created_at, created_by, updated_at, updated_by, is_deleted from m_unit_organisasi_kerja `,
		SelectDTO: `select u.id, u.kode, u.nama, u.id_bidang, b.kode kode_bidang, b.nama nama_bidang, u.created_at, u.created_by, u.updated_at, u.updated_by, u.is_deleted
				from m_unit_organisasi_kerja u
				left join m_bidang b on ( u.id_bidang = b.id and b.is_deleted = false )
				`,
		Insert: `insert into m_unit_organisasi_kerja
				(id, kode, nama, id_bidang, tenant_id, id_branch, created_at, created_by)
				values
				(:id, :kode, :nama, :id_bidang, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_unit_organisasi_kerja set
				id=:id,
				kode=:kode,
				nama=:nama,
				id_bidang=:id_bidang,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_unit_organisasi_kerja `,
		Count: `select count (id)
				from m_unit_organisasi_kerja `,
		Exist: `select count(id)>0 from m_unit_organisasi_kerja `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_pegawai pd 
			where id_unor = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type UnitKerjaRepository interface {
	Create(data UnitOrganisasiKerja) error
	GetAll(req model.StandardRequest) (data []UnitOrganisasiKerjaDTO, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data UnitOrganisasiKerja, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data UnitOrganisasiKerja) error
	ExistByNama(nama string, id string, idBranch string) (bool, error)
	ExistByKode(kode string, id string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type UnitKerjaRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideUnitKerjaRepositoryPostgreSQL(db *infras.PostgresqlConn) *UnitKerjaRepositoryPostgreSQL {
	s := new(UnitKerjaRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *UnitKerjaRepositoryPostgreSQL) Create(data UnitOrganisasiKerja) error {
	stmt, err := r.DB.Write.PrepareNamed(unorQuery.Insert)
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

func (r *UnitKerjaRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(u.is_deleted, false) = false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND u.id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.IdBidang != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" u.id_bidang = ? ")
		searchParams = append(searchParams, req.IdBidang)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(u.kode, u.nama, b.nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind("select count(*) from(" + unorQuery.SelectDTO + searchRoleBuff.String() + ")x")

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

	searchRoleBuff.WriteString("order by " + ColumnMappUnor[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchunorQuery := searchRoleBuff.String()
	searchunorQuery = r.DB.Read.Rebind(unorQuery.SelectDTO + searchunorQuery)
	fmt.Println("query", searchunorQuery)
	rows, err := r.DB.Read.Queryx(searchunorQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var unor UnitOrganisasiKerjaDTO
		err = rows.StructScan(&unor)
		if err != nil {
			return
		}

		data.Items = append(data.Items, unor)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *UnitKerjaRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []UnitOrganisasiKerjaDTO, err error) {
	criteria := ` where coalesce(u.is_deleted, false) = false `
	if req.IdBranch != "" {
		criteria += fmt.Sprintf(` and u.id_branch='%s'`, req.IdBranch)
	}
	if req.IdBidang != "" {
		criteria += fmt.Sprintf(` and u.id_bidang='%s'`, req.IdBidang)
	}
	criteria += fmt.Sprintf(` order by kode asc`)
	rows, err := r.DB.Read.Queryx(unorQuery.SelectDTO + criteria)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Unit Kerja NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList UnitOrganisasiKerjaDTO
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *UnitKerjaRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data UnitOrganisasiKerja, err error) {
	err = r.DB.Read.Get(&data, unorQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *UnitKerjaRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(unorQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *UnitKerjaRepositoryPostgreSQL) Update(data UnitOrganisasiKerja) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *UnitKerjaRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data UnitOrganisasiKerja) (err error) {
	stmt, err := tx.PrepareNamed(unorQuery.Update + " WHERE id=:id")
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

func (r *UnitKerjaRepositoryPostgreSQL) ExistByNama(nama string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch <> '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, unorQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *UnitKerjaRepositoryPostgreSQL) ExistByKode(kode string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(kode)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch <> '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, unorQuery.Exist+criteria, kode)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *UnitKerjaRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, unorQuery.ExistRelasi, id)

	return
}
