package auth

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
	userQuery = struct {
		Insert,
		Exist,
		Select,
		SelectDTO,
		SelectVerifikasi,
		Count,
		Update,
		UpdateFcmToken,
		UpdatePassword, resetImei string
	}{
		Insert: `INSERT INTO auth_user (
			id,
			nama,
			username,
			email,
			password,
			status, 
			id_role,
			id_unor,
			id_fungsionalitas,
			id_pegawai,
			active,
			created_at,
			tenant_id, 
			id_branch
		) VALUES (
			:id,
			:nama,
			:username,
			:email,
			:password,
			:status, 
			:id_role,
			:id_unor,
			:id_fungsionalitas,
			:id_pegawai,
			:active,
			:created_at,
			:tenant_id, 
			:id_branch 
		) `,
		Exist: `SELECT COUNT(u.id) > 0 FROM auth_user u`,
		Select: `SELECT u.id, u.nama, u.username, u.email, u.password, u.status, u.id_role, u.id_unor, u.id_fungsionalitas, u.id_pegawai, u.foto, u.active, u.created_at, u.updated_at, u.deleted_at, u.is_deleted, u.tenant_id, u.id_branch
		FROM auth_user u
			left join m_pegawai mp on u.id_pegawai = mp.id `,
		SelectDTO: `SELECT u.id, u.nama, u.username, u.email, u.password, u.status, u.id_role, u.id_unor, u.id_fungsionalitas, u.id_pegawai, u.active,
			mp.nip, mp.nama nama_pegawai, r.nama nama_role, muk.nama nama_unor, mf.nama nama_fungsionalitas, mf.jenis_approval, mp.id_bidang, mb.nama nama_bidang,
			u.tenant_id, u.id_branch, mb2.nama as nama_branch 
			FROM auth_user u
			left join auth_role r on r.id = u.id_role
			left join m_pegawai mp on mp.id = u.id_pegawai
			left join m_unit_organisasi_kerja muk on muk.id = u.id_unor
			left join m_bidang mb on mb.id = mp.id_bidang
			left join m_fungsionalitas mf on mf.id = u.id_fungsionalitas
			left join m_branch mb2 on mb2.id = u.id_branch  `,
		Count: `select count(u.id) from auth_user u 
			left join "auth_role" r on u.id_role = r.id 
			left join m_pegawai mp on u.id_pegawai = mp.id `,
		Update: `UPDATE auth_user SET 
		    id=:id,
			nama=:nama,
			username=:username,
			email=:email,
			status=:status, 
			id_role=:id_role,
			id_unor=:id_unor,
			id_fungsionalitas=:id_fungsionalitas,
			id_pegawai=:id_pegawai,
			active=:active,
			id_branch =:id_branch ,
			is_deleted=:is_deleted,
			deleted_at=:deleted_at,
			updated_at=:updated_at`,
		UpdateFcmToken: `UPDATE auth_user SET 
			firebase_token=:firebase_token, 
			updated_at=:updated_at `,
		UpdatePassword: `UPDATE auth_user SET
			password = :password`,
	}

	loginActivityQuery = struct {
		Insert string
	}{
		Insert: `INSERT INTO log_activity (
			id,
			username,
			jam
		) VALUES (
			:id,
			:username,
			:jam
		)`,
	}
)

// UserRepositoryPostgreSQL digunakan untuk Repository User
type UserRepositoryPostgreSQL struct {
	DB             *infras.PostgresqlConn
	roleRepisitory RoleRepository
}

// ProvideUserRepositoryPostgreSQL is the provider for this repository.
func ProvideUserRepositoryPostgreSQL(db *infras.PostgresqlConn, rr RoleRepository) *UserRepositoryPostgreSQL {
	return &UserRepositoryPostgreSQL{
		DB:             db,
		roleRepisitory: rr,
	}
}

