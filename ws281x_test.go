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
	color := color.RGBA{2, 42, 84, 168}

	u := colorToUint32(color)

	c.Assert(int(u>>24&255), Equals, 168)
	c.Assert(int(u>>16&255), Equals, 2)
	c.Assert(int(u>>8&255), Equals, 42)
	c.Assert(int(u>>0&255), Equals, 84)

	c.Assert(uint32(u), Equals, uint32(0xa8022a54))
}
