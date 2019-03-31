package main

import (
	"flag"
	"fmt"
	"image"
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
	fmt.Println(*bpath)

	baseF, err := os.Open(*bpath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer baseF.Close()
	fmt.Println(baseF)
	conf, format, err := image.DecodeConfig(baseF)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(format)
	fmt.Println(conf)
	fmt.Println(baseF)
	fmt.Println("---")
	baseImage, formatName, err := image.Decode(baseF)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(formatName)

	layerF, err := os.Open(*lpath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(layerF)
	layerImage, _, err := image.Decode(layerF)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	startPointLogo := image.Point{0, 0} //上に
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
	//----
	// reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	layerF2, err := os.Open(*lpath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer layerF2.Close()
	fmt.Println(layerF)
	fmt.Println("---")
	fmt.Println(layerF2)
	if layerF == layerF2 {
		fmt.Println("true")
	}
	config, format, err := image.DecodeConfig(layerF2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Width:", config.Width, "Height:", config.Height, "Format:", format)
}
