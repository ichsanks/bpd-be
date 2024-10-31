package master

import (
	"time"

	"github.com/gofrs/uuid"
)

// field untuk transaksi
type Pegawai struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	Nip              string     `db:"nip" json:"nip"`
	Nama             string     `db:"nama" json:"nama"`
	KodeJk           *string    `db:"jenis_kelamin" json:"kodeJk"`
	KodeAgama        *string    `db:"agama" json:"kodeAgama"`
	Alamat           *string    `db:"alamat" json:"alamat"`
	NoHp             *string    `db:"no_hp" json:"noHp"`
	Email            *string    `db:"email" json:"email"`
	IdUnor           *string    `db:"id_unor" json:"idUnor"`
	IdBidang         string     `db:"id_bidang" json:"idBidang"`
	IdJabatan        *string    `db:"id_jabatan" json:"idJabatan"`
	IdGolongan       *string    `db:"id_golongan" json:"idGolongan"`
	IdFungsionalitas *string    `db:"id_fungsionalitas" json:"idFungsionalitas"`
	Foto             *string    `db:"foto" json:"foto"`
	CreatedAt        *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy        *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt        *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy        *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted        bool       `db:"is_deleted" json:"isDeleted"`
	IdApprovalLine   *string    `db:"id_approval_line" json:"idApprovalLine"`
	IdManager        *string    `db:"id_manager" json:"idManager"`
	Nik              *string    `db:"nik" json:"nik"`
	FotoTtd          *string    `db:"foto_ttd" json:"fotoTtd"`
	TenantID         *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBranch         *string    `db:"id_branch" json:"idBranch"`
	IdJobGrade       *string    `db:"id_job_grade" json:"idJobGrade"`
	IdPersonGrade    *string    `db:"id_person_grade" json:"idPersonGrade"`
	IdLevelBod       *string    `db:"id_level_bod" json:"idLevelBod"`
	IdStatusPegawai  *string    `db:"id_status_pegawai" json:"idStatusPegawai"`
	IdStatusKontrak  *string    `db:"id_status_kontrak" json:"idStatusKontrak"`
	KodeVendor       *string    `db:"kode_vendor" json:"kodeVendor"`
}

type PegawaiDTO struct {
	ID                 uuid.UUID  `db:"id" json:"id"`
	Nip                string     `db:"nip" json:"nip"`
	Nama               string     `db:"nama" json:"nama"`
	KodeJk             *string    `db:"jenis_kelamin" json:"kodeJk"`
	NamaJk             *string    `db:"nama_jk" json:"namaJk"`
	KodeAgama          *string    `db:"agama" json:"kodeAgama"`
	NamaAgama          *string    `db:"nama_agama" json:"namaAgama"`
	Alamat             *string    `db:"alamat" json:"alamat"`
	NoHp               *string    `db:"no_hp" json:"noHp"`
	Email              *string    `db:"email" json:"email"`
	IdUnor             *string    `db:"id_unor" json:"idUnor"`
	KodeUnor           *string    `db:"kode_unor" json:"kodeUnor"`
	NamaUnor           *string    `db:"nama_unor" json:"namaUnor"`
	IdJabatan          *string    `db:"id_jabatan" json:"idJabatan"`
	NamaJabatan        *string    `db:"nama_jabatan" json:"namaJabatan"`
	IdGolongan         *string    `db:"id_golongan" json:"idGolongan"`
	NamaGolongan       *string    `db:"nama_golongan" json:"namaGolongan"`
	IdFungsionalitas   *string    `db:"id_fungsionalitas" json:"idFungsionalitas"`
	NamaFungsionalitas *string    `db:"nama_fungsionalitas" json:"namaFungsionalitas"`
	IdBidang           *string    `db:"id_bidang" json:"idBidang"`
	NamaBidang         *string    `db:"nama_bidang" json:"namaBidang"`
	CreatedAt          *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy          *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt          *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy          *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted          bool       `db:"is_deleted" json:"isDeleted"`
	IdApprovalLine     *string    `db:"id_approval_line" json:"idApprovalLine"`
	IdManager          *string    `db:"id_manager" json:"idManager"`
	Nik                *string    `db:"nik" json:"nik"`
	FotoTtd            *string    `db:"foto_ttd" json:"fotoTtd"`
	TenantID           *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBranch           *string    `db:"id_branch" json:"idBranch"`
	IdJobGrade         *string    `db:"id_job_grade" json:"idJobGrade"`
	IdPersonGrade      *string    `db:"id_person_grade" json:"idPersonGrade"`
	IdLevelBod         *string    `db:"id_level_bod" json:"idLevelBod"`
	IdStatusPegawai    *string    `db:"id_status_pegawai" json:"idStatusPegawai"`
	IdStatusKontrak    *string    `db:"id_status_kontrak" json:"idStatusKontrak"`
	KodeVendor         *string    `db:"kode_vendor" json:"kodeVendor"`
}

