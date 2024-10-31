package files

import (
	"image"
	"image/color"
	"io"
	"log"
	"os"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/random"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	// TemporaryFileFolder is a temporary folder file
	TemporaryFileFolder = "temp"
)

// CreateEmptyFile is create empty file (temporary)
func CreateEmptyFile() (fileName string, file *os.File, err error) {
	fileName = TemporaryFileFolder + "/" + random.RandStringBytes(5) + ".csv"
	file, err = os.Create(fileName)
	if err != nil {
		return
	}

	err = os.Chmod(fileName, 777)
	if err != nil {
		return
	}

	return
}

// CopyFile is copy file to empty file (temporary)
func CopyFile(dst io.Writer, src io.Reader) (err error) {

	_, err = io.Copy(dst, src)
	if err != nil {
		return
	}

	return
}

type RequestTextOnImg struct {
	BgImgPath string
	FontPath  string
	FontSize  float64
	Text      []string
}

func TextOnImg(request RequestTextOnImg) (image.Image, error) {
	bgImage, err := gg.LoadImage(request.BgImgPath)
	if err != nil {
		return nil, err
	}
	imgWidth := bgImage.Bounds().Dx()
	imgHeight := bgImage.Bounds().Dy()

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: request.FontSize})
	dc := gg.NewContext(imgWidth, imgHeight)
	dc.SetFontFace(face)
	dc.DrawImage(bgImage, 0, 0)

	// if err := dc.LoadFontFace(request.FontPath, request.FontSize); err != nil {
	// 	return nil, err
	// }

	x1 := float64(20)
	x2 := float64(imgWidth/2) - 10.0
	y1 := float64(imgHeight - 100)
	y2 := float64(imgHeight - 70)
	y3 := float64(imgHeight - 40)
	maxWidth := float64(imgWidth) - 60.0
	dc.SetColor(color.White)
	dc.DrawStringWrapped(request.Text[0], x1, y1, 0, 1, maxWidth, 1.5, gg.AlignLeft)
	dc.DrawStringWrapped(request.Text[1], x1, y2, 0, 1, maxWidth, 1.5, gg.AlignLeft)
	dc.DrawStringWrapped(request.Text[2], x2, y3, 0.5, 0.5, maxWidth, 1.3, gg.AlignLeft)
	dc.SavePNG(request.BgImgPath)

	return dc.Image(), nil
}
