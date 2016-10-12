package main

import (
	"flag"
	"image"
	"os"
	"time"

	"github.com/mcuadros/go-rpi-ws281x"

	"image/draw"
	_ "image/png"
)

var (
	pin        = flag.Int("gpio-pin", 18, "GPIO pin")
	width      = flag.Int("width", 8, "LED matrix width")
	height     = flag.Int("height", 32, "LED matrix height")
	brightness = flag.Int("brightness", 64, "Brightness (0-255)")
	img        = flag.String("image", "", "image path")
)

func main() {
	f, err := os.Open(*img)
	if err != nil {
		fatal(err)
	}

	m, _, err := image.Decode(f)
	if err != nil {
		fatal(err)
	}

	config := ws281x.DefaultConfig
	config.Brightness = *brightness
	config.Pin = *pin

	c, err := ws281x.NewCanvas(*width, *height, &config)
	if err != nil {
		fatal(err)
	}

	defer c.Close()

	err = c.Initialize()
	fatal(err)

	draw.Draw(c, c.Bounds(), m, image.ZP, draw.Over)

	c.Render()
	time.Sleep(time.Second * 5)
}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
