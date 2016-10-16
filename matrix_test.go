package ws281x

import (
	"image/color"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MatrixSuite struct{}

var _ = Suite(&MatrixSuite{})

func (s *MatrixSuite) TestColorToUint32(c *C) {
	u := colorToUint32(color.RGBA{2, 42, 84, 168})

	c.Assert(int(u>>24&255), Equals, 168)
	c.Assert(int(u>>16&255), Equals, 2)
	c.Assert(int(u>>8&255), Equals, 42)
	c.Assert(int(u>>0&255), Equals, 84)

	c.Assert(uint32(u), Equals, uint32(0xa8022a54))
}

func (s *MatrixSuite) TestUint32ToColor(c *C) {
	color := uint32ToColor(uint32(0xa8022a54)).(color.RGBA)

	c.Assert(color.A, Equals, uint8(168))
	c.Assert(color.R, Equals, uint8(2))
	c.Assert(color.G, Equals, uint8(42))
	c.Assert(color.B, Equals, uint8(84))
}
