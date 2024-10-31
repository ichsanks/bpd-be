package auth

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
)

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}
type User struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	TenantID         *uuid.UUID `json:"tenantId" db:"tenant_id"`
	IDBranch         *string    `json:"idBranch" db:"id_branch"`
	Nama             *string    `json:"nama" db:"nama"`
	Username         string     `json:"username" db:"username"`
	Email            *string    `json:"email" db:"email"`
	Password         string     `json:"password" db:"password"`
	Status           *string    `json:"status" db:"status"`
	IDRole           *string    `json:"idRole" db:"id_role"`
	IDUnor           *string    `json:"idUnor" db:"id_unor"`
	IDPegawai        *string    `json:"idPegawai" db:"id_pegawai"`
	IDFungsionalitas *string    `json:"idFungsionalitas" db:"id_fungsionalitas"`
	Active           bool       `db:"active" json:"active"`
	Foto             *string    `json:"foto" db:"foto"`
	CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt        *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt        *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted        bool       `db:"is_deleted" json:"isDeleted"`
}

// UserUpdateFormat
type UserUpdateFormat struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Nama             *string   `json:"nama" db:"nama"`
	Username         string    `json:"username" db:"username"`
	Email            *string   `json:"email" db:"email"`
	Status           *string   `json:"status" db:"status"`
	IDRole           *string   `json:"idRole" db:"id_role"`
	IDUnor           *string   `json:"idUnor" db:"id_unor"`
	IDPegawai        *string   `json:"idPegawai" db:"id_pegawai"`
	IDFungsionalitas *string   `json:"idFungsionalitas" db:"id_fungsionalitas"`
	IDBranch         *string   `json:"idBranch" db:"id_branch"`
	Active           bool      `db:"active" json:"active"`
}

// UserUpdateFcmTokenFormat
type UserUpdateFcmTokenFormat struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	FirebaseToken null.String `json:"firebaseToken" db:"firebase_token"`
}

// UserDTO digunakan untuk model join ke Role
type UserDTO struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	TenantID           *uuid.UUID `json:"tenantId" db:"tenant_id"`
	IDBranch           *string    `json:"idBranch" db:"id_branch"`
	NamaBranch         *string    `json:"namaBranch" db:"nama_branch"`
	Nama               *string    `json:"nama" db:"nama"`
	Username           string     `json:"username" db:"username"`
	Email              *string    `json:"email" db:"email"`
	Password           string     `json:"password" db:"password"`
	Status             *string    `json:"status" db:"status"`
	Nip                *string    `json:"nip" db:"nip"`
	IDPegawai          *string    `json:"idPegawai" db:"id_pegawai"`
	NamaPegawai        *string    `json:"namaPegawai" db:"nama_pegawai"`
	IDRole             *string    `json:"idRole" db:"id_role"`
	NamaRole           *string    `json:"namaRole" db:"nama_role"`
	IDUnor             *string    `json:"idUnor" db:"id_unor"`
	NamaUnor           *string    `json:"namaUnor" db:"nama_unor"`
	IDFungsionalitas   *string    `json:"idFungsionalitas" db:"id_fungsionalitas"`
	NamaFungsionalitas *string    `json:"namaFungsionalitas" db:"nama_fungsionalitas"`
	IDBidang           *string    `json:"idBidang" db:"id_bidang"`
	NamaBidang         *string    `json:"namaBidang" db:"nama_bidang"`
	JenisApproval      *string    `json:"jenisApproval" db:"jenis_approval"`
	Active             bool       `db:"active" json:"active"`
}

// UserDTO digunakan untuk model join ke Role
type UserRoleDTO struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	Username      string      `json:"username" db:"username"`
	Email         string      `json:"email" db:"email"`
	Status        null.String `json:"status" db:"status"`
	FirebaseToken null.String `json:"firebaseToken" db:"firebase_token"`
	IsDeleted     bool        `json:"isDeleted" db:"is_deleted"`
	RoleID        null.String `json:"roleId" db:"role_id"`
	Role          null.String `json:"role" db:"name"`
}

type StatusLogin string

const (
	SuccessLogin StatusLogin = "success"
	FailedLogin  StatusLogin = "failed"
)

type LoginActivity struct {
	ID       uuid.UUID   `json:"id" db:"id"`
	Username string      `json:"username" db:"username"`
	Status   StatusLogin `json:"status" db:"status"`
	Jam      time.Time   `json:"jam" db:"jam"`
}

func NewCreateActivityLogin(username string, status StatusLogin) LoginActivity {
	loginActivityID, _ := uuid.NewV4()
	return LoginActivity{
		ID:       loginActivityID,
		Username: username,
		Status:   status,
		Jam:      time.Now(),
	}
}

