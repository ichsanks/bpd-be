package report

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/xuri/excelize/v2"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/master"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
)

type ReportService interface {
	RptRekapBpd(req FilterReport) (data []ReportRekapBpd, err error)
	RptRekapBpdBagian(req FilterReport) (data []ReportRekapBpdBagian, err error)
	RptRekapTotalBpd(req FilterReport) (data []ReportRekapTotalBpd, err error)
	ExportRekapBpd(req FilterReport) (xlsx *excelize.File, err error)
	ExportRekapBpdBagian(req FilterReport) (xlsx *excelize.File, err error)
	ExportBiayaAkomodasiDetail(req FilterReport) (xlsx *excelize.File, err error)
	RptRekapReimbusment(req FilterReport) (data []ReportRekapAkReim, err error)
	RptRekapAkomodasi(req FilterReport) (data []RekapBiayaAkomodasi, err error)
	ExportRekapReimAkm(req FilterReport) (xlsx *excelize.File, err error)
}

type ReportServiceImpl struct {
	ReportRepository ReportRepository
	BidangRepository master.BidangRepository
	Config           *configs.Config
}

func ProvideReportServiceImpl(repository ReportRepository, bidangRepo master.BidangRepository, config *configs.Config) *ReportServiceImpl {
	s := new(ReportServiceImpl)
	s.ReportRepository = repository
	s.BidangRepository = bidangRepo
	s.Config = config
	return s
}

func (s *ReportServiceImpl) RptRekapBpd(req FilterReport) (data []ReportRekapBpd, err error) {
	return s.ReportRepository.RptRekapBpd(req)
}

func (s *ReportServiceImpl) RptRekapBpdBagian(req FilterReport) (data []ReportRekapBpdBagian, err error) {
	return s.ReportRepository.RptRekapBpdBagian(req)
}

func (s *ReportServiceImpl) RptRekapTotalBpd(req FilterReport) (data []ReportRekapTotalBpd, err error) {
	return s.ReportRepository.RptRekapTotalBpd(req)
}

