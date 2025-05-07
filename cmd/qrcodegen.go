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

// GenerateQRCodeImage generates a QR code image for the given URL and color hex
func GenerateQRCodeImage(url, colorHex string) (image.Image, error) {
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
	// Draw QR code modules with random rounded corners
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if bitmap[y][x] {
				cornerRadius := (rand.Intn(3) + 1) * scale
				drawRoundedRect(img, x*8*scale, y*8*scale, 8*scale, 8*scale, cornerRadius, customColor)
			}
		}
	}
	return img, nil
}
