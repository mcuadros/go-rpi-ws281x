package main

import (
	"flag"
	"image"
	"time"

	"github.com/mcuadros/go-rpi-ws281x"

	"image/color"
	"image/draw"
	_ "image/png"
)

var (
	pin        = flag.Int("gpio-pin", 18, "GPIO pin")
	width      = flag.Int("width", 8, "LED matrix width")
	height     = flag.Int("height", 32, "LED matrix height")
	brightness = flag.Int("brightness", 64, "Brightness (0-255)")
	freq       = flag.Duration("freq", time.Millisecond*5, "frequency")
)

func main() {
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

	colors := []color.Color{
		color.RGBA{255, 0, 0, 255},
		color.Black,
	}

	var i int
	for {
		draw.Draw(c, c.Bounds(), image.NewUniform(colors[i]), image.ZP, draw.Over)
		c.Render()
		time.Sleep(*freq)

		i++
		if i > 1 {
			i = 0
		}
	}

}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