func (s *ReportServiceImpl) ExportRekapBpd(req FilterReport) (xlsx *excelize.File, err error) {
	report, err := s.ReportRepository.RptRekapBpd(req)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Report not found")
	}

	namaBidang := "SEMUA BAGIAN"
	if req.IdBidang != "" {
		bidangID, _ := uuid.FromString(req.IdBidang)
		bidang, err := s.BidangRepository.ResolveByID(bidangID)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Bidang not found")
		}

		namaBidang = "BAGIAN " + bidang.Nama
	}

	d1, err := time.Parse(model.DefaultDateFormat, req.TglAwal)
	if err != nil {
		fmt.Println(err)
	}
	date1 := d1.Format("02/01/2006")

	d2, err := time.Parse(model.DefaultDateFormat, req.TglAkhir)
	if err != nil {
		fmt.Println(err)
	}
	date2 := d2.Format("02/01/2006")
	ketTgl := "TANGGAL " + date1 + " s.d " + date2

	xlsx = excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet1Name := "Rekap BPD"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)
	index, err := xlsx.NewSheet(sheet1Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Style
	exp := "#,##0"
	SHeader, _ := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
	})

	SHeader3, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	SChild, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	SChildNum, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &exp,
	})

	colStart := 6
	xlsx.SetCellValue(sheet1Name, "A1", "LAPORAN REKAP PERJALANAN DINAS")
	xlsx.SetCellValue(sheet1Name, "A2", namaBidang)
	xlsx.SetCellValue(sheet1Name, "A3", ketTgl)

	xlsx.SetCellValue(sheet1Name, "A5", "No")
	xlsx.SetCellValue(sheet1Name, "B5", "Nomor BPD")
	xlsx.SetCellValue(sheet1Name, "C5", "Tgl Berangkat")
	xlsx.SetCellValue(sheet1Name, "D5", "Tgl Kembali")
	xlsx.SetCellValue(sheet1Name, "E5", "Nama")
	xlsx.SetCellValue(sheet1Name, "F5", "Bagian")
	xlsx.SetCellValue(sheet1Name, "G5", "Jabatan")
	xlsx.SetCellValue(sheet1Name, "H5", "Total Akomodasi")
	xlsx.SetCellValue(sheet1Name, "I5", "Total Deklarasi")

	// Style
	xlsx.SetCellStyle(sheet1Name, "A1", "A1", SHeader)
	xlsx.SetCellStyle(sheet1Name, "A2", "A2", SHeader)
	xlsx.SetCellStyle(sheet1Name, "A3", "A3", SHeader)
	xlsx.SetCellStyle(sheet1Name, "A5", "A5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "B5", "B5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "C5", "C5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "D5", "D5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "E5", "E5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "F5", "F5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "G5", "G5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "H5", "H5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "I5", "I5", SHeader3)

	xlsx.SetRowHeight(sheet1Name, 1, 20)
	xlsx.SetRowHeight(sheet1Name, 2, 20)
	xlsx.SetRowHeight(sheet1Name, 3, 20)
	xlsx.SetRowHeight(sheet1Name, 4, 20)

	// Loop Report
	var tAkomodasi float64
	var tDeklarasi float64
	for i, v := range report {
		// header
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+colStart), i+1)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+colStart), model.ParseString(v.Nomor))

		if v.TglBerangkat != nil {
			date, errr := time.Parse(time.RFC3339, *v.TglBerangkat)
			if errr != nil {
				fmt.Println(errr)
			}
			d := date.Format("02-01-2006")
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+colStart), d)
		}

		if v.TglKembali != nil {
			date, errr := time.Parse(time.RFC3339, *v.TglKembali)
			if errr != nil {
				fmt.Println(errr)
			}
			d := date.Format("02-01-2006")
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+colStart), d)
		}

		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", i+colStart), model.ParseString(v.Nama))
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+colStart), model.ParseString(v.NamaBidang))
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+colStart), model.ParseString(v.NamaJabatan))

		totalAkomodasi := model.ParseFloat(v.Akomodasi)
		totalDeklarasi := model.ParseFloat(v.BiayaDinas) + model.ParseFloat(v.Reimbursement)
		tAkomodasi += totalAkomodasi
		tDeklarasi += totalDeklarasi
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("H%d", i+colStart), totalAkomodasi)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("I%d", i+colStart), totalDeklarasi)

		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("A%d", i+colStart), fmt.Sprintf("A%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("B%d", i+colStart), fmt.Sprintf("B%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("C%d", i+colStart), fmt.Sprintf("C%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("D%d", i+colStart), fmt.Sprintf("D%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("E%d", i+colStart), fmt.Sprintf("E%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("F%d", i+colStart), fmt.Sprintf("F%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("G%d", i+colStart), fmt.Sprintf("G%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("H%d", i+colStart), fmt.Sprintf("H%d", i+colStart), SChildNum)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("I%d", i+colStart), fmt.Sprintf("I%d", i+colStart), SChildNum)
	}

	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", len(report)+colStart), "TOTAL")
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("H%d", len(report)+colStart), tAkomodasi)
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("I%d", len(report)+colStart), tDeklarasi)

	xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("G%d", len(report)+colStart), fmt.Sprintf("G%d", len(report)+colStart), SChild)
	xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("H%d", len(report)+colStart), fmt.Sprintf("H%d", len(report)+colStart), SChildNum)
	xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("I%d", len(report)+colStart), fmt.Sprintf("I%d", len(report)+colStart), SChildNum)

	// Set Width
	xlsx.SetColWidth(sheet1Name, "A", "A", 5)
	xlsx.SetColWidth(sheet1Name, "B", "B", 20)
	xlsx.SetColWidth(sheet1Name, "C", "C", 15)
	xlsx.SetColWidth(sheet1Name, "D", "D", 15)
	xlsx.SetColWidth(sheet1Name, "E", "E", 20)
	xlsx.SetColWidth(sheet1Name, "F", "F", 20)
	xlsx.SetColWidth(sheet1Name, "G", "G", 20)
	xlsx.SetColWidth(sheet1Name, "H", "H", 15)
	xlsx.SetColWidth(sheet1Name, "I", "I", 15)

	xlsx.SetActiveSheet(index)

	return
}

