//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/auth"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/bpd"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/master"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/report"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/handlers"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/router"
)

// Wiring for configurations.
var configurations = wire.NewSet(
	configs.Get,
)

// Wiring for persistences.
var persistences = wire.NewSet(
	infras.ProvidePostgreSQLConn,
)

// Wiring for all domains.
var domains = wire.NewSet(
	domainAuth,
	domainMaster,
	domainBpd,
	domainReport,
)

// Wiring for domain Auth
var domainAuth = wire.NewSet(
	// Log System interface and implementation
	auth.ProvideLogSystemServiceImpl,
	wire.Bind(new(auth.LogSystemService), new(*auth.LogSystemServiceImpl)),
	// LogSystemRepository interface and implementation
	auth.ProvideLogSystemRepositoryPostgreSQL,
	wire.Bind(new(auth.LogSystemRepository), new(*auth.LogSystemRepositoryPostgreSQL)),

	// Menu interface and implementation
	auth.ProvideMenuServiceImpl,
	wire.Bind(new(auth.MenuService), new(*auth.MenuServiceImpl)),
	// MenuRepository interface and implementation
	auth.ProvideMenuRepositoryPostgreSQL,
	wire.Bind(new(auth.MenuRepository), new(*auth.MenuRepositoryPostgreSQL)),

	// Role interface and implementation
	auth.ProvideRoleServiceImpl,
	wire.Bind(new(auth.RoleService), new(*auth.RoleServiceImpl)),
	// RoleRepository interface and implementation
	auth.ProvideRoleRepositoryPostgreSQL,
	wire.Bind(new(auth.RoleRepository), new(*auth.RoleRepositoryPostgreSQL)),

	// UserService interface and implementation
	auth.ProvideUserServiceImpl,
	wire.Bind(new(auth.UserService), new(*auth.UserServiceImpl)),
	// UserRepository interface and implementation
	auth.ProvideUserRepositoryPostgreSQL,
	wire.Bind(new(auth.UserRepository), new(*auth.UserRepositoryPostgreSQL)),

	// DashboardService interface and implementation
	auth.ProvideDashboardServiceImpl,
	wire.Bind(new(auth.DashboardService), new(*auth.DashboardServiceImpl)),
	// DashboardRepository interface and implementation
	auth.ProvideDashboardRepositoryPostgreSQL,
	wire.Bind(new(auth.DashboardRepository), new(*auth.DashboardRepositoryPostgreSQL)),
)

