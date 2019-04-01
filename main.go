package main

import (
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"os"

	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	var (
		bpath = flag.String("p", "", "path to image(base image)")
		lpath = flag.String("l", "", "path to image(layer image)")
	)
	flag.Parse()

	baseConfig := getImageConfig(*bpath)
	layerFileConfig := getImageConfig(*lpath)
	height, width := calcStartPoint(baseConfig, layerFileConfig)
	fmt.Println(height)
	fmt.Println(width)

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

func calcStartPoint(baseConfig image.Config, layerConfig image.Config) (int, int) {
	return (baseConfig.Height / 2) - (layerConfig.Height / 2), (baseConfig.Width / 2) - (layerConfig.Width / 2)
}