func (s *ReportServiceImpl) ExportRekapBpdBagian(req FilterReport) (xlsx *excelize.File, err error) {
	report, err := s.ReportRepository.RptRekapBpd(req)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Report not found")
	}

	namaBidang := "SEMUA BAGIAN"
	if req.IdBidang != "" {
		bidangID, _ := uuid.FromString(req.IdBidang)
		bidang, err := s.BidangRepository.ResolveByID(bidangID)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Bidang not found")
		}

		namaBidang = "BAGIAN " + bidang.Nama
	}

	d1, err := time.Parse(model.DefaultDateFormat, req.TglAwal)
	if err != nil {
		fmt.Println(err)
	}
	date1 := d1.Format("02/01/2006")

	d2, err := time.Parse(model.DefaultDateFormat, req.TglAkhir)
	if err != nil {
		fmt.Println(err)
	}
	date2 := d2.Format("02/01/2006")
	ketTgl := "TANGGAL " + date1 + " s.d " + date2

	xlsx = excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet1Name := "Rekap Perbagian"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)
	index, err := xlsx.NewSheet(sheet1Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Style
	exp := "#,##0"
	SHeader, _ := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
	})

	SHeader3, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	SChild, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	SChildNum, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &exp,
	})

	colStart := 6
	xlsx.SetCellValue(sheet1Name, "A1", "LAPORAN REKAP BIAYA PERJALANAN DINAS")
	xlsx.SetCellValue(sheet1Name, "A2", namaBidang)
	xlsx.SetCellValue(sheet1Name, "A3", ketTgl)

	xlsx.SetCellValue(sheet1Name, "A5", "No")
	xlsx.SetCellValue(sheet1Name, "B5", "Nama")
	xlsx.SetCellValue(sheet1Name, "C5", "Sub Bagian")
	xlsx.SetCellValue(sheet1Name, "D5", "Nomor BPD")
	xlsx.SetCellValue(sheet1Name, "E5", "Tgl Berangkat")
	xlsx.SetCellValue(sheet1Name, "F5", "Tgl Kembali")
	xlsx.SetCellValue(sheet1Name, "G5", "Tujuan")
	xlsx.SetCellValue(sheet1Name, "H5", "Total Hari")
	xlsx.SetCellValue(sheet1Name, "I5", "Status Approval")
	xlsx.SetCellValue(sheet1Name, "J5", "Penyelesaian")
	xlsx.SetCellValue(sheet1Name, "K5", "Total Biaya")

	// Style
	xlsx.SetCellStyle(sheet1Name, "A1", "A1", SHeader)
	xlsx.SetCellStyle(sheet1Name, "A2", "A2", SHeader)
	xlsx.SetCellStyle(sheet1Name, "A3", "A3", SHeader)
	xlsx.SetCellStyle(sheet1Name, "A5", "A5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "B5", "B5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "C5", "C5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "D5", "D5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "E5", "E5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "F5", "F5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "G5", "G5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "H5", "H5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "I5", "I5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "J5", "J5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "K5", "K5", SHeader3)

	xlsx.SetRowHeight(sheet1Name, 1, 20)
	xlsx.SetRowHeight(sheet1Name, 2, 20)
	xlsx.SetRowHeight(sheet1Name, 3, 20)
	xlsx.SetRowHeight(sheet1Name, 4, 20)

	// Loop Report
	var totalBiaya float64 = 0
	for i, v := range report {
		// header
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+colStart), i+1)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+colStart), model.ParseString(v.Nama))
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+colStart), model.ParseString(v.NamaUnor))
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+colStart), model.ParseString(v.Nomor))

		if v.TglBerangkat != nil {
			date, errr := time.Parse(time.RFC3339, *v.TglBerangkat)
			if errr != nil {
				fmt.Println(errr)
			}
			d := date.Format("02-01-2006")
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", i+colStart), d)
		}

		if v.TglKembali != nil {
			date, errr := time.Parse(time.RFC3339, *v.TglKembali)
			if errr != nil {
				fmt.Println(errr)
			}
			d := date.Format("02-01-2006")
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+colStart), d)
		}

		// xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+colStart), model.ParseString(v.Tujuan))
		// xlsx.SetCellValue(sheet1Name, fmt.Sprintf("H%d", i+colStart), model.ParseInt(v.JmlHari))
		// xlsx.SetCellValue(sheet1Name, fmt.Sprintf("I%d", i+colStart), model.ParseString(v.KetStatus))
		// xlsx.SetCellValue(sheet1Name, fmt.Sprintf("J%d", i+colStart), model.ParseString(v.Penyelesaian))
		// if v.Total != nil {
		// 	totalBiaya += *v.Total
		// 	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("K%d", i+colStart), *v.Total)
		// }

		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("A%d", i+colStart), fmt.Sprintf("A%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("B%d", i+colStart), fmt.Sprintf("B%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("C%d", i+colStart), fmt.Sprintf("C%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("D%d", i+colStart), fmt.Sprintf("D%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("E%d", i+colStart), fmt.Sprintf("E%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("F%d", i+colStart), fmt.Sprintf("F%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("G%d", i+colStart), fmt.Sprintf("G%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("H%d", i+colStart), fmt.Sprintf("H%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("I%d", i+colStart), fmt.Sprintf("I%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("J%d", i+colStart), fmt.Sprintf("J%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("K%d", i+colStart), fmt.Sprintf("K%d", i+colStart), SChildNum)
	}

	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("J%d", len(report)+colStart), "TOTAL")
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("K%d", len(report)+colStart), totalBiaya)

	xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("J%d", len(report)+colStart), fmt.Sprintf("J%d", len(report)+colStart), SChild)
	xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("K%d", len(report)+colStart), fmt.Sprintf("K%d", len(report)+colStart), SChildNum)

	// Set Width
	xlsx.SetColWidth(sheet1Name, "A", "A", 5)
	xlsx.SetColWidth(sheet1Name, "B", "B", 20)
	xlsx.SetColWidth(sheet1Name, "C", "C", 20)
	xlsx.SetColWidth(sheet1Name, "D", "E", 15)
	xlsx.SetColWidth(sheet1Name, "F", "F", 15)
	xlsx.SetColWidth(sheet1Name, "G", "G", 20)
	xlsx.SetColWidth(sheet1Name, "H", "H", 15)
	xlsx.SetColWidth(sheet1Name, "I", "I", 15)
	xlsx.SetColWidth(sheet1Name, "J", "J", 15)
	xlsx.SetColWidth(sheet1Name, "K", "K", 15)

	xlsx.SetActiveSheet(index)

	return
}