// Wiring for domain Master
var domainMaster = wire.NewSet(
	// LogService interface and implementation
	master.ProvideLogServiceImpl,
	wire.Bind(new(master.LogService), new(*master.LogServiceImpl)),
	// LogRepository interface and implementation
	master.ProvideLogRepositoryPostgreSQL,
	wire.Bind(new(master.LogRepository), new(*master.LogRepositoryPostgreSQL)),

	// Pegawai interface and implementation
	master.ProvidePegawaiServiceImpl,
	wire.Bind(new(master.PegawaiService), new(*master.PegawaiServiceImpl)),
	// PegawaiRepository interface and implementation
	master.ProvidePegawaiRepositoryPostgreSQL,
	wire.Bind(new(master.PegawaiRepository), new(*master.PegawaiRepositoryPostgreSQL)),

	// Agama interface and implementation
	master.ProvideAgamaServiceImpl,
	wire.Bind(new(master.AgamaService), new(*master.AgamaServiceImpl)),
	// AgamaRepository interface and implementation
	master.ProvideAgamaRepositoryPostgreSQL,
	wire.Bind(new(master.AgamaRepository), new(*master.AgamaRepositoryPostgreSQL)),

	// JenisKelamin interface and implementation
	master.ProvideJenisKelaminServiceImpl,
	wire.Bind(new(master.JenisKelaminService), new(*master.JenisKelaminServiceImpl)),
	// JenisKelaminRepository interface and implementation
	master.ProvideJenisKelaminRepositoryPostgreSQL,
	wire.Bind(new(master.JenisKelaminRepository), new(*master.JenisKelaminRepositoryPostgreSQL)),

	// UnitKerja interface and implementation
	master.ProvideUnitKerjaServiceImpl,
	wire.Bind(new(master.UnitKerjaService), new(*master.UnitKerjaServiceImpl)),
	// UnitKerjaRepository interface and implementation
	master.ProvideUnitKerjaRepositoryPostgreSQL,
	wire.Bind(new(master.UnitKerjaRepository), new(*master.UnitKerjaRepositoryPostgreSQL)),

	// Jabatan interface and implementation
	master.ProvideJabatanServiceImpl,
	wire.Bind(new(master.JabatanService), new(*master.JabatanServiceImpl)),
	// JabatanRepository interface and implementation
	master.ProvideJabatanRepositoryPostgreSQL,
	wire.Bind(new(master.JabatanRepository), new(*master.JabatanRepositoryPostgreSQL)),

	// Bidang interface and implementation
	master.ProvideBidangServiceImpl,
	wire.Bind(new(master.BidangService), new(*master.BidangServiceImpl)),
	// BidangRepository interface and implementation
	master.ProvideBidangRepositoryPostgreSQL,
	wire.Bind(new(master.BidangRepository), new(*master.BidangRepositoryPostgreSQL)),

	// Golongan interface and implementation
	master.ProvideGolonganServiceImpl,
	wire.Bind(new(master.GolonganService), new(*master.GolonganServiceImpl)),
	// GolonganRepository interface and implementation
	master.ProvideGolonganRepositoryPostgreSQL,
	wire.Bind(new(master.GolonganRepository), new(*master.GolonganRepositoryPostgreSQL)),

	// Fungsionalitas interface and implementation
	master.ProvideFungsionalitasServiceImpl,
	wire.Bind(new(master.FungsionalitasService), new(*master.FungsionalitasServiceImpl)),
	// FungsionalitasRepository interface and implementation
	master.ProvideFungsionalitasRepositoryPostgreSQL,
	wire.Bind(new(master.FungsionalitasRepository), new(*master.FungsionalitasRepositoryPostgreSQL)),

	// Kendaraan interface and implementation
	master.ProvideKendaraanServiceImpl,
	wire.Bind(new(master.KendaraanService), new(*master.KendaraanServiceImpl)),
	// KendaraanRepository interface and implementation
	master.ProvideKendaraanRepositoryPostgreSQL,
	wire.Bind(new(master.KendaraanRepository), new(*master.KendaraanRepositoryPostgreSQL)),

	// JenisBiaya interface and implementation
	master.ProvideJenisBiayaServiceImpl,
	wire.Bind(new(master.JenisBiayaService), new(*master.JenisBiayaServiceImpl)),
	// JenisBiayaRepository interface and implementation
	master.ProvideJenisBiayaRepositoryPostgreSQL,
	wire.Bind(new(master.JenisBiayaRepository), new(*master.JenisBiayaRepositoryPostgreSQL)),

	// JenisKendaraan interface and implementation
	master.ProvideJenisKendaraanServiceImpl,
	wire.Bind(new(master.JenisKendaraanService), new(*master.JenisKendaraanServiceImpl)),
	// JenisKendaraanRepository interface and implementation
	master.ProvideJenisKendaraanRepositoryPostgreSQL,
	wire.Bind(new(master.JenisKendaraanRepository), new(*master.JenisKendaraanRepositoryPostgreSQL)),

	// JenisPerjalananDinas interface and implementation
	master.ProvideJenisPerjalananDinasServiceImpl,
	wire.Bind(new(master.JenisPerjalananDinasService), new(*master.JenisPerjalananDinasServiceImpl)),
	// JenisPerjalananDinasRepository interface and implementation
	master.ProvideJenisPerjalananDinasRepositoryPostgreSQL,
	wire.Bind(new(master.JenisPerjalananDinasRepository), new(*master.JenisPerjalananDinasRepositoryPostgreSQL)),

	// RuleApproval interface and implementation
	master.ProvideRuleApprovalServiceImpl,
	wire.Bind(new(master.RuleApprovalService), new(*master.RuleApprovalServiceImpl)),
	// RuleApprovalRepository interface and implementation
	master.ProvideRuleApprovalRepositoryPostgreSQL,
	wire.Bind(new(master.RuleApprovalRepository), new(*master.RuleApprovalRepositoryPostgreSQL)),

	// JenisApproval interface and implementation
	master.ProvideJenisApprovalServiceImpl,
	wire.Bind(new(master.JenisApprovalService), new(*master.JenisApprovalServiceImpl)),
	// JenisApprovalRepository interface and implementation
	master.ProvideJenisApprovalRepositoryPostgreSQL,
	wire.Bind(new(master.JenisApprovalRepository), new(*master.JenisApprovalRepositoryPostgreSQL)),

	// StatusPegawai interface and implementation
	master.ProvideStatusPegawaiServiceImpl,
	wire.Bind(new(master.StatusPegawaiService), new(*master.StatusPegawaiServiceImpl)),
	// StatusPegawaiRepository interface and implementation
	master.ProvideStatusPegawaiRepositoryPostgreSQL,
	wire.Bind(new(master.StatusPegawaiRepository), new(*master.StatusPegawaiRepositoryPostgreSQL)),

	// JobGrade interface and implementation
	master.ProvideJobGradeServiceImpl,
	wire.Bind(new(master.JobGradeService), new(*master.JobGradeServiceImpl)),
	// StatusPegawaiRepository interface and implementation
	master.ProvideJobGradeRepositoryPostgreSQL,
	wire.Bind(new(master.JobGradeRepository), new(*master.JobGradeRepositoryPostgreSQL)),

	// PersonGrade interface and implementation
	master.ProvidePersonGradeServiceImpl,
	wire.Bind(new(master.PersonGradeService), new(*master.PersonGradeServiceImpl)),
	// StatusPegawaiRepository interface and implementation
	master.ProvidePersonGradeRepositoryPostgreSQL,
	wire.Bind(new(master.PersonGradeRepository), new(*master.PersonGradeRepositoryPostgreSQL)),

	// LevelBod interface and implementation
	master.ProvideLevelBodServiceImpl,
	wire.Bind(new(master.LevelBodService), new(*master.LevelBodServiceImpl)),
	// LevelBodRepository interface and implementation
	master.ProvideLevelBodRepositoryPostgreSQL,
	wire.Bind(new(master.LevelBodRepository), new(*master.LevelBodRepositoryPostgreSQL)),

	// JenisTujuan interface and implementation
	master.ProvideJenisTujuanServiceImpl,
	wire.Bind(new(master.JenisTujuanService), new(*master.JenisTujuanServiceImpl)),
	// JenisTujuanRepository interface and implementation
	master.ProvideJenisTujuanRepositoryPostgreSQL,
	wire.Bind(new(master.JenisTujuanRepository), new(*master.JenisTujuanRepositoryPostgreSQL)),

	// FasilitasTransport interface and implementation
	master.ProvideFasilitasTransportServiceImpl,
	wire.Bind(new(master.FasilitasTransportService), new(*master.FasilitasTransportServiceImpl)),
	// FasilitasTransportRepository interface and implementation
	master.ProvideFasilitasTransportRepositoryPostgreSQL,
	wire.Bind(new(master.FasilitasTransportRepository), new(*master.FasilitasTransportRepositoryPostgreSQL)),

	// Branch interface and implementation
	master.ProvideBranchServiceImpl,
	wire.Bind(new(master.BranchService), new(*master.BranchServiceImpl)),
	// BranchRepository interface and implementation
	master.ProvideBranchRepositoryPostgreSQL,
	wire.Bind(new(master.BranchRepository), new(*master.BranchRepositoryPostgreSQL)),

	// Dokumen interface and implementation
	master.ProvideDokumenServiceImpl,
	wire.Bind(new(master.DokumenService), new(*master.DokumenServiceImpl)),
	// DokumenRepository interface and implementation
	master.ProvideDokumenRepositoryPostgreSQL,
	wire.Bind(new(master.DokumenRepository), new(*master.DokumenRepositoryPostgreSQL)),

	// SyaratDokumen interface and implementation
	master.ProvideSyaratDokumenServiceImpl,
	wire.Bind(new(master.SyaratDokumenService), new(*master.SyaratDokumenServiceImpl)),
	// SyaratDokumenRepository interface and implementation
	master.ProvideSyaratDokumenRepositoryPostgreSQL,
	wire.Bind(new(master.SyaratDokumenRepository), new(*master.SyaratDokumenRepositoryPostgreSQL)),

	// StatusKontrak interface and implementation
	master.ProvideStatusKontrakServiceImpl,
	wire.Bind(new(master.StatusKontrakService), new(*master.StatusKontrakServiceImpl)),
	// StatusKontrakRepository interface and implementation
	master.ProvideStatusKontrakRepositoryPostgreSQL,
	wire.Bind(new(master.StatusKontrakRepository), new(*master.StatusKontrakRepositoryPostgreSQL)),

	// SettingBiaya interface and implementation
	master.ProvideSettingBiayaServiceImpl,
	wire.Bind(new(master.SettingBiayaService), new(*master.SettingBiayaServiceImpl)),
	// SettingBiayaRepository interface and implementation
	master.ProvideSettingBiayaRepositoryPostgreSQL,
	wire.Bind(new(master.SettingBiayaRepository), new(*master.SettingBiayaRepositoryPostgreSQL)),

	// KategoriBiaya interface and implementation
	master.ProvideKategoriBiayaServiceImpl,
	wire.Bind(new(master.KategoriBiayaService), new(*master.KategoriBiayaServiceImpl)),
	// KategoriBiayaRepository interface and implementation
	master.ProvideKategoriBiayaRepositoryPostgreSQL,
	wire.Bind(new(master.KategoriBiayaRepository), new(*master.KategoriBiayaRepositoryPostgreSQL)),

	// STtd interface and implementation
	master.ProvideSTtdServiceImpl,
	wire.Bind(new(master.STtdService), new(*master.STtdServiceImpl)),
	// STtdRepository interface and implementation
	master.ProvideSTtdRepositoryPostgreSQL,
	wire.Bind(new(master.STtdRepository), new(*master.STtdRepositoryPostgreSQL)),

	// JenisSppd interface and implementation
	master.ProvideJenisSppdServiceImpl,
	wire.Bind(new(master.JenisSppdService), new(*master.JenisSppdServiceImpl)),
	// JenisSppdRepository interface and implementation
	master.ProvideJenisSppdRepositoryPostgreSQL,
	wire.Bind(new(master.JenisSppdRepository), new(*master.JenisSppdRepositoryPostgreSQL)),
)

