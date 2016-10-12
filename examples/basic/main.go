package main

import (
	"flag"
	"image/color"
	"time"

	"github.com/mcuadros/go-rpi-ws281x"
)

var gpioPin = flag.Int("gpio-pin", 18, "GPIO pin")
var width = flag.Int("width", 8, "LED matrix width")
var height = flag.Int("height", 4, "LED matrix height")
var brightness = flag.Int("brightness", 64, "Brightness (0-255)")

const (
	pixelColor uint32 = 255 << 16 // green
)

func main() {
	config := ws281x.DefaultConfig
	config.Brightness = *brightness

	c, err := ws281x.NewCanvas(*width, *height, &config)
	if err != nil {
		fatal(err)
	}

	defer c.Close()
	if err := c.Begin(); err != nil {
		fatal(err)
	}

	color := color.RGBA{255, 0, 0, 255}
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			c.Set(x, y, color)
			c.Show()
			time.Sleep(10 * time.Millisecond)
		}
	}

	c.Clear()
}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
