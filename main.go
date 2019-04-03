package main

import (
	"flag"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"os"

	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var TEMP_PNG = "temp.png"

func main() {

	var (
		bpath = flag.String("p", "", "path to image(base image)")
		lpath = flag.String("l", "", "path to image(layer image)")
		text  = flag.String("t", "sample", "text to insert")
	)
	flag.Parse()

	baseFile, err := os.Open(*bpath)
	if err != nil {
		log.Fatal(err)
	}
	defer baseFile.Close()
	baseImage, formatName, err := image.Decode(baseFile)
	if err != nil {
		panic(err)
		// log.Fatal(err)
	}
	fmt.Println(formatName)

	if len(*lpath) != 0 {
		baseConfig := getImageConfig(*bpath)
		layerFileConfig := getImageConfig(*lpath)
		// height, width := calcCenterPoint(baseConfig, layerFileConfig)
		height, width := calcLowerLeftPoint(baseConfig, layerFileConfig)

		layerFile, err := os.Open(*lpath)
		if err != nil {
			log.Fatal(err)
		}
		defer layerFile.Close()
		layerImage, _, err := image.Decode(layerFile)
		if err != nil {
			panic(err)
			// log.Fatal(err)
		}

		startPointLogo := image.Point{width, height}
		logoRectangle := image.Rectangle{startPointLogo, startPointLogo.Add(layerImage.Bounds().Size())}
		originRectangle := image.Rectangle{image.Point{0, 0}, baseImage.Bounds().Size()}
		rgba := image.NewRGBA(originRectangle)
		draw.Draw(rgba, originRectangle, baseImage, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, logoRectangle, layerImage, image.Point{0, 0}, draw.Over)
		out, err := os.Create("res.jpeg")
		if err != nil {
			fmt.Println(err)
		}

		var opt jpeg.Options
		opt.Quality = 100

		jpeg.Encode(out, rgba, &opt)
	} else if len(*text) != 0 {
		writeTextImg(*text)
		baseConfig := getImageConfig(*bpath)
		layerFileConfig := getImageConfig(TEMP_PNG)
		// height, width := calcCenterPoint(baseConfig, layerFileConfig)
		height, width := calcCenterPoint(baseConfig, layerFileConfig)
		layerFile, err := os.Open(TEMP_PNG)
		if err != nil {
			log.Fatal(err)
		}
		defer layerFile.Close()
		layerImage, _, err := image.Decode(layerFile)
		if err != nil {
			panic(err)
			// log.Fatal(err)
		}
		startPointLogo := image.Point{height, width}
		logoRectangle := image.Rectangle{startPointLogo, startPointLogo.Add(layerImage.Bounds().Size())}
		originRectangle := image.Rectangle{image.Point{0, 0}, baseImage.Bounds().Size()}
		rgba := image.NewRGBA(originRectangle)
		draw.Draw(rgba, originRectangle, baseImage, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, logoRectangle, layerImage, image.Point{0, 0}, draw.Over)
		out, err := os.Create("res.jpeg")
		if err != nil {
			fmt.Println(err)
		}

		var opt jpeg.Options
		opt.Quality = 100

		jpeg.Encode(out, rgba, &opt)

	}

}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func writeTextImg(text string) {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
	// 	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
	// 		xrate := float32(x) / float32(img.Rect.Max.X)
	// 		yrate := float32(y) / float32(img.Rect.Max.Y)
	// 		img.Set(x, y, color.RGBA{uint8(xrate * 255 * yrate), 0, 0, uint8(yrate * 255)})
	// 	}
	// }

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{fixed.Int26_6(24 * 64), fixed.Int26_6(18 * 64)},
	}
	d.DrawString(text)

	file, err := os.Create(TEMP_PNG)
	if err != nil {
		panic(err.Error())
	}

	if err := png.Encode(file, img); err != nil {
		panic(err.Error())
	}
	defer file.Close()
}

func getImageConfig(path string) image.Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	baseConfig := getImgConfig(file)
	return baseConfig
}

func getImgConfig(file io.Reader) image.Config {
	config, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return config
}

func calcCenterPoint(baseConfig image.Config, layerConfig image.Config) (int, int) {
	return (baseConfig.Height / 2) - (layerConfig.Height / 2), (baseConfig.Width / 2) - (layerConfig.Width / 2)
}

func calcLowerLeftPoint(baseConfig image.Config, layerConfig image.Config) (int, int) {
	return (baseConfig.Height - layerConfig.Height), 0
}

// func calcLowerRightPoint(baseConfig image.Config, layerConfig image.Config) (int, int) {
// 	return (baseConfig.Height - layerConfig.Height), (baseConfig.Width - layerConfig.Width)
// }

func fontload(fname string) []byte {
	file, err := os.Open(fname)
	defer file.Close()
	if err != nil {
		fmt.Println("error:file\n", err)
		return nil
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("error:fileread\n", err)
		return nil
	}

	return bytes
}