// field param di swagger
type PegawaiFormat struct {
	ID               uuid.UUID `db:"id" json:"id"`
	Nip              string    `db:"nip" json:"nip"`
	Nama             string    `db:"nama" json:"nama"`
	KodeJk           *string   `db:"jenis_kelamin" json:"kodeJk"`
	KodeAgama        *string   `db:"agama" json:"kodeAgama"`
	Alamat           *string   `db:"alamat" json:"alamat"`
	NoHp             *string   `db:"no_hp" json:"noHp"`
	Email            *string   `db:"email" json:"email"`
	IdUnor           *string   `db:"id_unor" json:"idUnor"`
	IdBidang         string    `db:"id_bidang" json:"idBidang"`
	IdJabatan        *string   `db:"id_jabatan" json:"idJabatan"`
	IdGolongan       *string   `db:"id_golongan" json:"idGolongan"`
	IdFungsionalitas *string   `db:"id_fungsionalitas" json:"idFungsionalitas"`
	Nik              *string   `db:"nik" json:"nik"`
	FotoTtd          *string   `db:"foto_ttd" json:"fotoTtd"`
	IdBranch         *string   `db:"id_branch" json:"idBranch"`
	IdJobGrade       *string   `db:"id_job_grade" json:"idJobGrade"`
	IdPersonGrade    *string   `db:"id_person_grade" json:"idPersonGrade"`
	IdLevelBod       *string   `db:"id_level_bod" json:"idLevelBod"`
	IdStatusPegawai  *string   `db:"id_status_pegawai" json:"idStatusPegawai"`
	IdStatusKontrak  *string   `db:"id_status_kontrak" json:"idStatusKontrak"`
	KodeVendor       *string   `db:"kode_vendor" json:"kodeVendor"`
}

// alis dari json ke db untuk sort table fe
var ColumnMappPegawai = map[string]interface{}{
	"id":                 "mp.id",
	"nip":                "mp.nip",
	"nama":               "mp.nama",
	"namaBidang":         "mb.nama",
	"namaUnor":           "muok.nama",
	"namaFungsionalitas": "mf.nama",
	"namaGolongan":       "mg.nama",
	"namaJabatan":        "mj.nama",
	"createdAt":          "mp.created_at",
	"createdBy":          "mp.created_by",
	"updatedAt":          "mp.updated_at",
	"updatedBy":          "mp.updated_by",
	"isDeleted":          "mp.is_deleted",
}

// field create dan update
func (p *Pegawai) PegawaiFormat(reqFormat PegawaiFormat, userId uuid.UUID, tenantId uuid.UUID) (newPegawai Pegawai, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newPegawai = Pegawai{
			ID:               newID,
			Nip:              reqFormat.Nip,
			Nama:             reqFormat.Nama,
			KodeJk:           reqFormat.KodeJk,
			KodeAgama:        reqFormat.KodeAgama,
			Alamat:           reqFormat.Alamat,
			NoHp:             reqFormat.NoHp,
			Email:            reqFormat.Email,
			IdUnor:           reqFormat.IdUnor,
			IdBidang:         reqFormat.IdBidang,
			IdJabatan:        reqFormat.IdJabatan,
			IdGolongan:       reqFormat.IdGolongan,
			IdFungsionalitas: reqFormat.IdFungsionalitas,
			Nik:              reqFormat.Nik,
			FotoTtd:          reqFormat.FotoTtd,
			CreatedAt:        &now,
			CreatedBy:        &userId,
			TenantID:         &tenantId,
			IdBranch:         reqFormat.IdBranch,
			IdJobGrade:       reqFormat.IdJobGrade,
			IdPersonGrade:    reqFormat.IdPersonGrade,
			IdLevelBod:       reqFormat.IdLevelBod,
			IdStatusPegawai:  reqFormat.IdStatusPegawai,
			IdStatusKontrak:  reqFormat.IdStatusKontrak,
			KodeVendor:       reqFormat.KodeVendor,
		}
	} else {
		newPegawai = Pegawai{
			ID:               reqFormat.ID,
			Nip:              reqFormat.Nip,
			Nama:             reqFormat.Nama,
			KodeJk:           reqFormat.KodeJk,
			KodeAgama:        reqFormat.KodeAgama,
			Alamat:           reqFormat.Alamat,
			NoHp:             reqFormat.NoHp,
			Email:            reqFormat.Email,
			IdUnor:           reqFormat.IdUnor,
			IdBidang:         reqFormat.IdBidang,
			IdJabatan:        reqFormat.IdJabatan,
			IdGolongan:       reqFormat.IdGolongan,
			IdFungsionalitas: reqFormat.IdFungsionalitas,
			Nik:              reqFormat.Nik,
			FotoTtd:          reqFormat.FotoTtd,
			UpdatedAt:        &now,
			UpdatedBy:        &userId,
			TenantID:         &tenantId,
			IdBranch:         reqFormat.IdBranch,
			IdJobGrade:       reqFormat.IdJobGrade,
			IdPersonGrade:    reqFormat.IdPersonGrade,
			IdLevelBod:       reqFormat.IdLevelBod,
			IdStatusPegawai:  reqFormat.IdStatusPegawai,
			IdStatusKontrak:  reqFormat.IdStatusKontrak,
			KodeVendor:       reqFormat.KodeVendor,
		}
	}
	return
}

// field delete soft
func (p *Pegawai) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	p.IsDeleted = true
	p.Nama = p.Nama + " [Deleted at " + cTime.Format("2006-01-02 15:04:05") + "]"
}

type PegawaiParams struct {
	IdPegawai        string `db:"id_pegawai" json:"idPegawai"`
	IdFungsionalitas string `db:"id_fungsionalitas" json:"idFungsionalitas"`
	IdBidang         string `db:"id_bidang" json:"idBidang"`
	IdUnor           string `db:"id_unor" json:"idUnor"`
}