type UserRepository interface {
	ResolveAll(req model.StandardRequestUser) (dataProyek pagination.Response, err error)
	CreateLoginActivity(loginActivity LoginActivity) error
	ExistByUsername(username string) (exist bool, err error)
	ResolveUserByUsername(username string) ([]User, error)
	ResolveUserByUsernameDTO(username string) (UserDTO, error)
	ResolveUserByID(id uuid.UUID) (User, error)
	ResolveUserByRole(roleName string, idBidang string) (data []User, err error)
	ResolveUserByIDDTO(id uuid.UUID) (UserRoleDTO, error)
	TransactionCreateUser(user User) error
	UpdateUser(id uuid.UUID, user User) error
	UpdateUserFcmToken(id uuid.UUID, user User) error
	UpdateUserPassword(id uuid.UUID, user User) error
	ExistByUsernameId(username string, id string) (exist bool, err error)
}

// TransactionCreateUser digunakan untuk menambahkan user baru dalam blok transaction
func (u *UserRepositoryPostgreSQL) TransactionCreateUser(user User) error {
	return u.DB.WithTransaction(func(db *sqlx.Tx, errs chan error) {
		err := u.createUser(user)
		if err != nil {
			errs <- err
			return
		}

		errs <- nil
	})
}

// createUser is method to create a new user
func (u *UserRepositoryPostgreSQL) createUser(user User) error {
	stmt, err := u.DB.Read.PrepareNamed(userQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	_, err = stmt.Exec(user)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return nil
}

// ExistByUsername is function to check that username exist or not
func (u *UserRepositoryPostgreSQL) ExistByUsername(username string) (exist bool, err error) {
	err = u.DB.Read.Get(&exist, userQuery.Exist+" WHERE username = $1 AND u.active is true AND u.is_deleted is false ", username)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return exist, err
}

// ResolveUserByUsername is function resolving user data by username
func (u *UserRepositoryPostgreSQL) ResolveUserByUsername(username string) (user []User, err error) {
	err = u.DB.Read.Select(&user, userQuery.Select+" WHERE u.username = $1 AND u.deleted_at is null", username)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return user, nil
}

// ResolveUserByUsername is function resolving user data by username and role id
func (u *UserRepositoryPostgreSQL) ResolveUserByUsernameDTO(username string) (UserDTO, error) {
	var user UserDTO
	err := u.DB.Read.Get(&user, userQuery.SelectDTO+" WHERE u.username = $1 and u.active = true and u.is_deleted = false  ", username)
	if err != nil {
		logger.ErrorWithStack(err)
		return UserDTO{}, err
	}

	return user, nil
}

// ResolveUserByID is function resolving user data by id
func (u *UserRepositoryPostgreSQL) ResolveUserByID(id uuid.UUID) (User, error) {
	var user User
	err := u.DB.Read.Get(&user, userQuery.Select+" WHERE u.id = $1 AND u.active = true", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, err
		}
		logger.ErrorWithStack(err)
		return User{}, err
	}
	return user, nil
}

// ResolveUserByRole is function resolving user data by role
func (u *UserRepositoryPostgreSQL) ResolveUserByRole(roleName string, idBidang string) (data []User, errr error) {
	rows, err := u.DB.Read.Queryx(userQuery.Select+" WHERE u.role_id = $1  AND u.active = true", roleName)
	if err == sql.ErrNoRows {
		errr = failure.NotFound("User")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var master User
		err = rows.StructScan(&master)
		if err != nil {
			return
		}
		data = append(data, master)
	}
	return
}

// ResolveUserByID is function resolving user data by email
func (u *UserRepositoryPostgreSQL) ResolveUserByIDDTO(id uuid.UUID) (UserRoleDTO, error) {
	var user UserRoleDTO
	err := u.DB.Read.Get(&user, userQuery.SelectDTO+" WHERE u.id = $1 AND u.deleted_at is null", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserRoleDTO{}, err
		}
		logger.ErrorWithStack(err)
		return UserRoleDTO{}, err
	}
	return user, nil
}

// func (u *UserRepositoryPostgreSQL) UpdateUser(id uuid.UUID, user User) error {
// 	return u.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
// 		if err := txUpdateUser(tx, user); err != nil {
// 			e <- err
// 			return
// 		}
// 		e <- nil
// 	})
// }

