package ws2811

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
	"unsafe"
)

type StripType int

const (
	StripRGB StripType = 0x100800
	StripRBG StripType = 0x100008
	StripGRB StripType = 0x081000
	StripGBR StripType = 0x080010
	StripBRG StripType = 0x001008
	StripBGR StripType = 0x000810
)

var DefaultConfig = HardwareConfig{
	Pin:        18,
	Frequency:  800000,
	DMA:        5,
	Brightness: 255,
	StripType:  StripRGB,
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

	leds *C.ws2811_t
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

// GetPixelColor return an Color which allows access to the LED display data as
// if it were a sequence of 24-bit RGB values.
func (c *Canvas) GetPixelColor(x, y int) uint32 {
	color := C.ws2811_get_led(c.leds, C.int(c.position(x, y)))
	return uint32(color)
}

// SetPixelColor set LED at position x,y to the provided 24-bit color value.
func (c *Canvas) SetPixelColor(x, y int, color uint32) {
	C.ws2811_set_led(c.leds, C.int(c.position(x, y)), C.uint32_t(color))
}

func (c *Canvas) position(x, y int) int {
	return x + (y * c.Width)
}

// NumPixels return the number of pixels in the display
func (c *Canvas) NumPixels() int {
	return int(c.leds.channel[c.Config.Channel].count)
}

func (c *Canvas) Clear() error {
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			c.SetPixelColor(x, y, 0)
		}
	}

	return c.Show()
}

func (c *Canvas) Close() {
	C.ws2811_fini(c.leds)
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
