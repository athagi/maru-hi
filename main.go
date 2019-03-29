package main

import (
	"flag"
	"fmt"
	"image"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	var (
		path = flag.String("p", "", "path to image")
	)
	flag.Parse()
	fmt.Println(*path)

	f, err := os.Open(*path)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer f.Close()
	conf, format, err := image.DecodeConfig(f)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(format)
	fmt.Println(conf)
}
