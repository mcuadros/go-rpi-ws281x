package ws281x

/*
#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lws2811
#include <stdint.h>
#include <string.h>
#include <ws2811.h>
void ws2811_set_led(ws2811_t *ws2811, int index, uint32_t value) {
	ws2811->channel[0].leds[index] = value;
}

uint32_t ws2811_get_led(ws2811_t *ws2811, int index) {
    return ws2811->channel[0].leds[index];
}
*/
import "C"

import (
	"fmt"
	"image"
	"image/color"
	"unsafe"
)

type StripType int

const (
	// 4 color R, G, B and W ordering
	StripRGBW StripType = 0x18100800
	StripRBGW StripType = 0x18100008
	StripGRBW StripType = 0x18081000
	StripGBRW StripType = 0x18080010
	StripBRGW StripType = 0x18001008
	StripBGRW StripType = 0x18000810

	// 3 color R, G and B ordering
	StripRGB StripType = 0x00100800
	StripRBG StripType = 0x00100008
	StripGRB StripType = 0x00081000
	StripGBR StripType = 0x00080010
	StripBRG StripType = 0x00001008
	StripBGR StripType = 0x00000810
)

var DefaultConfig = HardwareConfig{
	Pin:        18,
	Frequency:  800000,
	DMA:        5,
	Brightness: 30,
	StripType:  StripGRB,
}

type HardwareConfig struct {
	Pin        int
	Frequency  int
	DMA        int
	Invert     bool
	Channel    int
	Brightness int
	StripType  StripType
}

type Canvas struct {
	Width, Height int
	Config        *HardwareConfig

	leds   *C.ws2811_t
	closed bool
}

func NewCanvas(w, h int, config *HardwareConfig) (*Canvas, error) {
	c := &Canvas{
		Width:  w,
		Height: h,
		Config: config,

		leds: (*C.ws2811_t)(C.malloc(C.sizeof_ws2811_t)),
	}

	if c.leds == nil {
		return nil, fmt.Errorf("unable to allocate memory")
	}

	C.memset(unsafe.Pointer(c.leds), 0, C.sizeof_ws2811_t)

	c.initializeChannels()
	c.initializeController()
	return c, nil
}

func (c *Canvas) initializeChannels() {
	for ch := 0; ch < 2; ch++ {
		c.setChannel(0, 0, 0, 0, false, StripRGB)
	}

	c.setChannel(
		c.Config.Channel,
		c.Width*c.Height,
		c.Config.Pin,
		c.Config.Brightness,
		c.Config.Invert,
		c.Config.StripType,
	)

}

func (c *Canvas) setChannel(ch, count, pin, brightness int, inverse bool, t StripType) {
	c.leds.channel[ch].count = C.int(count)
	c.leds.channel[ch].gpionum = C.int(pin)
	c.leds.channel[ch].brightness = C.int(brightness)
	c.leds.channel[ch].invert = C.int(btoi(inverse))
	c.leds.channel[ch].strip_type = C.int(int(t))
}

func (c *Canvas) initializeController() {
	c.leds.freq = C.uint32_t(c.Config.Frequency)
	c.leds.dmanum = C.int(c.Config.DMA)
}

// Begin Initialize library, must be called once before other functions are	called.
func (c *Canvas) Begin() error {
	if resp := int(C.ws2811_init(c.leds)); resp != 0 {
		return fmt.Errorf("ws2811_init failed with code: %d", resp)
	}

	return nil
}

// Update the display with the data from the LED buffer
func (c *Canvas) Show() error {
	if resp := int(C.ws2811_render(c.leds)); resp != 0 {
		return fmt.Errorf("ws2811_render failed with code: %d", resp)
	}

	return nil
}

func (c *Canvas) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds return the topology of the Matrix
func (c *Canvas) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.Width, c.Height)
}

// At return an Color which allows access to the LED display data as
// if it were a sequence of 24-bit RGB values.
func (c *Canvas) At(x, y int) color.Color {
	//	color := C.ws2811_get_led(c.leds, C.int(c.position(x, y)))
	//	return uint32(color)
	return color.Black
}

// Set set LED at position x,y to the provided 24-bit color value.
func (c *Canvas) Set(x, y int, color color.Color) {
	C.ws2811_set_led(c.leds, C.int(c.position(x, y)), C.uint32_t(colorToUint32(color)))
}

func (c *Canvas) position(x, y int) int {
	return x + (y * c.Width)
}

// Size return the number of pixels in the display
func (c *Canvas) Size() int {
	return int(c.leds.channel[c.Config.Channel].count)
}

func (c *Canvas) Clear() error {
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			c.Set(x, y, color.Black)
		}
	}

	return c.Show()
}

func (c *Canvas) Close() {
	if c.closed {
		return
	}

	c.closed = true
	C.ws2811_fini(c.leds)
}

func colorToUint32(c color.Color) uint32 {
	// A color's RGBA method returns values in the range [0, 65535]
	red, green, blue, alpha := c.RGBA()

	return (alpha>>8)<<24 | (red>>8)<<16 | (green>>8)<<8 | blue>>8
}

func btoi(b bool) int {
	if b {
		return 1
	}

	return 0
}