var domainBpd = wire.NewSet(
	// PerjalananDinas interface and implementation
	bpd.ProvidePerjalananDinasServiceImpl,
	wire.Bind(new(bpd.PerjalananDinasService), new(*bpd.PerjalananDinasServiceImpl)),
	// PerjalananDinasRepository interface and implementation
	bpd.ProvidePerjalananDinasRepositoryPostgreSQL,
	wire.Bind(new(bpd.PerjalananDinasRepository), new(*bpd.PerjalananDinasRepositoryPostgreSQL)),

	// PerjalananDinasBiaya interface and implementation
	bpd.ProvidePerjalananDinasBiayaServiceImpl,
	wire.Bind(new(bpd.PerjalananDinasBiayaService), new(*bpd.PerjalananDinasBiayaServiceImpl)),
	// PerjalananDinasBiayaRepository interface and implementation
	bpd.ProvidePerjalananDinasBiayaRepositoryPostgreSQL,
	wire.Bind(new(bpd.PerjalananDinasBiayaRepository), new(*bpd.PerjalananDinasBiayaRepositoryPostgreSQL)),

	// PengajuanBpdHistori interface and implementation
	bpd.ProvidePengajuanBpdHistoriServiceImpl,
	wire.Bind(new(bpd.PengajuanBpdHistoriService), new(*bpd.PengajuanBpdHistoriServiceImpl)),
	// PengajuanBpdHistoriRepository interface and implementation
	bpd.ProvidePengajuanBpdHistoriRepositoryPostgreSQL,
	wire.Bind(new(bpd.PengajuanBpdHistoriRepository), new(*bpd.PengajuanBpdHistoriRepositoryPostgreSQL)),

	// PerjalananDinasKendaraan interface and implementation
	bpd.ProvidePerjalananDinasKendaraanServiceImpl,
	wire.Bind(new(bpd.PerjalananDinasKendaraanService), new(*bpd.PerjalananDinasKendaraanServiceImpl)),
	// PerjalananDinasKendaraanRepository interface and implementation
	bpd.ProvidePerjalananDinasKendaraanRepositoryPostgreSQL,
	wire.Bind(new(bpd.PerjalananDinasKendaraanRepository), new(*bpd.PerjalananDinasKendaraanRepositoryPostgreSQL)),

	// LogKegiatan interface and implementation
	bpd.ProvideLogKegiatanServiceImpl,
	wire.Bind(new(bpd.LogKegiatanService), new(*bpd.LogKegiatanServiceImpl)),
	// LogKegiatanRepository interface and implementation
	bpd.ProvideLogKegiatanRepositoryPostgreSQL,
	wire.Bind(new(bpd.LogKegiatanRepository), new(*bpd.LogKegiatanRepositoryPostgreSQL)),
	// Superman
	// SppSuperman interface and implementation
	bpd.ProvideSppSupermanServiceImpl,
	wire.Bind(new(bpd.SppSupermanService), new(*bpd.SppSupermanServiceImpl)),
	// SppSupermanRepository interface and implementation
	bpd.ProvideSppSupermanRepositoryPostgreSQL,
	wire.Bind(new(bpd.SppSupermanRepository), new(*bpd.SppSupermanRepositoryPostgreSQL)),

	// Surat Perjalanan Dinas
	// Surat Perjalanan Dinas interface and implementation
	bpd.ProvideSuratPerjalananDinasServiceImpl,
	wire.Bind(new(bpd.SuratPerjalananDinasService), new(*bpd.SuratPerjalananDinasServiceImpl)),
	// SuratPerjalananDinasRepository interface and implementation
	bpd.ProvideSuratPerjalananDinasRepositoryPostgreSQL,
	wire.Bind(new(bpd.SuratPerjalananDinasRepository), new(*bpd.SuratPerjalananDinasRepositoryPostgreSQL)),

	// Sppd Dokumen interface and implementation
	bpd.ProvideSppdDokumenServiceImpl,
	wire.Bind(new(bpd.SppdDokumenService), new(*bpd.SppdDokumenServiceImpl)),
	// SppdDokumenRepository interface and implementation
	bpd.ProvideSppdDokumenRepositoryPostgreSQL,
	wire.Bind(new(bpd.SppdDokumenRepository), new(*bpd.SppdDokumenRepositoryPostgreSQL)),

	// PengajuanSppdHistori interface and implementation
	bpd.ProvidePengajuanSppdHistoriServiceImpl,
	wire.Bind(new(bpd.PengajuanSppdHistoriService), new(*bpd.PengajuanSppdHistoriServiceImpl)),
	// PengajuanBpdHistoriRepository interface and implementation
	bpd.ProvidePengajuanSppdHistoriRepositoryPostgreSQL,
	wire.Bind(new(bpd.PengajuanSppdHistoriRepository), new(*bpd.PengajuanSppdHistoriRepositoryPostgreSQL)),
)

