package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

var imagePath = flag.String("i", "", "Path to image")
var messages = flag.Bool("m", false, "Is a messages icon")

type Dim struct {
	name      string   // prefix for outout file
	pointSize float64  // point size for image
	scale     []string // point scaling  (1x,2x,3x)
}

var dims = []Dim{
	Dim{name: "20pt", pointSize: 20, scale: []string{"1x", "2x", "3x"}},
	Dim{name: "29pt", pointSize: 29, scale: []string{"1x", "2x", "3x"}},
	Dim{name: "40pt", pointSize: 40, scale: []string{"1x", "2x", "3x"}},
	Dim{name: "60pt", pointSize: 60, scale: []string{"2x", "3x"}},
	Dim{name: "76pt", pointSize: 76, scale: []string{"1x", "2x"}},
	Dim{name: "85.5pt", pointSize: 83.5, scale: []string{"2x"}},
	Dim{name: "1024pt", pointSize: 1024, scale: []string{"1x"}},
}

func main() {
	flag.Parse()
	var _ = resize.Resize
	var path = *imagePath
	var isMessagesIcon = *messages

	if isMessagesIcon {
		dims = []Dim{
			Dim{name: "60x45pt", pointSize: 60, scale: []string{"2x", "3x"}},
			Dim{name: "67x50pt", pointSize: 67, scale: []string{"2x"}},
			Dim{name: "74x55pt", pointSize: 74, scale: []string{"2x"}},
			Dim{name: "27x20pt", pointSize: 27, scale: []string{"2x", "3x"}},
			Dim{name: "32x24pt", pointSize: 32, scale: []string{"2x", "3x"}},
			Dim{name: "1024x768pt", pointSize: 1024, scale: []string{"1x"}},
		}
	}

	if len(path) == 0 {
		return
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	outDir := filepath.Dir(path)
	fName := filepath.Base(path)
	fExt := filepath.Ext(fName)
	fName = strings.TrimSuffix(fName, fExt)
	fName = strings.Replace(fName, " ", "_", -1)

	for _, dim := range dims {
		for _, scale := range dim.scale {
			pointScale := fmt.Sprint(dim.pointSize)
			// if is decimal
			if math.Mod(dim.pointSize, 1.0) != 0 {
				pointScale = strings.Replace(pointScale, ".", "_", -1)
			}
			outPath := fmt.Sprintf("%v/%v_%s", outDir, fName, pointScale)

			var newImg image.Image

			switch scale {
			case "1x":
				outPath = fmt.Sprintf("%v%v", outPath, fExt)
				var dimxy = uint(dim.pointSize)
				newImg = resize.Resize(dimxy, 0, img, resize.Lanczos3)
			case "2x":
				outPath = fmt.Sprintf("%v@2x%v", outPath, fExt)
				var dimxy = uint(dim.pointSize * 2)
				newImg = resize.Resize(dimxy, 0, img, resize.Lanczos3)
			case "3x":
				outPath = fmt.Sprintf("%v@3x%v", outPath, fExt)
				var dimxy = uint(dim.pointSize * 3)
				newImg = resize.Resize(dimxy, 0, img, resize.Lanczos3)
			default:
				fmt.Println("Unrecognized scaling factor", scale)
			}

			fmt.Println("Writing", outPath)
			out, err := os.Create(outPath)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			if newImg != nil {
				png.Encode(out, newImg)
			}
		}
	}

}
