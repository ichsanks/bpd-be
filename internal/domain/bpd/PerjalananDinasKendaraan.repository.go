package bpd

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
	pdKendaraanQuery = struct {
		Select    string
		SelectDTO string
		Insert    string
		Update    string
		Delete    string
		Count     string
		Exist     string
	}{
		Select: `select id, id_perjalanan_dinas, id_kendaraan, id_pegawai, nama_pengemudi, created_at, created_by, updated_at, updated_by, is_deleted from perjalanan_dinas_kendaraan`,
		SelectDTO: `select pk.id, pk.id_perjalanan_dinas, pk.id_kendaraan, pk.id_pegawai, coalesce(pk.nama_pengemudi, mp.nama) nama_pengemudi, pk.created_at, pk.created_by, pk.updated_at, pk.updated_by, pk.is_deleted,
				mk.nama nama_kendaraan, mk.nopol, mk.id_jenis_kendaraan, mj.nama nama_jenis_kendaraan
				from perjalanan_dinas_kendaraan pk
				left join m_kendaraan mk on mk.id = pk.id_kendaraan
				left join m_jenis_kendaraan mj on mj.id = mk.id_jenis_kendaraan
				left join m_pegawai mp on mp.id = pk.id_pegawai
		`,
		Insert: `insert into perjalanan_dinas_kendaraan
				(id, id_perjalanan_dinas, id_kendaraan, id_pegawai, nama_pengemudi, created_at, created_by) 
				values 
				(:id, :id_perjalanan_dinas, :id_kendaraan, :id_pegawai, :nama_pengemudi, :created_at, :created_by) `,
		Update: `update perjalanan_dinas_kendaraan set
				id=:id,
				id_perjalanan_dinas=:id_perjalanan_dinas,
				id_kendaraan=:id_kendaraan,
				id_pegawai=:id_pegawai,
				nama_pengemudi=:nama_pengemudi,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from perjalanan_dinas_kendaraan `,
		Count: `select count (id)
				from perjalanan_dinas_kendaraan `,
		Exist: `select count(id)>0 from perjalanan_dinas_kendaraan `,
	}
)

type PerjalananDinasKendaraanRepository interface {
	Create(data PerjalananDinasKendaraan) error
	GetAll(idPerjalananDinas string) (data []PerjalananDinasKendaraanDTO, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data PerjalananDinasKendaraan, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data PerjalananDinasKendaraan) error
}

type PerjalananDinasKendaraanRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvidePerjalananDinasKendaraanRepositoryPostgreSQL(db *infras.PostgresqlConn) *PerjalananDinasKendaraanRepositoryPostgreSQL {
	s := new(PerjalananDinasKendaraanRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *PerjalananDinasKendaraanRepositoryPostgreSQL) Create(data PerjalananDinasKendaraan) error {
	stmt, err := r.DB.Write.PrepareNamed(pdKendaraanQuery.Insert)
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

func (r *PerjalananDinasKendaraanRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(pk.is_deleted, false) = false ")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(mk.nama, mk.nopol, mj.nama, coalesce(pk.nama_pengemudi, mp.nama)) ilike ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	if req.IdPerjalananDinas != "" {
		searchRoleBuff.WriteString(" AND pk.id_perjalanan_dinas = ? ")
		searchParams = append(searchParams, req.IdPerjalananDinas)
	}

	query := r.DB.Read.Rebind("select count(*) from(" + pdKendaraanQuery.SelectDTO + searchRoleBuff.String() + ")x")

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

	searchRoleBuff.WriteString("order by " + ColumnMappPDKendaraan[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchPdKendaraanQuery := searchRoleBuff.String()
	searchPdKendaraanQuery = r.DB.Read.Rebind(pdKendaraanQuery.SelectDTO + searchPdKendaraanQuery)
	fmt.Println("query", searchPdKendaraanQuery)
	rows, err := r.DB.Read.Queryx(searchPdKendaraanQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var items PerjalananDinasKendaraanDTO
		err = rows.StructScan(&items)
		if err != nil {
			return
		}

		data.Items = append(data.Items, items)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *PerjalananDinasKendaraanRepositoryPostgreSQL) GetAll(idPerjalananDinas string) (data []PerjalananDinasKendaraanDTO, err error) {
	rows, err := r.DB.Read.Queryx(pdKendaraanQuery.SelectDTO+" where coalesce(pk.is_deleted, false) = false and pk.id_perjalanan_dinas=$1 order by pk.created_at desc", idPerjalananDinas)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Perjalanan Dinas Kendaraan NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items PerjalananDinasKendaraanDTO
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	return
}

func (r *PerjalananDinasKendaraanRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data PerjalananDinasKendaraan, err error) {
	err = r.DB.Read.Get(&data, pdKendaraanQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *PerjalananDinasKendaraanRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(pdKendaraanQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *PerjalananDinasKendaraanRepositoryPostgreSQL) Update(data PerjalananDinasKendaraan) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *PerjalananDinasKendaraanRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data PerjalananDinasKendaraan) (err error) {
	stmt, err := tx.PrepareNamed(pdKendaraanQuery.Update + " WHERE id=:id")
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