func (s *ReportServiceImpl) ExportBiayaAkomodasiDetail(req FilterReport) (xlsx *excelize.File, err error) {
	report, err := s.ReportRepository.RptRekapAkomodasi(req)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Report not found")
	}

	namaBidang := "SEMUA BAGIAN"
	if req.IdBidang != "" {
		bidangID, _ := uuid.FromString(req.IdBidang)
		bidang, err := s.BidangRepository.ResolveByID(bidangID)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Bidang not found")
		}

		namaBidang = "BAGIAN " + bidang.Nama
	}

	d1, err := time.Parse(model.DefaultDateFormat, req.TglAwal)
	if err != nil {
		fmt.Println(err)
	}
	date1 := d1.Format("02/01/2006")

	d2, err := time.Parse(model.DefaultDateFormat, req.TglAkhir)
	if err != nil {
		fmt.Println(err)
	}
	date2 := d2.Format("02/01/2006")
	ketTgl := "TANGGAL " + date1 + " s.d " + date2

	xlsx = excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet1Name := "Rekap Total"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)
	index, err := xlsx.NewSheet(sheet1Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Style
	exp := "#,##0"
	SHeader, _ := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
	})

	SHeader3, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	SChild, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	SChildNum, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &exp,
	})

	colStart := 6
	xlsx.SetCellValue(sheet1Name, "A1", "LAPORAN BIAYA AKOMODASI DETAIL")
	xlsx.SetCellValue(sheet1Name, "A2", namaBidang)
	xlsx.SetCellValue(sheet1Name, "A3", ketTgl)

	xlsx.SetCellValue(sheet1Name, "A5", "No")
	xlsx.SetCellValue(sheet1Name, "B5", "Nomor BPD")
	xlsx.SetCellValue(sheet1Name, "C5", "Nama")
	xlsx.SetCellValue(sheet1Name, "D5", "Jabatan")
	xlsx.SetCellValue(sheet1Name, "E5", "Tanggal Berangkat")
	xlsx.SetCellValue(sheet1Name, "F5", "Tanggal Kembali")
	xlsx.SetCellValue(sheet1Name, "G5", "Jenis Akomodasi")
	xlsx.SetCellValue(sheet1Name, "H5", "Kode Booking")
	xlsx.SetCellValue(sheet1Name, "I5", "Harga")
	xlsx.SetCellValue(sheet1Name, "J5", "Total")

	// Style
	xlsx.SetCellStyle(sheet1Name, "A1", "A1", SHeader)
	xlsx.SetCellStyle(sheet1Name, "A2", "A2", SHeader)
	xlsx.SetCellStyle(sheet1Name, "A3", "A3", SHeader)
	xlsx.SetCellStyle(sheet1Name, "A5", "A5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "B5", "B5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "C5", "C5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "D5", "D5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "E5", "E5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "F5", "F5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "G5", "G5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "H5", "H5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "I5", "I5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "J5", "J5", SHeader3)

	xlsx.SetRowHeight(sheet1Name, 1, 20)
	xlsx.SetRowHeight(sheet1Name, 2, 20)
	xlsx.SetRowHeight(sheet1Name, 3, 20)
	xlsx.SetRowHeight(sheet1Name, 4, 20)

	// Loop Report
	iCol := 0
	var totalBiaya float64 = 0
	no := 0
	for _, h := range report {
		var totalSub float64 = 0
		if h.Details != nil {
			var akmdet []RekapBiayaAkomodasiDet
			err := json.Unmarshal(*h.Details, &akmdet)
			if err != nil {
				fmt.Println(err)
			}

			for _, d := range akmdet {
				totalSub += model.ParseFloat(d.Nominal)
			}

			for i, v := range akmdet {
				// header
				nomor := ""
				nama := ""
				jabatan := ""
				tgl1 := ""
				tgl2 := ""
				if i == 0 {
					no++
					nomor = model.ParseString(h.Nomor)
					nama = model.ParseString(h.Nama)
					jabatan = model.ParseString(h.NamaJabatan)
					tgl1 = model.ParseString(h.TglBerangkat)
					tgl2 = model.ParseString(h.TglKembali)
					xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", iCol+colStart), no)
				}
				xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", iCol+colStart), nomor)
				xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", iCol+colStart), nama)
				xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", iCol+colStart), jabatan)

				if tgl1 != "" {
					date, errr := time.Parse(time.RFC3339, tgl1)
					if errr != nil {
						fmt.Println(errr)
					}
					d := date.Format("02-01-2006")
					xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", iCol+colStart), d)
				}

				if tgl2 != "" {
					date, errr := time.Parse(time.RFC3339, tgl2)
					if errr != nil {
						fmt.Println(errr)
					}
					d := date.Format("02-01-2006")
					xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", iCol+colStart), d)
				}

				xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", iCol+colStart), model.ParseString(v.NamaBiaya))
				xlsx.SetCellValue(sheet1Name, fmt.Sprintf("H%d", iCol+colStart), model.ParseString(v.Keterangan))
				xlsx.SetCellValue(sheet1Name, fmt.Sprintf("I%d", iCol+colStart), model.ParseFloat(v.Nominal))

				if i == 0 {
					xlsx.SetCellValue(sheet1Name, fmt.Sprintf("J%d", iCol+colStart), totalSub)
				}

				xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("A%d", iCol+colStart), fmt.Sprintf("A%d", iCol+colStart), SChild)
				xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("B%d", iCol+colStart), fmt.Sprintf("B%d", iCol+colStart), SChild)
				xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("C%d", iCol+colStart), fmt.Sprintf("C%d", iCol+colStart), SChild)
				xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("D%d", iCol+colStart), fmt.Sprintf("D%d", iCol+colStart), SChild)
				xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("E%d", iCol+colStart), fmt.Sprintf("E%d", iCol+colStart), SChild)
				xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("F%d", iCol+colStart), fmt.Sprintf("F%d", iCol+colStart), SChild)
				xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("G%d", iCol+colStart), fmt.Sprintf("G%d", iCol+colStart), SChild)
				xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("H%d", iCol+colStart), fmt.Sprintf("H%d", iCol+colStart), SChild)
				xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("I%d", iCol+colStart), fmt.Sprintf("I%d", iCol+colStart), SChildNum)
				xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("J%d", iCol+colStart), fmt.Sprintf("J%d", iCol+colStart), SChildNum)

				iCol++
			}
		}

		// sub total
		totalBiaya += totalSub
	}

	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("I%d", iCol+colStart), "TOTAL")
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("J%d", iCol+colStart), totalBiaya)

	xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("I%d", iCol+colStart), fmt.Sprintf("I%d", iCol+colStart), SChild)
	xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("J%d", iCol+colStart), fmt.Sprintf("J%d", iCol+colStart), SChildNum)

	// Set Width
	xlsx.SetColWidth(sheet1Name, "A", "A", 5)
	xlsx.SetColWidth(sheet1Name, "B", "B", 20)
	xlsx.SetColWidth(sheet1Name, "C", "C", 20)
	xlsx.SetColWidth(sheet1Name, "D", "D", 15)
	xlsx.SetColWidth(sheet1Name, "E", "E", 15)
	xlsx.SetColWidth(sheet1Name, "F", "F", 15)
	xlsx.SetColWidth(sheet1Name, "G", "G", 15)
	xlsx.SetColWidth(sheet1Name, "H", "H", 15)
	xlsx.SetColWidth(sheet1Name, "I", "I", 15)
	xlsx.SetColWidth(sheet1Name, "J", "J", 15)

	xlsx.SetActiveSheet(index)

	return
}

