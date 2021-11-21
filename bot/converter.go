package bot

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/qeesung/image2ascii/convert"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	imglib "image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const BASE_TEXT_WIDTH = 152
const BASE_TEXT_HEIGHT = 55

func loadFont() (*truetype.Font, error) {
	fontFile := "static/UbuntuMono-R.ttf"
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	parsedFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	return parsedFont, nil
}

func DownloadImage(url string) (imglib.Image, string, int, int, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, "", 0, 0, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	img, fileName, err := imglib.Decode(response.Body)
	bounds := img.Bounds()

	return img, fileName, bounds.Max.X, bounds.Max.Y, err
}

func ConvertImageToAscii(image imglib.Image) string {
	bounds := image.Bounds()
	ratio := float64(bounds.Max.X) / float64(bounds.Max.Y)


	convertOptions := convert.DefaultOptions
	convertOptions.Colored = false

	convertOptions.FixedWidth = BASE_TEXT_WIDTH
	convertOptions.FixedHeight = BASE_TEXT_HEIGHT

	if ratio < 1 {
		convertOptions.FixedHeight = int(float64(BASE_TEXT_HEIGHT) * (float64(bounds.Max.Y) / float64(bounds.Max.X)))
	} else {
		convertOptions.FixedWidth = int(float64(BASE_TEXT_WIDTH) * ratio)
	}

	log.Printf("image width: %d", bounds.Max.X)
	log.Printf("image height: %d", bounds.Max.Y)
	log.Printf("text width: %d", convertOptions.FixedWidth)
	log.Printf("text height: %d", convertOptions.FixedHeight)
	log.Printf("ratio: %e", ratio)

	converter := convert.NewImageConverter()
	asciiArt := converter.Image2ASCIIString(image, &convertOptions)
	return asciiArt
}

func GenerateImageFromText(
	textContent string,
	fgColorHex string,
	bgColorHex string,
	fontSize float64,
	width int,
	height int,
) ([]byte, error) {
	fgColor := color.RGBA{0xff, 0xff, 0xff, 0xff}
	if len(fgColorHex) == 7 {
		_, err := fmt.Sscanf(fgColorHex, "#%02x%02x%02x", &fgColor.R, &fgColor.G, &fgColor.B)
		if err != nil {
			log.Println(err)
			fgColor = color.RGBA{0x2e, 0x34, 0x36, 0xff}
		}
	}

	bgColor := color.RGBA{0x30, 0x0a, 0x24, 0xff}
	if len(bgColorHex) == 7 {
		_, err := fmt.Sscanf(bgColorHex, "#%02x%02x%02x", &bgColor.R, &bgColor.G, &bgColor.B)
		if err != nil {
			log.Println(err)
			bgColor = color.RGBA{0x30, 0x0a, 0x24, 0xff}
		}
	}

	loadedFont, err := loadFont()
	if err != nil {
		return nil, err
	}

	code := strings.Replace(textContent, "\t", "    ", -1) // convert tabs into spaces
	text := strings.Split(code, "\n") // split newlines into arrays

	fg := imglib.NewUniform(fgColor)
	bg := imglib.NewUniform(bgColor)
	rgba := imglib.NewRGBA(imglib.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), bg, imglib.Pt(0, 0), draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(loadedFont)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	textXOffset := 50
	textYOffset := 10 + int(c.PointToFixed(fontSize)>>6)

	pt := freetype.Pt(textXOffset, textYOffset)
	for _, s := range text {
		_, err = c.DrawString(strings.Replace(s, "\r", "", -1), pt)
		if err != nil {
			return nil, err
		}
		pt.Y += c.PointToFixed(fontSize * 1.5)
	}

	b := new(bytes.Buffer)
	if err := png.Encode(b, rgba); err != nil {
		log.Println("unable to encode image.")
		return nil, err
	}
	return b.Bytes(), nil
}
