package cmd

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/urfave/cli/v2"
)

// GenerateCommand returns the generate command
func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "Generate a QR code from a URL with a centered logo overlay.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Usage:    "URL to generate QR code for",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Usage:   "Output file path (default: qrcode.png)",
				Value:   "qrcode.png",
				Aliases: []string{"o"},
			},
			&cli.StringFlag{
				Name:    "color",
				Usage:   "Hex color for QR code modules (e.g. #ffffff)",
				Value:   "#ffffff",
				Aliases: []string{"c"},
			},
		},
		Action: func(c *cli.Context) error {
			url := c.String("url")
			outputPath := c.String("output")
			colorHex := c.String("color")

			// Initialize random seed
			rand.Seed(time.Now().UnixNano())

			// Create QR code with highest error correction
			qr, err := qrcode.New(url, qrcode.Highest)
			if err != nil {
				return fmt.Errorf("failed to create QR code: %w", err)
			}

			// Scale factor for high resolution
			scale := 10

			bitmap := qr.Bitmap()
			size := len(bitmap)
			imgSize := size * 8 * scale

			img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))
			// Fill background with black
			draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
			// Set QR code module color to user-specified color
			customColor, err := parseHexColor(colorHex)
			if err != nil {
				return fmt.Errorf("invalid color hex: %v", err)
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

			// No logo overlay, just the QR code

			outputDir := filepath.Dir(outputPath)
			if outputDir != "." {
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf("failed to create output directory: %w", err)
				}
			}

			f, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer f.Close()

			if err := png.Encode(f, img); err != nil {
				return fmt.Errorf("failed to save QR code: %w", err)
			}

			fmt.Printf("QR code with centered logo generated successfully: %s\n", outputPath)
			return nil
		},
	}
}

// drawRoundedRect draws a rounded rectangle on the image
func drawRoundedRect(img *image.RGBA, x, y, width, height, radius int, c color.Color) {
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

// loadAndResizeLogo loads logo.png and resizes it to the given width and height
func loadAndResizeLogo(path string, width, height int) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	// Resize (nearest neighbor)
	resized := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := x * img.Bounds().Dx() / width
			srcY := y * img.Bounds().Dy() / height
			resized.Set(x, y, img.At(srcX, srcY))
		}
	}
	return resized, nil
}

// toColorSilhouette converts an image to a colored silhouette, preserving alpha
func toColorSilhouette(src image.Image, col color.Color) image.Image {
	bounds := src.Bounds()
	c := color.NRGBAModel.Convert(col).(color.NRGBA)
	coloredLogo := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := src.At(x, y).RGBA()
			coloredLogo.Set(x, y, color.NRGBA{c.R, c.G, c.B, uint8(a >> 8)})
		}
	}
	return coloredLogo
}

// circleMask implements image.Image for a circular mask
// Used to draw a white circle background for the logo
// Center is at (radius, radius)
type circleMask struct{ radius int }

func (m *circleMask) ColorModel() color.Model { return color.AlphaModel }
func (m *circleMask) Bounds() image.Rectangle { return image.Rect(0, 0, m.radius*2, m.radius*2) }
func (m *circleMask) At(x, y int) color.Color {
	dx := x - m.radius
	dy := y - m.radius
	if dx*dx+dy*dy <= m.radius*m.radius {
		return color.Alpha{A: 255}
	}
	return color.Alpha{A: 0}
}

// parseHexColor parses a hex color string (e.g. #fff or #ffffff) to color.Color
func parseHexColor(s string) (color.Color, error) {
	s = strings.TrimPrefix(s, "#")
	var r, g, b uint8
	if len(s) == 6 {
		_, err := fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b)
		if err != nil {
			return nil, err
		}
		return color.RGBA{r, g, b, 255}, nil
	} else if len(s) == 3 {
		_, err := fmt.Sscanf(s, "%1x%1x%1x", &r, &g, &b)
		if err != nil {
			return nil, err
		}
		return color.RGBA{r * 17, g * 17, b * 17, 255}, nil
	}
	return nil, fmt.Errorf("invalid hex color length")
}
