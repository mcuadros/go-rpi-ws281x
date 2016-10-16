# go-rpi-ws281x [![GoDoc](https://godoc.org/github.com/mcuadros/go-rpi-ws281x?status.svg)](http://godoc.org/github.com/mcuadros/go-rpi-ws281x) 

golang binding for [rpi_ws281x](https://github.com/jgarff/rpi_ws281x), userspace Raspberry Pi PWM library for WS281X LEDs. Supports any Raspberry and  WS2812, SK6812RGB and [SK6812RGBW](https://www.adafruit.com/category/168?q=SK6812RGBW&) LEDs strips, this includes [Unicorn pHAT](https://shop.pimoroni.com/products/unicorn-phat) and [NeoPixels](https://www.adafruit.com/category/168)


Installation
------------

The recommended way to install `go-rpi-ws281x` is:

```sh
go get github.com/mcuadros/go-rpi-ws281x
cd $GOPATH/src/github.com/mcuadros/go-rpi-ws281x/vendor/rpi_ws281x
scons
```

Requires having install golang and scons (on raspbian, apt-get install scons).

Examples
--------

```go
// create a new canvas with the given width and height, and the config, in this
// case the configuration is for a Unicorn pHAT (8x4 pixels matrix) with the
// default configuration
c, _ := ws281x.NewCanvas(8, 4, &ws281x.DefaultConfig)


// initialize the canvas and the matrix
c.Initialize()

// since ws281x implements image.Image any function like draw.Draw from the std
// library may be used with it.
// 
// now we copy a white image into the ws281x.Canvas, this turn on all the leds
// to white
draw.Draw(c, c.Bounds(), image.NewUniform(color.White), image.ZP, draw.Over)

// render and sleep to see the leds on
c.Render()
time.Sleep(time.Second * 5)

// don't forget close the canvas, if not you leds may remain on
c.Close()
```

Check the folder [`examples`](https://github.com/mcuadros/go-rpi-ws281x/tree/master/examples) folder for more examples

License
-------

MIT, see [LICENSE](LICENSE)