// func txUpdateUser(tx *sqlx.Tx, data User) (err error) {
// 	stmt, err := tx.PrepareNamed(userQuery.Update + " WHERE id=:id")
// 	if err != nil {
// 		logger.ErrorWithStack(err)
// 		return
// 	}
// 	defer stmt.Close()
// 	_, err = stmt.Exec(data)
// 	if err != nil {
// 		logger.ErrorWithStack(err)
// 	}
// 	return
// }

// UpdateUser is function to update the user entity
func (u *UserRepositoryPostgreSQL) UpdateUser(id uuid.UUID, user User) error {
	fmt.Println("user : ", user)
	stmt, err := u.DB.Read.PrepareNamed(userQuery.Update + " WHERE id = :id")
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	_, err = stmt.Exec(user)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return nil
}

// UpdateUserFcmToken is function to update the user fcm token
func (u *UserRepositoryPostgreSQL) UpdateUserFcmToken(id uuid.UUID, user User) error {
	stmt, err := u.DB.Read.PrepareNamed(userQuery.UpdateFcmToken + " WHERE id = :id")
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	_, err = stmt.Exec(user)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return nil
}

// UpdateUserPassword is function to update the user password
func (u *UserRepositoryPostgreSQL) UpdateUserPassword(id uuid.UUID, user User) error {
	stmt, err := u.DB.Read.PrepareNamed(userQuery.UpdatePassword + " WHERE id = :id")
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	_, err = stmt.Exec(user)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return nil
}

// CreateLoginActivity is function to create log from login activity
func (u *UserRepositoryPostgreSQL) CreateLoginActivity(loginActivity LoginActivity) error {
	stmt, err := u.DB.Read.PrepareNamed(loginActivityQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	_, err = stmt.Exec(loginActivity)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	return nil
}

// ResolveAll digunakan untuk menampilkan semua data
func (r *UserRepositoryPostgreSQL) ResolveAll(req model.StandardRequestUser) (response pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer

	searchRoleBuff.WriteString(" WHERE u.active = ? ")
	searchParams = append(searchParams, true)

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString("  concat (u.nama, u.username, u.email, mp.nama) ilike ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	if req.IdRole != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" u.id_role = ?  ")
		searchParams = append(searchParams, req.IdRole)
	}
	if req.IdBidang != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" mp.id_bidang = ?  ")
		searchParams = append(searchParams, req.IdBidang)
	}
	if req.IdUnor != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" u.id_unor = ?  ")
		searchParams = append(searchParams, req.IdUnor)
	}

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND u.id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	// query := userQuery.Count + searchRoleBuff.String()
	// queryDTO := userQuery.SelectDTO + searchRoleBuff.String()
	// query = r.DB.Read.Rebind(query)
	query := r.DB.Read.Rebind("select count(*) from (" + userQuery.SelectDTO + searchRoleBuff.String() + ")s")
	var totalData int
	err = r.DB.Read.QueryRow(query, searchParams...).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		response.Items = make([]interface{}, 0)
		return
	}

	searchRoleBuff.WriteString("order by " + ColumnMappUser[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchUserQuery := searchRoleBuff.String()
	searchUserQuery = r.DB.Read.Rebind(userQuery.SelectDTO + searchUserQuery)
	rows, err := r.DB.Read.Queryx(searchUserQuery, searchParams...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var userDTO UserDTO
		err = rows.StructScan(&userDTO)
		if err != nil {
			return
		}

		response.Items = append(response.Items, userDTO)
	}

	response.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}

// ExistByUsername is function to check that username exist or not
func (u *UserRepositoryPostgreSQL) ExistByUsernameId(username string, id string) (exist bool, err error) {
	err = u.DB.Read.Get(&exist, userQuery.Exist+" WHERE username = $1 AND u.active is true AND u.is_deleted is false and  id <> $2 ", username, id)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return exist, err
}
