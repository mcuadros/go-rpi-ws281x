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
	Size       int
	Pin        int
	Frequency  int
	DMA        int
	Invert     bool
	Channel    int
	Brightness int
	StripType  StripType
}

type WS281x struct {
	Config *HardwareConfig

	size   int
	leds   *C.ws2811_t
	closed bool
}

func NewWS281x(size int, config *HardwareConfig) (Matrix, error) {
	c := &WS281x{
		Config: config,

		size: size,
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

func (c *WS281x) initializeChannels() {
	for ch := 0; ch < 2; ch++ {
		c.setChannel(0, 0, 0, 0, false, StripRGB)
	}

	c.setChannel(
		c.Config.Channel,
		c.size,
		c.Config.Pin,
		c.Config.Brightness,
		c.Config.Invert,
		c.Config.StripType,
	)

}

func (c *WS281x) setChannel(ch, count, pin, brightness int, inverse bool, t StripType) {
	c.leds.channel[ch].count = C.int(count)
	c.leds.channel[ch].gpionum = C.int(pin)
	c.leds.channel[ch].brightness = C.int(brightness)
	c.leds.channel[ch].invert = C.int(btoi(inverse))
	c.leds.channel[ch].strip_type = C.int(int(t))
}

func (c *WS281x) initializeController() {
	c.leds.freq = C.uint32_t(c.Config.Frequency)
	c.leds.dmanum = C.int(c.Config.DMA)
}

// Initialize initialize library, must be called once before other functions are
// called.
func (c *WS281x) Initialize() error {
	if resp := int(C.ws2811_init(c.leds)); resp != 0 {
		return fmt.Errorf("ws2811_init failed with code: %d", resp)
	}

	return nil
}

// Render update the display with the data from the LED buffer
func (c *WS281x) Render() error {
	if resp := int(C.ws2811_render(c.leds)); resp != 0 {
		return fmt.Errorf("ws2811_render failed with code: %d", resp)
	}

	return nil
}

// At return an Color which allows access to the LED display data as
// if it were a sequence of 24-bit RGB values.
func (c *WS281x) At(position int) color.Color {
	color := C.ws2811_get_led(c.leds, C.int(position))
	return uint32ToColor(uint32(color))
}

// Set set LED at position x,y to the provided 24-bit color value.
func (c *WS281x) Set(position int, color color.Color) {
	C.ws2811_set_led(c.leds, C.int(position), C.uint32_t(colorToUint32(color)))
}

func (c *WS281x) Close() error {
	if c.closed {
		return nil
	}

	c.closed = true
	C.ws2811_fini(c.leds)
	return nil
}

func colorToUint32(c color.Color) uint32 {
	// A color's RGBA method returns values in the range [0, 65535]
	red, green, blue, alpha := c.RGBA()

	return (alpha>>8)<<24 | (red>>8)<<16 | (green>>8)<<8 | blue>>8
}

func uint32ToColor(u uint32) color.Color {
	return color.RGBA{
		uint8(u>>16) & 255,
		uint8(u>>8) & 255,
		uint8(u>>0) & 255,
		uint8(u>>24) & 255,
	}
}

func btoi(b bool) int {
	if b {
		return 1
	}

	return 0
}
