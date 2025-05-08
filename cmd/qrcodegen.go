package cmd

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"time"

	"github.com/skip2/go-qrcode"
)

// drawTriangle draws a filled upward-pointing triangle for each QR code module
func drawTriangle(img *image.RGBA, x, y, width, height int, c color.Color) {
	// Vertices: top (x+width/2, y), bottom left (x, y+height), bottom right (x+width, y+height)
	for dy := 0; dy < height; dy++ {
		left := x + width/2 - (width*dy)/(2*height)
		right := x + width/2 + (width*dy)/(2*height)
		for dx := left; dx <= right; dx++ {
			img.Set(dx, y+dy, c)
		}
	}
}

// isInFinderPattern returns true if (x, y) is inside a finder pattern area
func isInFinderPattern(x, y, size int) bool {
	finderSize := 7
	margin := 0
	// Top-left
	if x >= margin && x < finderSize+margin && y >= margin && y < finderSize+margin {
		return true
	}
	// Top-right
	if x >= size-finderSize-margin && x < size-margin && y >= margin && y < finderSize+margin {
		return true
	}
	// Bottom-left
	if x >= margin && x < finderSize+margin && y >= size-finderSize-margin && y < size-margin {
		return true
	}
	return false
}

// drawSquare draws a filled square for a QR code module
func drawSquare(img *image.RGBA, x, y, width, height int, c color.Color) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			img.Set(x+dx, y+dy, c)
		}
	}
}

// isInTimingPattern returns true if (x, y) is in a timing pattern (row 6 or column 6, but not in finder pattern)
func isInTimingPattern(x, y, size int) bool {
	// Exclude finder pattern areas
	if isInFinderPattern(x, y, size) {
		return false
	}
	// Timing patterns are at row 6 and column 6
	return x == 6 || y == 6
}

// drawSquareCentered draws a filled square, centered in the module cell
func drawSquareCentered(img *image.RGBA, x, y, cellSize, squareSize int, c color.Color) {
	offset := (cellSize - squareSize) / 2
	drawSquare(img, x+offset, y+offset, squareSize, squareSize, c)
}

// GenerateQRCodeImage generates a QR code image for the given URL and color hex
func GenerateQRCodeImage(url, colorHex, style string) (image.Image, error) {
	rand.Seed(time.Now().UnixNano())
	qr, err := qrcode.New(url, qrcode.Highest)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}
	bitmap := qr.Bitmap()
	size := len(bitmap)
	scale := 10
	imgSize := size * 8 * scale
	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))
	// Fill background with black
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	// Parse color
	customColor, err := parseHexColor(colorHex)
	if err != nil {
		return nil, fmt.Errorf("invalid color hex: %v", err)
	}
	// Draw QR code modules as larger squares for finder/timing patterns and style for others
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if bitmap[y][x] {
				if isInFinderPattern(x, y, size) || isInTimingPattern(x, y, size) {
					drawSquareCentered(img, x*8*scale, y*8*scale, 8*scale, 10*scale, customColor)
				} else {
					switch style {
					case "triangle":
						drawTriangle(img, x*8*scale, y*8*scale, 8*scale, 8*scale, customColor)
					case "rounded":
						drawRoundedSquare(img, x*8*scale, y*8*scale, 8*scale, 8*scale, 3*scale, customColor)
					default:
						drawSquare(img, x*8*scale, y*8*scale, 8*scale, 8*scale, customColor)
					}
				}
			}
		}
	}
	return img, nil
}

// drawRoundedSquare draws a filled rounded square for a QR code module
func drawRoundedSquare(img *image.RGBA, x, y, width, height, radius int, c color.Color) {
	draw.Draw(img, image.Rect(x+radius, y, x+width-radius, y+height), &image.Uniform{c}, image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(x, y+radius, x+width, y+height-radius), &image.Uniform{c}, image.Point{}, draw.Src)
	for i := 0; i < radius; i++ {
		for j := 0; j < radius; j++ {
			if (i-radius)*(i-radius)+(j-radius)*(j-radius) <= radius*radius {
				img.Set(x+i, y+j, c)
				img.Set(x+width-1-i, y+j, c)
				img.Set(x+i, y+height-1-j, c)
				img.Set(x+width-1-i, y+height-1-j, c)
			}
		}
	}
}
