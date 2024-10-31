package master

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

var (
	ruleApprovalQuery = struct {
		Select        string
		SelectDTO     string
		SelectListDTO string
		Insert        string
		Update        string
		SelectTtd     string
	}{
		Select: `select id, nama, id_fungsionalitas, jenis, updated_at, created_at, created_by, updated_by, is_deleted, group_rule, id_pegawai, tenant_id, id_branch from m_rule_approval`,
		SelectDTO: `select r.id, r.nama, r.id_fungsionalitas, f.nama nama_fungsionalitas, r.jenis, r.created_at, r.created_by, r.updated_at, r.updated_by, r.is_deleted, r.group_rule, r.id_pegawai from m_rule_approval r
			left join m_fungsionalitas f on r.id_fungsionalitas = f.id
		`,
		SelectListDTO: `SELECT * FROM (
			select r.id, r.nama, r.id_fungsionalitas, f.nama nama_fungsionalitas, r.jenis, r.created_at, 
			r.created_by, r.updated_at, r.updated_by, r.is_deleted, r.group_rule, r.id_pegawai, j.nama nama_jenis, r.id_branch,
			(case when group_rule=1 then 'ALL PEGAWAI'
				when group_rule=2 then 'FUNGSIONALITAS'
				when group_rule=3 then 'ID PEGAWAI'
			else '' end
			) nama_group_rule,
			(case when group_rule=1 then '-'
				when group_rule=2 then f.nama
				when group_rule=3 then p.nama
			else '' end
			) group_value 
			from m_rule_approval r
			left join m_fungsionalitas f on r.id_fungsionalitas = f.id
			left join m_jenis_approval j on j.id=r.jenis
			left join m_pegawai p on p.id=r.id_pegawai
			where coalesce(r.is_deleted, false)=false
		)r `,
		Insert: `INSERT INTO m_rule_approval (id, nama, id_fungsionalitas, jenis, created_by, created_at, group_rule, id_pegawai, tenant_id, id_branch) values(:id, :nama, :id_fungsionalitas, :jenis, :created_by, :created_at, :group_rule, :id_pegawai, :tenant_id, :id_branch) `,
		Update: `UPDATE m_rule_approval SET 
			id=:id, 
			nama=:nama, 
			id_fungsionalitas=:id_fungsionalitas, 
			jenis=:jenis,
			group_rule=:group_rule,
			id_pegawai=:id_pegawai,
			tenant_id=:tenant_id,
			id_branch=:id_branch,
			updated_at=:updated_at,
			updated_by=:updated_by, 
			is_deleted=:is_deleted`,
		SelectTtd: `select a.id, a.id_pegawai, b.nama nama_pegawai, b.nip , c.nama nama_bidang, d.nama nama_jabatan  from m_rule_approval_detail a
		left join m_pegawai b on b.id = a.id_pegawai
		left join m_bidang c on c.id = b.id_bidang
		left join m_jabatan d on d.id = b.id_jabatan`,
	}
)

var (
	ruleApprovalDetailQuery = struct {
		Select                string
		SelectDTO             string
		SelectNext            string
		InsertBulk            string
		InsertBulkPlaceholder string
	}{
		Select: `select id, id_rule_approval, id_fungsionalitas, id_unor, id_bidang, type_approval, urut, group_approval, feedback_tolak, approval_line, id_pegawai, esign, ket_ttd from m_rule_approval_detail`,
		SelectNext: `select r.id, r.id_rule_approval, r.id_fungsionalitas, r.id_unor, r.id_bidang, r.type_approval, r.urut, r.group_approval, r.feedback_tolak, r.approval_line, r.id_pegawai, f.is_head, r.esign, r.ket_ttd
				from m_rule_approval_detail r
				left join m_fungsionalitas f on r.id_fungsionalitas = f.id `,
		SelectDTO: `select r.id, r.id_rule_approval, r.id_fungsionalitas, f.nama nama_fungsionalitas, r.id_unor, u.kode kode_unor, u.nama nama_unor,
				r.id_bidang, b.kode kode_bidang, b.nama nama_bidang, r.type_approval, r.urut, r.group_approval, r.feedback_tolak, r.approval_line, r.id_pegawai, p.nama pegawai, f.is_head, r.esign , r.ket_ttd
				from m_rule_approval_detail r
				left join m_rule_approval ap on ap.id = r.id_rule_approval
				left join m_fungsionalitas f on r.id_fungsionalitas = f.id
				left join m_bidang b on r.id_bidang = b.id
				left join m_unit_organisasi_kerja u on r.id_unor = u.id
				left join m_pegawai p on p.id = r.id_pegawai
		`,
		InsertBulk:            `INSERT INTO public.m_rule_approval_detail(id, id_rule_approval, id_fungsionalitas, id_unor, id_bidang, type_approval, urut, group_approval, feedback_tolak, approval_line, id_pegawai, esign, ket_ttd) values `,
		InsertBulkPlaceholder: ` (:id, :id_rule_approval, :id_fungsionalitas, :id_unor, :id_bidang, :type_approval, :urut, :group_approval, :feedback_tolak, :approval_line, :id_pegawai, :esign, :ket_ttd) `,
	}
)