// InputUser is struct as register json body
type InputUser struct {
	Nama             *string    `json:"nama" db:"nama"`
	Username         string     `json:"username" db:"username"`
	Email            *string    `json:"email" db:"email"`
	Password         string     `json:"password" db:"password"`
	Status           *string    `json:"status" db:"status"`
	IDRole           *string    `json:"idRole" db:"id_role"`
	IDUnor           *string    `json:"idUnor" db:"id_unor"`
	IDPegawai        *string    `json:"idPegawai" db:"id_pegawai"`
	IDFungsionalitas *string    `json:"idFungsionalitas" db:"id_fungsionalitas"`
	TenantID         *uuid.UUID `json:"-"`
	IDBranch         *string    `json:"-"`
	Active           bool       `db:"active" json:"active"`
}

// Validate digunakan untuk memvalidasi inputan user
func (i InputUser) Validate() error {
	v := shared.GetValidator()
	v.RegisterValidation("alphaspace", shared.AlphaSpace)
	v.RegisterValidation("alphanumspace", shared.AlphaNumSpace)

	return v.Struct(i)
}

// CreateUser is function to parse from user input to user struct
func (i InputUser) CreateUser() User {
	userID, _ := uuid.NewV4()

	hash, _ := bcrypt.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	return User{
		ID:               userID,
		IDRole:           i.IDRole,
		Username:         i.Username,
		Nama:             i.Nama,
		Status:           i.Status,
		Password:         string(hash),
		Email:            i.Email,
		IDPegawai:        i.IDPegawai,
		IDFungsionalitas: i.IDFungsionalitas,
		IDUnor:           i.IDUnor,
		Active:           i.Active,
		CreatedAt:        time.Now(),
		TenantID:         i.TenantID,
		IDBranch:         i.IDBranch,
	}
}

// CreateUser is function to parse from user input to user struct
func (i InputUser) Registrasi() User {
	userID, _ := uuid.NewV4()

	hash, _ := bcrypt.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	return User{
		ID:               userID,
		IDRole:           i.IDRole,
		Username:         i.Username,
		Password:         string(hash),
		IDPegawai:        i.IDPegawai,
		IDUnor:           i.IDUnor,
		IDFungsionalitas: i.IDFungsionalitas,
	}
}