var domainReport = wire.NewSet(
	// Report interface and implementation
	report.ProvideReportServiceImpl,
	wire.Bind(new(report.ReportService), new(*report.ReportServiceImpl)),
	// ReportRepository interface and implementation
	report.ProvideReportRepositoryPostgreSQL,
	wire.Bind(new(report.ReportRepository), new(*report.ReportRepositoryPostgreSQL)),
)

// Wiring for HTTP routing.
var routing = wire.NewSet(
	wire.Struct(new(router.DomainHandlers), "*"),
	// Auth
	handlers.ProvideLogSystemHandler,
	handlers.ProvideMenuHandler,
	handlers.ProvideRoleHandler,
	handlers.ProvideUserHandler,
	handlers.ProvideDashboardHandler,
	// Master
	handlers.ProvideLogHandler,
	handlers.ProvidePegawaiHandler,
	handlers.ProvideAgamaHandler,
	handlers.ProvideJenisKelaminHandler,
	handlers.ProvideUnitKerjaHandler,
	handlers.ProvideJabatanHandler,
	handlers.ProvideBidangHandler,
	handlers.ProvideGolonganHandler,
	handlers.ProvideFungsionalitasHandler,
	handlers.ProvideKendaraanHandler,
	handlers.ProvideJenisBiayaHandler,
	handlers.ProvideJenisKendaraanHandler,
	handlers.ProvideJenisPerjalananDinasHandler,
	handlers.ProvideRuleApprovalHandler,
	handlers.ProvideJenisApprovalHandler,
	handlers.ProvideStatusPegawaiHandler,
	handlers.ProvideJobGradeHandler,
	handlers.ProvidePersonGradeHandler,
	handlers.ProvideLevelBodHandler,
	handlers.ProvideJenisTujuanHandler,
	handlers.ProvideFasilitasTransportHandler,
	handlers.ProvideBranchHandler,
	handlers.ProvideDokumenHandler,
	handlers.ProvideSyaratDokumenHandler,
	handlers.ProvideStatusKontrakHandler,
	handlers.ProvideSettingBiayaHandler,
	handlers.ProvideKategoriBiayaHandler,
	handlers.ProvideSTtdHandler,
	handlers.ProvideJenisSppdHandler,
	//Bpd
	handlers.ProvidePerjalananDinasHandler,
	handlers.ProvidePerjalananDinasBiayaHandler,
	handlers.ProvidePengajuanBpdHistoriHandler,
	handlers.ProvidePerjalananDinasKendaraanHandler,
	handlers.ProvideLogKegiatanHandler,
	handlers.ProvideSppSupermanHandler,
	handlers.ProvideSuratPerjalananDinasHandler,
	handlers.ProvideSppdDokumenHandler,
	handlers.ProvidePengajuanSppdHistoriHandler,
	// Report
	handlers.ProvideReportHandler,
	// File
	handlers.ProvideFileHandler,
	// JWT
	middleware.ProvideJWTMiddleware,
	router.ProvideRouter,
)

// Wiring for everything.
func InitializeService() *http.HTTP {
	wire.Build(
		// configurations
		configurations,
		// persistences
		persistences,
		// domains
		domains,
		// routing
		routing,
		// selected transport layer
		http.ProvideHTTP)
	return &http.HTTP{}
}
