package bot

import (
	"fmt"
	"github.com/qeesung/image2ascii/convert"
	"image"
	"net/http"
)

func DownloadImage(url string) (image.Image, string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}

	defer response.Body.Close()

	img, fileName, err := image.Decode(response.Body)

	return img, fileName, err
}

func ConvertImageToAscii(image image.Image) {
	convertOptions := convert.DefaultOptions
	convertOptions.Colored = false
	convertOptions.FixedWidth = 200
	convertOptions.FixedHeight = 80

	converter := convert.NewImageConverter()
	fmt.Print(converter.Image2ASCIIString(image, &convertOptions))
}
