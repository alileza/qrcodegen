package cmd

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"

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
			&cli.StringFlag{
				Name:    "style",
				Usage:   "Module style: square, rounded, triangle (default: square)",
				Value:   "square",
				Aliases: []string{"s"},
			},
		},
		Action: func(c *cli.Context) error {
			url := c.String("url")
			outputPath := c.String("output")
			colorHex := c.String("color")
			style := c.String("style")

			img, err := GenerateQRCodeImage(url, colorHex, style)
			if err != nil {
				return fmt.Errorf("failed to generate QR code: %w", err)
			}

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

			fmt.Printf("QR code generated successfully: %s\n", outputPath)
			return nil
		},
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