type RuleApprovalRepository interface {
	Create(data RuleApproval) error
	Update(data RuleApproval) error
	UpdateRuleApproval(data RuleApproval) error
	ResolveByIDDTO(id string) (data RuleApprovalDTO, err error)
	ResolveAll(req model.StandardRequestRuleApproval) (data pagination.Response, err error)
	ResolveByID(id string) (data RuleApproval, err error)
	GetAll(idFungsionalitas string) (data []RuleApprovalDTO, err error)
	GetAllRuleApprovalDetail(idRule string, typeApproval string) (data []RuleApprovalDetailDTO, err error)
	ResolveRuleApprovalDetail(id string) (data RuleApprovalDetail, err error)
	ResolveRuleApprovalDetailDTO(id string) (data RuleApprovalDetailDTO, err error)
	GetNextRuleApprovalDetail(id string, typeApproval string) (data RuleApprovalDetailDTO, err error)
	ResolveByKode(req RuleParams) (data RuleApproval, err error)
	GetAllRuleApprovalDetailByKode(req RuleParams) (data []RuleApprovalDetailDTO, err error)
	ResolveTtd(typeApproval string) (data RuleApprovalTtd, err error)
}

type RuleApprovalRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideRuleApprovalRepositoryPostgreSQL(db *infras.PostgresqlConn) *RuleApprovalRepositoryPostgreSQL {
	s := new(RuleApprovalRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *RuleApprovalRepositoryPostgreSQL) ResolveAll(req model.StandardRequestRuleApproval) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(r.is_deleted, false) = false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND r.id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(r.nama_jenis, r.nama_group_rule, r.group_value) ilike ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	if req.Jenis != "" {
		searchRoleBuff.WriteString(" AND r.jenis=? ")
		searchParams = append(searchParams, req.Jenis)
	}

	query := r.DB.Read.Rebind("select count(*) from (" + ruleApprovalQuery.SelectListDTO + searchRoleBuff.String() + ")s")

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

	searchRoleBuff.WriteString("order by " + ColumnMappRuleApproval[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchSiswaQuery := searchRoleBuff.String()
	searchSiswaQuery = r.DB.Read.Rebind(ruleApprovalQuery.SelectListDTO + searchSiswaQuery)
	rows, err := r.DB.Read.Queryx(searchSiswaQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var items RuleApprovalDTO
		err = rows.StructScan(&items)
		if err != nil {
			return
		}

		data.Items = append(data.Items, items)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}

// Function digunakan untuk create with transaction
func (r *RuleApprovalRepositoryPostgreSQL) Create(data RuleApproval) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		// Function create table m_rule_approval
		if err := r.CreateTxRule(tx, data); err != nil {
			e <- err
			return
		}

		// Function Insert Bulk table rule_approval_detail
		if err := txCreateRuleApprovalDetail(tx, data.Detail); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

// Function digunakan untuk update with transaction
func (r *RuleApprovalRepositoryPostgreSQL) UpdateRuleApproval(data RuleApproval) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		// Function update table m_rule_approval
		if err := r.UpdateTxRule(tx, data); err != nil {
			e <- err
			return
		}

		// Function delete not in table m_rule_approval_detail
		ids := make([]string, 0)
		for _, d := range data.Detail {
			ids = append(ids, d.ID.String())
		}

		if err := r.txDeleteDetailNotIn(tx, data.ID.String(), ids); err != nil {
			e <- err
			return
		}

		// Function Insert Bulk table m_rule_approval_detail
		if err := txCreateRuleApprovalDetail(tx, data.Detail); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *RuleApprovalRepositoryPostgreSQL) CreateTxRule(tx *sqlx.Tx, data RuleApproval) error {
	stmt, err := tx.PrepareNamed(ruleApprovalQuery.Insert)
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

func (r *RuleApprovalRepositoryPostgreSQL) Update(data RuleApproval) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		// Function update table m_rule_approval
		if err := r.UpdateTxRule(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *RuleApprovalRepositoryPostgreSQL) UpdateTxRule(tx *sqlx.Tx, data RuleApproval) error {
	stmt, err := tx.PrepareNamed(ruleApprovalQuery.Update + " WHERE id=:id")
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return nil
}

func txCreateRuleApprovalDetail(tx *sqlx.Tx, details []RuleApprovalDetail) (err error) {
	if len(details) == 0 {
		return
	}
	query, args, err := composeBulkUpsertRuleApprovalDetailQuery(details)
	if err != nil {
		return
	}

	query = tx.Rebind(query)
	stmt, err := tx.Preparex(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Stmt.Exec(args...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func composeBulkUpsertRuleApprovalDetailQuery(details []RuleApprovalDetail) (qResult string, params []interface{}, err error) {
	values := []string{}
	for _, d := range details {
		param := map[string]interface{}{
			"id":                d.ID,
			"id_rule_approval":  d.IdRuleApproval,
			"id_fungsionalitas": d.IdFungsionalitas,
			"id_unor":           d.IdUnor,
			"id_bidang":         d.IdBidang,
			"type_approval":     d.TypeApproval,
			"urut":              d.Urut,
			"group_approval":    d.GroupApproval,
			"feedback_tolak":    d.FeedbackTolak,
			"approval_line":     d.ApprovalLine,
			"id_pegawai":        d.IdPegawai,
			"esign":             d.Esign,
			"ket_ttd":           d.KetTtd,
		}
		q, args, err := sqlx.Named(ruleApprovalDetailQuery.InsertBulkPlaceholder, param)
		if err != nil {
			return qResult, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}
	qResult = fmt.Sprintf(`%v %v 
						ON CONFLICT (id) 
						DO UPDATE SET 
						id_fungsionalitas=EXCLUDED.id_fungsionalitas, 
						id_unor=EXCLUDED.id_unor, 
						id_bidang=EXCLUDED.id_bidang, 
						type_approval=EXCLUDED.type_approval, 
						urut=EXCLUDED.urut,
						group_approval=EXCLUDED.group_approval, 
						feedback_tolak=EXCLUDED.feedback_tolak, 
						approval_line=EXCLUDED.approval_line, 
						id_pegawai=EXCLUDED.id_pegawai,
						esign=EXCLUDED.esign, 
						ket_ttd=EXCLUDED.ket_ttd`, ruleApprovalDetailQuery.InsertBulk, strings.Join(values, ","))
	return
}

func (r *RuleApprovalRepositoryPostgreSQL) txDeleteDetailNotIn(tx *sqlx.Tx, idRule string, ids []string) (err error) {
	query, args, err := sqlx.In("delete from m_rule_approval_detail where id_rule_approval = ? AND id NOT IN (?)", idRule, ids)
	query = tx.Rebind(query)

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	res, err := r.DB.Write.Exec(query, args...)
	_, err = res.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (r *RuleApprovalRepositoryPostgreSQL) ResolveByID(id string) (data RuleApproval, err error) {
	err = r.DB.Read.Get(&data, ruleApprovalQuery.Select+" WHERE id=$1 ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *RuleApprovalRepositoryPostgreSQL) ResolveByIDDTO(id string) (data RuleApprovalDTO, err error) {
	err = r.DB.Read.Get(&data, ruleApprovalQuery.SelectListDTO+" WHERE r.id=$1 ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	ruleDetail, err := r.GetAllRuleApprovalDetail(id, "")
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	data.Detail = ruleDetail
	return
}

func (r *RuleApprovalRepositoryPostgreSQL) GetAll(idFungsionalitas string) (data []RuleApprovalDTO, err error) {
	criteria := ` where coalesce(r.is_deleted, false)=false`
	if idFungsionalitas != "" {
		criteria += fmt.Sprintf(` and r.id_fungsionalitas='%s' `, idFungsionalitas)
	}

	rows, err := r.DB.Read.Queryx(ruleApprovalQuery.SelectDTO + criteria)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Rule Approval NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items RuleApprovalDTO
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}

	if data == nil {
		data = make([]RuleApprovalDTO, 0)
	}

	return
}

func (r *RuleApprovalRepositoryPostgreSQL) GetAllRuleApprovalDetail(idRule string, typeApproval string) (data []RuleApprovalDetailDTO, err error) {
	criteria := ` where r.id_rule_approval=$1 `
	if typeApproval != "" {
		criteria += fmt.Sprintf(` and r.type_approval='%s' `, typeApproval)
	}
	criteria += ` order by r.urut asc `

	rows, err := r.DB.Read.Queryx(ruleApprovalDetailQuery.SelectDTO+criteria, idRule)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var master RuleApprovalDetailDTO
		err = rows.StructScan(&master)

		if err != nil {
			return
		}

		data = append(data, master)
	}
	return
}

func (r *RuleApprovalRepositoryPostgreSQL) ResolveRuleApprovalDetail(id string) (data RuleApprovalDetail, err error) {
	err = r.DB.Read.Get(&data, ruleApprovalDetailQuery.Select+" where id = $1", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *RuleApprovalRepositoryPostgreSQL) ResolveRuleApprovalDetailDTO(id string) (data RuleApprovalDetailDTO, err error) {
	err = r.DB.Read.Get(&data, ruleApprovalDetailQuery.SelectDTO+" where r.id = $1", id)
	if err != nil {
		// logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *RuleApprovalRepositoryPostgreSQL) GetNextRuleApprovalDetail(id string, typeApproval string) (data RuleApprovalDetailDTO, err error) {
	criteria := ` 
		where r.id_rule_approval = (select id_rule_approval from m_rule_approval_detail where id = $1)
		and r.urut = (
			SELECT coalesce(SUM(case when coalesce(urut,0)>0 then urut+1 else 1 end), 1) FROM (
				select urut urut from m_rule_approval_detail 
				where id = $1
				and type_approval = $2
			)x
		)
		and r.type_approval = $2
		order by r.urut asc
		limit 1	
	`
	err = r.DB.Read.Get(&data, ruleApprovalDetailQuery.SelectNext+criteria, id, typeApproval)
	if err != nil {
		// logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *RuleApprovalRepositoryPostgreSQL) ResolveByKode(req RuleParams) (data RuleApproval, err error) {
	where := " WHERE jenis=$1 "
	if req.IdPegawai != "" {
		where += fmt.Sprintf(" and id_pegawai='%v' ", req.IdPegawai)
	}

	if req.IdFungsionalitas != "" {
		where += fmt.Sprintf(" and id_fungsionalitas='%v' ", req.IdFungsionalitas)
	}

	if req.GroupRule != 0 {
		where += fmt.Sprintf(" and group_rule='%v' ", req.GroupRule)
	}

	err = r.DB.Read.Get(&data, ruleApprovalQuery.Select+where, req.Jenis)
	if err != nil {
		return
	}

	return
}

func (r *RuleApprovalRepositoryPostgreSQL) GetAllRuleApprovalDetailByKode(req RuleParams) (data []RuleApprovalDetailDTO, err error) {
	criteria := ` where ap.is_deleted=false 
		and ap.jenis=$1 
		and r.type_approval=$2 
		and ap.group_rule=$3 `

	criteria += " order by r.urut asc"
	fmt.Println("req:", req)
	rows, err := r.DB.Read.Queryx(ruleApprovalDetailQuery.SelectDTO+criteria, req.Jenis, req.TypeApproval, req.GroupRule)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Data Not Found")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var master RuleApprovalDetailDTO
		err = rows.StructScan(&master)

		if err != nil {
			return
		}

		data = append(data, master)
	}
	return
}

func (r *RuleApprovalRepositoryPostgreSQL) ResolveTtd(id string) (data RuleApprovalTtd, err error) {
	err = r.DB.Read.Get(&data, ruleApprovalQuery.SelectTtd+" where a.type_approval=$1 ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}
