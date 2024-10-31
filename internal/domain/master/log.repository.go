package master

import (
	"bytes"
	"fmt"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

var (
	logQuery = struct {
		Select string
		Count  string
	}{
		Select: `select la.id, la.actions, la.jam, la.kode, la.keterangan, la.id_user, u.username, la.platform, la.ip_address, la.user_agent 
			from dapodik.log_activity la 
			left join dapodik.users u on la.id_user = u.id `,
		Count: `select count (la.id)
			from dapodik.log_activity la 
			left join dapodik.users u on la.id_user = u.id `,
	}
)

type LogRepository interface {
	ResolveAll(req model.StandardRequestPegawai) (data pagination.Response, err error)
}

type LogRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideLogRepositoryPostgreSQL(db *infras.PostgresqlConn) *LogRepositoryPostgreSQL {
	s := new(LogRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *LogRepositoryPostgreSQL) ResolveAll(req model.StandardRequestPegawai) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" where la.jam::date between ? and ? ")
	searchParams = append(searchParams, req.StartDate)
	searchParams = append(searchParams, req.EndDate)

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(la.actions, la.keterangan, la.kode, u.username, la.platform) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(logQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappLog[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchLogQuery := searchRoleBuff.String()
	searchLogQuery = r.DB.Read.Rebind(logQuery.Select + searchLogQuery)
	fmt.Println("query", searchLogQuery)
	rows, err := r.DB.Read.Queryx(searchLogQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var log Log
		err = rows.StructScan(&log)
		if err != nil {
			return
		}

		data.Items = append(data.Items, log)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}
