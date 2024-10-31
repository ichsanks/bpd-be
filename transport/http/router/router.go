package router

import (
	"github.com/go-chi/chi"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/handlers"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
)

// DomainHandlers is a struct that contains all domain-specific handlers.
type DomainHandlers struct {
	// Auth
	LogSystemHandler handlers.LogSystemHandler
	MenuHandler      handlers.MenuHandler
	RoleHandler      handlers.RoleHandler
	UserHandler      handlers.UserHandler
	DashboardHandler handlers.DashboardHandler
	// Master
	LogHandler                  handlers.LogHandler
	PegawaiHandler              handlers.PegawaiHandler
	AgamaHandler                handlers.AgamaHandler
	JenisKelaminHandler         handlers.JenisKelaminHandler
	JabatanHandler              handlers.JabatanHandler
	BidangHandler               handlers.BidangHandler
	GolonganHandler             handlers.GolonganHandler
	UnitKerjaHandler            handlers.UnitKerjaHandler
	FungsionalitasHandler       handlers.FungsionalitasHandler
	KendaraanHandler            handlers.KendaraanHandler
	JenisBiayaHandler           handlers.JenisBiayaHandler
	JenisKendaraanHandler       handlers.JenisKendaraanHandler
	JenisPerjalananDinasHandler handlers.JenisPerjalananDinasHandler
	RuleApprovalHandler         handlers.RuleApprovalHandler
	StatusPegawaiHandler        handlers.StatusPegawaiHandler
	JobGradeHandler             handlers.JobGradeHandler
	PersonGradeHandler          handlers.PersonGradeHandler
	LevelBodHandler             handlers.LevelBodHandler
	JenisTujuanHandler          handlers.JenisTujuanHandler
	FasilitasTransportHandler   handlers.FasilitasTransportHandler
	BranchHandler               handlers.BranchHandler
	DokumenHandler              handlers.DokumenHandler
	SyaratDokumenHandler        handlers.SyaratDokumenHandler
	StatusKontrakHandler        handlers.StatusKontrakHandler
	SettingBiayaHandler         handlers.SettingBiayaHandler
	KategoriBiayaHandler        handlers.KategoriBiayaHandler
	STtdHandler                 handlers.STtdHandler
	JenisSppdHandler            handlers.JenisSppdHandler
	//BPD
	PerjalananDinasHandler          handlers.PerjalananDinasHandler
	PerjalananDinasBiayaHandler     handlers.PerjalananDinasBiayaHandler
	PengajuanBpdHistoriHandler      handlers.PengajuanBpdHistoriHandler
	PerjalananDinasKendaraanHandler handlers.PerjalananDinasKendaraanHandler
	LogKegiatanHandler              handlers.LogKegiatanHandler
	JenisApprovalHandler            handlers.JenisApprovalHandler
	SppSupermanHandler              handlers.SppSupermanHandler
	SuratPerjalananDinasHandler     handlers.SuratPerjalananDinasHandler
	SppdDokumenHandler              handlers.SppdDokumenHandler
	PengajuanSppdHistoriHandler     handlers.PengajuanSppdHistoriHandler
	// Report
	ReportHandler handlers.ReportHandler
	// File
	FileHandler handlers.FileHandler
}

// Router is the router struct containing handlers.
type Router struct {
	JwtMiddleware  *middleware.JWT
	DomainHandlers DomainHandlers
}

// ProvideRouter is the provider function for this router.
func ProvideRouter(domainHandlers DomainHandlers, jwtMiddleware *middleware.JWT) Router {
	return Router{
		DomainHandlers: domainHandlers,
		JwtMiddleware:  jwtMiddleware,
	}
}

// SetupRoutes sets up all routing for this server.
func (r *Router) SetupRoutes(mux *chi.Mux) {
	mux.Route("/v1", func(rc chi.Router) {
		// Auth
		r.DomainHandlers.LogSystemHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.MenuHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.RoleHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.UserHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.DashboardHandler.Router(rc, r.JwtMiddleware)
		// Master
		r.DomainHandlers.LogHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.PegawaiHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.AgamaHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.JenisKelaminHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.UnitKerjaHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.JabatanHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.BidangHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.GolonganHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.FungsionalitasHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.KendaraanHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.JenisBiayaHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.JenisKendaraanHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.JenisPerjalananDinasHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.RuleApprovalHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.JenisApprovalHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.StatusPegawaiHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.JobGradeHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.PersonGradeHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.LevelBodHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.JenisTujuanHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.FasilitasTransportHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.BranchHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.DokumenHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.SyaratDokumenHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.StatusKontrakHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.SettingBiayaHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.KategoriBiayaHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.STtdHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.JenisSppdHandler.Router(rc, r.JwtMiddleware)
		// BPD
		r.DomainHandlers.PerjalananDinasHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.PerjalananDinasBiayaHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.PengajuanBpdHistoriHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.PerjalananDinasKendaraanHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.LogKegiatanHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.SppSupermanHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.SuratPerjalananDinasHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.SppdDokumenHandler.Router(rc, r.JwtMiddleware)
		r.DomainHandlers.PengajuanSppdHistoriHandler.Router(rc, r.JwtMiddleware)
		// Report
		r.DomainHandlers.ReportHandler.Router(rc, r.JwtMiddleware)
		// File
		r.DomainHandlers.FileHandler.Router(rc, r.JwtMiddleware)
	})
}
