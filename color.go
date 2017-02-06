// color.go
package main

import (
	"fmt"
	"image/color"
	"math"
)

type Color struct {
	r, g, b float64
}

func (c Color) Add(u Color) Color {
	return Color{c.r + u.r, c.g + u.g, c.b + u.b}
}

func (c Color) Mul(f float64) Color {
	return Color{c.r * f, c.g * f, c.b * f}
}

func (c Color) String() string {
	return fmt.Sprintf("<Col: %.2f %.2f %.2f>", c.r, c.g, c.b)
}

func (c Color) ToPixel() color.RGBA {
	c.r = math.Max(0.0, math.Min(c.r*255.0, 255.0))
	c.g = math.Max(0.0, math.Min(c.g*255.0, 255.0))
	c.b = math.Max(0.0, math.Min(c.b*255.0, 255.0))
	return color.RGBA{uint8(c.r), uint8(c.g), uint8(c.b), 255}
}
