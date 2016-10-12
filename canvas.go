package ws281x

import (
	"image"
	"image/color"
	"image/draw"
)

type Canvas struct {
	w, h   int
	m      Matrix
	closed bool
}

func NewCanvas(w, h int, config *HardwareConfig) (*Canvas, error) {
	m, err := NewWS281x(w*h, config)
	if err != nil {
		return nil, err
	}

	return &Canvas{
		w: w,
		h: h,
		m: m,
	}, nil
}

// Initialize initialize the matrix and the canvas
func (c *Canvas) Initialize() error {
	return c.m.Initialize()
}

// Render update the display with the data from the LED buffer
func (c *Canvas) Render() error {
	return c.m.Render()
}

func (c *Canvas) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds return the topology of the Matrix
func (c *Canvas) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.w, c.h)
}

// At return an Color which allows access to the LED display data as
// if it were a sequence of 24-bit RGB values.
func (c *Canvas) At(x, y int) color.Color {
	return c.m.At(c.position(x, y))
}

// Set set LED at position x,y to the provided 24-bit color value.
func (c *Canvas) Set(x, y int, color color.Color) {
	c.m.Set(c.position(x, y), color)
}

func (c *Canvas) position(x, y int) int {
	return x + (y * c.w)
}

func (c *Canvas) Clear() error {
	draw.Draw(c, c.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)
	return c.m.Render()
}

func (c *Canvas) Close() error {
	c.Clear()

	return c.m.Close()
}

type Matrix interface {
	Initialize() error
	At(position int) color.Color
	Set(position int, c color.Color)
	Render() error
	Close() error
}
