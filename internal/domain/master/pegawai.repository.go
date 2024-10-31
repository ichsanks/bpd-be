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
	pegawaiQuery = struct {
		Select      string
		SelectDTO   string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `SELECT id, nip, nama, jenis_kelamin, agama, alamat, no_hp, email, id_unor, id_jabatan, id_golongan, created_at, updated_at, created_by, updated_by, is_deleted, foto, id_fungsionalitas, id_approval_line, id_manager, nik, foto_ttd, tenant_id, id_branch, id_status_pegawai, id_job_grade, id_person_grade, id_level_bod, id_status_kontrak, kode_vendor
		FROM public.m_pegawai `,
		SelectDTO: `select mp.id, mp.nip, mp.nama, mp.jenis_kelamin, mjk.nama nama_jk, mp.agama, ma.nama nama_agama, mp.alamat, mp.no_hp, mp.email, mp.id_unor, muok.kode kode_unor, muok.nama nama_unor, mp.id_jabatan, mj.nama nama_jabatan, mp.id_golongan, mg.nama nama_golongan, 
		mp.id_fungsionalitas, mf.nama nama_fungsionalitas, mp.id_bidang, mb.nama nama_bidang, mp.created_at, mp.created_by, mp.updated_at, mp.updated_by, mp.is_deleted, mp.id_approval_line, mp.id_manager, mp.nik, mp.foto_ttd, mp.tenant_id, mp.id_branch, mp.id_status_pegawai, mp.id_job_grade, mp.id_person_grade, mp.id_level_bod, mp.id_status_kontrak, mp.kode_vendor
		from m_pegawai mp
		left join m_agama ma on mp.agama = ma.kode
		left join m_jenis_kelamin mjk on mp.jenis_kelamin = mjk.kode
		left join m_unit_organisasi_kerja muok on mp.id_unor = muok.id
		left join m_bidang mb on mb.id = mp.id_bidang
		left join m_jabatan mj on mp.id_jabatan = mj.id
		left join m_golongan mg on mp.id_golongan = mg.id
		left join m_fungsionalitas mf on mp.id_fungsionalitas = mf.id		
		`,
		Insert: `insert into m_pegawai
				(id, nip, nama, jenis_kelamin, agama, alamat, no_hp, email, id_unor, id_bidang, id_jabatan, id_golongan, id_fungsionalitas, created_at, created_by, nik, foto_ttd, tenant_id, id_branch, id_status_pegawai, id_job_grade, id_person_grade, id_level_bod, id_status_kontrak, kode_vendor)
				values
				(:id, :nip, :nama, :jenis_kelamin, :agama, :alamat, :no_hp, :email, :id_unor, :id_bidang, :id_jabatan, :id_golongan, :id_fungsionalitas, :created_at, :created_by, :nik, :foto_ttd, :tenant_id, :id_branch, :id_status_pegawai, :id_job_grade, :id_person_grade, :id_level_bod, :id_status_kontrak, :kode_vendor) `,
		Update: `update m_pegawai set
				id=:id,
				nip=:nip,
				nama=:nama,
				jenis_kelamin=:jenis_kelamin,
				agama=:agama,
				alamat=:alamat,
				no_hp=:no_hp,
				email=:email,
				id_unor=:id_unor,
				id_bidang=:id_bidang,
				id_jabatan=:id_jabatan,
				id_golongan=:id_golongan,
				id_fungsionalitas=:id_fungsionalitas,
				nik=:nik,
				foto_ttd=:foto_ttd,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				id_status_pegawai=:id_status_pegawai,
				id_job_grade=:id_job_grade,
				id_person_grade=:id_person_grade,
				id_level_bod=:id_level_bod,
				id_status_kontrak=:id_status_kontrak,
				kode_vendor=:kode_vendor
				 `,
		Delete: `delete from m_pegawai `,
		Count: `select count (id)
				from m_pegawai `,
		Exist: `select count(id)>0 from m_pegawai `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from surat_perjalanan_dinas pd 
			where id_pegawai = $1
			and id_branch = $2
			and coalesce(is_deleted, false) is false 
			for update 
		) x `,
	}
)

type PegawaiRepository interface {
	Create(data Pegawai) error
	GetAll(req model.StandardRequest) (data []Pegawai, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Pegawai, err error)
	ResolveByIDDTO(id uuid.UUID) (Pegawai PegawaiDTO, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data Pegawai) error
	ExistByNip(nip string, idBranch string) (bool, error)
	ExistByNama(nama string, idBranch string) (bool, error)
	ExistByNipID(id uuid.UUID, nip string, idBranch string) (bool, error)
	ExistByNamaID(id uuid.UUID, nama string, idBranch string) (bool, error)
	GetAllPegawai(req PegawaiParams) (data []PegawaiDTO, errr error)
	ExistRelasiStatus(id uuid.UUID, idBranch string) (exist bool)
}

type PegawaiRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvidePegawaiRepositoryPostgreSQL(db *infras.PostgresqlConn) *PegawaiRepositoryPostgreSQL {
	s := new(PegawaiRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *PegawaiRepositoryPostgreSQL) Create(data Pegawai) error {
	stmt, err := r.DB.Read.PrepareNamed(pegawaiQuery.Insert)
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

func (r *PegawaiRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE mp.is_deleted is false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND mp.id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(mp.nip, mp.nama, muok.nama, mj.nama, mg.nama, mf.nama, mb.nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	if req.IdFungsionalitas != "" {
		searchRoleBuff.WriteString(" AND mp.id_fungsionalitas = ? ")
		searchParams = append(searchParams, req.IdFungsionalitas)
	}

	if req.IdBidang != "" {
		searchRoleBuff.WriteString(" AND mp.id_bidang = ? ")
		searchParams = append(searchParams, req.IdBidang)
	}

	if req.IdUnor != "" {
		searchRoleBuff.WriteString(" AND mp.id_unor = ? ")
		searchParams = append(searchParams, req.IdUnor)
	}

	query := r.DB.Read.Rebind("select count(x.id) from(" + pegawaiQuery.SelectDTO + searchRoleBuff.String() + ")x")
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

	searchRoleBuff.WriteString("order by " + ColumnMappPegawai[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchPegawaiQuery := searchRoleBuff.String()
	searchPegawaiQuery = r.DB.Read.Rebind(pegawaiQuery.SelectDTO + searchPegawaiQuery)
	rows, err := r.DB.Read.Queryx(searchPegawaiQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var Pegawai PegawaiDTO
		err = rows.StructScan(&Pegawai)
		if err != nil {
			return
		}

		data.Items = append(data.Items, Pegawai)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *PegawaiRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []Pegawai, err error) {

	where := " where coalesce(is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and id_branch='%s' ", req.IdBranch)
	}

	rows, err := r.DB.Read.Queryx(pegawaiQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Pegawai NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList Pegawai
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *PegawaiRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (Pegawai Pegawai, err error) {
	err = r.DB.Read.Get(&Pegawai, pegawaiQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *PegawaiRepositoryPostgreSQL) ResolveByIDDTO(id uuid.UUID) (Pegawai PegawaiDTO, err error) {
	err = r.DB.Read.Get(&Pegawai, pegawaiQuery.SelectDTO+" WHERE mp.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *PegawaiRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(pegawaiQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *PegawaiRepositoryPostgreSQL) Update(data Pegawai) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdatePegawai(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txUpdatePegawai(tx *sqlx.Tx, data Pegawai) (err error) {
	stmt, err := tx.PrepareNamed(pegawaiQuery.Update + " WHERE id=:id")
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

func (r *PegawaiRepositoryPostgreSQL) ExistByNip(nip string, idBranch string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, pegawaiQuery.Exist+" where upper(nip)=upper($1) and coalesce(is_deleted, false)=false and id_branch=$2 ", nip, idBranch)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *PegawaiRepositoryPostgreSQL) ExistByNama(nama string, idBranch string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, pegawaiQuery.Exist+" where upper(nama)=upper($1) and coalesce(is_deleted, false)=false and id_branch=$2 ", nama, idBranch)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *PegawaiRepositoryPostgreSQL) ExistByNipID(id uuid.UUID, nip string, idBranch string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, pegawaiQuery.Exist+" where id <> $1 and upper(nip)=upper($2) and coalesce(is_deleted, false)=false and id_branch=$3 ", id, nip, idBranch)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *PegawaiRepositoryPostgreSQL) ExistByNamaID(id uuid.UUID, nama string, idBranch string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, pegawaiQuery.Exist+" where id <> $1 and upper(nama)=upper($2) and coalesce(is_deleted, false)=false and id_branch=$3 ", id, nama, idBranch)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *PegawaiRepositoryPostgreSQL) GetAllPegawai(req PegawaiParams) (data []PegawaiDTO, errr error) {
	criteria := " WHERE coalesce(mp.is_deleted,false)=false "
	if req.IdPegawai != "" {
		criteria += fmt.Sprintf(" AND mp.id='%v' ", req.IdPegawai)
	}

	if req.IdFungsionalitas != "" {
		criteria += fmt.Sprintf(" AND mp.id_fungsionalitas='%v' ", req.IdFungsionalitas)
	}

	if req.IdBidang != "" {
		criteria += fmt.Sprintf(" AND mp.id_bidang='%v' ", req.IdBidang)
	}

	if req.IdUnor != "" {
		criteria += fmt.Sprintf(" AND mp.id_unor='%v' ", req.IdUnor)
	}

	rows, err := r.DB.Read.Queryx(pegawaiQuery.SelectDTO + criteria)
	if err == sql.ErrNoRows {
		errr = failure.NotFound("Pegawai")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var master PegawaiDTO
		err = rows.StructScan(&master)
		if err != nil {
			return
		}

		data = append(data, master)
	}
	return
}

func (r *PegawaiRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID, idBranch string) (exist bool) {
	r.DB.Read.Get(&exist, pegawaiQuery.ExistRelasi, id, idBranch)

	return
}