func (s *ReportServiceImpl) RptRekapReimbusment(req FilterReport) (data []ReportRekapAkReim, err error) {
	return s.ReportRepository.RptRekapReimbusment(req)
}

func (s *ReportServiceImpl) RptRekapAkomodasi(req FilterReport) (data []RekapBiayaAkomodasi, err error) {
	return s.ReportRepository.RptRekapAkomodasi(req)
}
func (s *ReportServiceImpl) ExportRekapReimAkm(req FilterReport) (xlsx *excelize.File, err error) {
	report, err := s.ReportRepository.RptRekapBpd(req)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Report not found")
	}

	namaBidang := "SEMUA BAGIAN"
	if req.IdBidang != "" {
		bidangID, _ := uuid.FromString(req.IdBidang)
		bidang, err := s.BidangRepository.ResolveByID(bidangID)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Bidang not found")
		}

		namaBidang = "BAGIAN " + bidang.Nama
	}

	d1, err := time.Parse(model.DefaultDateFormat, req.TglAwal)
	if err != nil {
		fmt.Println(err)
	}
	date1 := d1.Format("02/01/2006")

	d2, err := time.Parse(model.DefaultDateFormat, req.TglAkhir)
	if err != nil {
		fmt.Println(err)
	}
	date2 := d2.Format("02/01/2006")
	ketTgl := "TANGGAL " + date1 + " s.d " + date2

	xlsx = excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet1Name := ""
	if req.Type == "AKOMODASI" {
		sheet1Name = "Rekap Biaya Akomodasi"
	} else if req.Type == "REIMBURSEMENT" {
		sheet1Name = "Rekap Biaya Reimbursement"
	}

	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)
	index, err := xlsx.NewSheet(sheet1Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Style
	exp := "#,##0"
	SHeader, _ := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
	})

	SHeader3, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	SChild, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	SChildNum, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &exp,
	})

	colStart := 6
	if req.Type == "AKOMODASI" {
		xlsx.SetCellValue(sheet1Name, "A1", "LAPORAN REKAP BIAYA AKOMODASI")
	} else if req.Type == "REIMBURSEMENT" {
		xlsx.SetCellValue(sheet1Name, "A1", "LAPORAN REKAP BIAYA REIMBURSEMENT")
	}
	xlsx.SetCellValue(sheet1Name, "A2", namaBidang)
	xlsx.SetCellValue(sheet1Name, "A3", ketTgl)

	xlsx.SetCellValue(sheet1Name, "A5", "No")
	xlsx.SetCellValue(sheet1Name, "B5", "Nomor BPD")
	xlsx.SetCellValue(sheet1Name, "C5", "Nama")
	xlsx.SetCellValue(sheet1Name, "D5", "Bidang")
	xlsx.SetCellValue(sheet1Name, "E5", "Tgl Berangkat")
	xlsx.SetCellValue(sheet1Name, "F5", "Tgl Kembali")
	xlsx.SetCellValue(sheet1Name, "G5", "Total")

	// Style

	xlsx.SetCellStyle(sheet1Name, "A1", "A1", SHeader)
	xlsx.SetCellStyle(sheet1Name, "A2", "A2", SHeader)
	xlsx.SetCellStyle(sheet1Name, "A3", "A3", SHeader)
	xlsx.SetCellStyle(sheet1Name, "A5", "A5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "B5", "B5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "C5", "C5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "D5", "D5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "E5", "E5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "F5", "F5", SHeader3)
	xlsx.SetCellStyle(sheet1Name, "G5", "G5", SHeader3)

	xlsx.SetRowHeight(sheet1Name, 1, 20)
	xlsx.SetRowHeight(sheet1Name, 2, 20)
	xlsx.SetRowHeight(sheet1Name, 3, 20)
	xlsx.SetRowHeight(sheet1Name, 4, 20)

	// Loop Report
	var totalBiaya float64 = 0
	for i, v := range report {
		// header
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+colStart), i+1)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+colStart), model.ParseString(v.Nomor))
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+colStart), model.ParseString(v.Nama))
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+colStart), model.ParseString(v.NamaBidang))

		if v.TglBerangkat != nil {
			date, errr := time.Parse(time.RFC3339, *v.TglBerangkat)
			if errr != nil {
				fmt.Println(errr)
			}
			d := date.Format("02-01-2006")
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", i+colStart), d)
		}

		if v.TglKembali != nil {
			date, errr := time.Parse(time.RFC3339, *v.TglKembali)
			if errr != nil {
				fmt.Println(errr)
			}
			d := date.Format("02-01-2006")
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+colStart), d)
		}

		if req.Type == "AKOMODASI" {
			totalBiaya += model.ParseFloat(v.Akomodasi)
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+colStart), model.ParseFloat(v.Akomodasi))
		} else if req.Type == "REIMBURSEMENT" {
			totalBiaya += model.ParseFloat(v.Reimbursement)
			xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+colStart), model.ParseFloat(v.Reimbursement))
		}

		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("A%d", i+colStart), fmt.Sprintf("A%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("B%d", i+colStart), fmt.Sprintf("B%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("C%d", i+colStart), fmt.Sprintf("C%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("D%d", i+colStart), fmt.Sprintf("D%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("E%d", i+colStart), fmt.Sprintf("E%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("F%d", i+colStart), fmt.Sprintf("F%d", i+colStart), SChild)
		xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("G%d", i+colStart), fmt.Sprintf("G%d", i+colStart), SChildNum)
	}
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", len(report)+colStart), "TOTAL")
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", len(report)+colStart), totalBiaya)

	xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("F%d", len(report)+colStart), fmt.Sprintf("F%d", len(report)+colStart), SChild)
	xlsx.SetCellStyle(sheet1Name, fmt.Sprintf("G%d", len(report)+colStart), fmt.Sprintf("G%d", len(report)+colStart), SChildNum)

	// Set Width
	xlsx.SetColWidth(sheet1Name, "A", "A", 5)
	xlsx.SetColWidth(sheet1Name, "B", "B", 20)
	xlsx.SetColWidth(sheet1Name, "C", "C", 20)
	xlsx.SetColWidth(sheet1Name, "D", "E", 15)
	xlsx.SetColWidth(sheet1Name, "F", "F", 15)
	xlsx.SetColWidth(sheet1Name, "G", "G", 20)

	xlsx.SetActiveSheet(index)

	return
}