type InputChangePassword struct {
	OldPassword string `json:"oldPassword" validate:"required,min=6"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

// Update is function to transform into to User entity
func (i InputChangePassword) Update(user User) (User, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(i.OldPassword))
	if err != nil {
		return User{}, failure.Conflict("update password", "password", "old password does not match with the current password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(i.NewPassword))
	if err == nil {
		return User{}, failure.Conflict("update password", "password", "new password match with the current password")
	}

	newPassword, _ := bcrypt.GenerateFromPassword([]byte(i.NewPassword), bcrypt.DefaultCost)
	now := time.Now()
	return User{
		ID:               user.ID,
		IDRole:           user.IDRole,
		Username:         user.Username,
		Password:         string(newPassword),
		Email:            user.Email,
		IDPegawai:        user.IDPegawai,
		IDFungsionalitas: user.IDFungsionalitas,
		IDUnor:           user.IDUnor,
		Active:           user.Active,
		IDBranch:         user.IDBranch,
		UpdatedAt:        &now,
	}, nil
}

// ResetPasswdUpdate is function to transform into to User entity
func (i InputChangePassword) ResetPasswdUpdate(user User) (User, error) {
	newPassword, _ := bcrypt.GenerateFromPassword([]byte(i.NewPassword), bcrypt.DefaultCost)
	return User{
		ID:        user.ID,
		IDRole:    user.IDRole,
		IDPegawai: user.IDPegawai,
		Username:  user.Username,
		Password:  string(newPassword),
	}, nil
}

// Update is function to transform into to User entity
func (i UserUpdateFormat) Update(user UserUpdateFormat) (User, error) {
	// newPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	now := time.Now()
	return User{
		ID:               user.ID,
		IDRole:           user.IDRole,
		Username:         user.Username,
		Nama:             user.Nama,
		Status:           user.Status,
		Email:            user.Email,
		IDPegawai:        user.IDPegawai,
		IDFungsionalitas: user.IDFungsionalitas,
		IDUnor:           user.IDUnor,
		Active:           user.Active,
		IDBranch:         user.IDBranch,
		UpdatedAt:        &now,
	}, nil
}

// Update is function to transform into to User entity untuk update token fcm
func (i UserUpdateFcmTokenFormat) UpdateFcmToken(user UserUpdateFcmTokenFormat) (User, error) {
	return User{
		ID: user.ID,
	}, nil
}

// InputLogin is struct as login json body
type InputLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	// CapchaTOken string `json:"capchaTOken"`
	// RoleID   string `json:"roleId"`
}

// Response is represent respond login
func (r *InputLogin) Response(user UserDTO, role Role, accessToken string) ResponseLogin {
	return ResponseLogin{
		Token: ResponseLoginToken{
			AccessToken: accessToken,
		},
		User: ResponseLoginUser{
			ID:                 user.ID,
			IDRole:             user.IDRole,
			Username:           user.Username,
			Nama:               user.Nama,
			Status:             user.Status,
			Email:              user.Email,
			Nip:                user.Nip,
			IDPegawai:          user.IDPegawai,
			IDFungsionalitas:   user.IDFungsionalitas,
			IDUnor:             user.IDUnor,
			IDBidang:           user.IDBidang,
			NamaUnor:           user.NamaUnor,
			NamaBidang:         user.NamaBidang,
			NamaFungsionalitas: user.NamaFungsionalitas,
			JenisApproval:      user.JenisApproval,
			Role:               role,
			TenantID:           user.TenantID,
			IDBranch:           user.IDBranch,
			NamaBranch:         user.NamaBranch,
		},
	}
}

// ResponseLogin is result processing from login process
type ResponseLogin struct {
	Token ResponseLoginToken `json:"token"`
	User  ResponseLoginUser  `json:"user"`
}

// ResponseLoginUser deliver result of user entity
type ResponseLoginUser struct {
	ID                 uuid.UUID   `json:"id"`
	TenantID           *uuid.UUID  `json:"tenantId" db:"tenant_id"`
	IDBranch           *string     `json:"idBranch" db:"id_branch"`
	NamaBranch         *string     `json:"namaBranch" db:"nama_branch"`
	Nama               *string     `json:"nama" db:"nama"`
	Username           string      `json:"username" db:"username"`
	Email              *string     `json:"email" db:"email"`
	Status             *string     `json:"status" db:"status"`
	IDRole             *string     `json:"idRole" db:"id_role"`
	Nip                *string     `json:"nip" db:"nip"`
	IDPegawai          *string     `json:"idPegawai" db:"id_pegawai"`
	IDUnor             *string     `json:"idUnor" db:"id_unor"`
	NamaUnor           *string     `json:"namaUnor" db:"nama_unor"`
	IDFungsionalitas   *string     `json:"idFungsionalitas" db:"id_fungsionalitas"`
	NamaFungsionalitas *string     `json:"namaFungsionalitas" db:"nama_fungsionalitas"`
	IDBidang           *string     `json:"idBidang" db:"id_bidang"`
	NamaBidang         *string     `json:"namaBidang" db:"nama_bidang"`
	JenisApproval      *string     `json:"jenisApproval" db:"jenis_approval"`
	FirebaseToken      null.String `json:"firebaseToken"`
	Role               Role        `json:"role"`
}

// ResponseLoginToken deliver result of user token
type ResponseLoginToken struct {
	AccessToken string
}

// NewUserLoginClaims digunakan untuk mengeset nilai dari JWT
func NewUserLoginClaims(user UserDTO, expiredIn int) jwt.MapClaims {
	claims := jwt.MapClaims{}
	claims["userId"] = user.ID
	claims["pegawaiId"] = user.IDPegawai
	claims["roleId"] = user.IDRole
	claims["unorId"] = user.IDUnor
	claims["fungsionalitasId"] = user.IDFungsionalitas
	claims["tenantId"] = user.TenantID
	claims["branchId"] = user.IDBranch
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Duration(expiredIn) * time.Hour).Unix()

	return claims
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// SoftDelete untuk mengeset flag isDeleted
func (u *User) SoftDelete(userID uuid.UUID) {
	now := time.Now()
	u.Active = false
	u.IsDeleted = true
	u.DeletedAt = &now
	u.Username = u.Username + " [Deleted at " + now.Format("2006-01-02 15:04:05.000000") + "]"
}

var ColumnMappUser = map[string]interface{}{
	"id":                 "u.id",
	"username":           "u.username",
	"nama":               "u.nama",
	"email":              "u.email",
	"namaUnor":           "muk.nama",
	"namaFungsionalitas": "mf.nama",
	"namaBidang":         "mb.nama",
	"namaRole":           "r.nama",
	"createdAt":          "u.created_at",
	"createdBy":          "u.created_by",
	"updatedAt":          "u.updated_at",
	"updatedBy":          "u.updated_by",
	"isDeleted":          "u.is_deleted",
}
