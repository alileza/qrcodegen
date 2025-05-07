# QRCodeGen

A simple command-line and web tool to generate QR codes from URLs.

## Features
- Generate QR codes from the command line or a web interface
- Customizable QR code color (hex)
- High-resolution output (great for print)
- Modern, rounded QR code style

## Installation

```bash
go install qrcodegen@latest
```

## CLI Usage

Generate a QR code from a URL:

```bash
qrcodegen generate --url "https://example.com" --color "#f54b37" --output "qrcode.png"
```

### Options
- `--url` (required): The URL to generate a QR code for
- `--output` or `-o`: Output file path (default: qrcode.png)
- `--color` or `-c`: Hex color for QR code modules (default: #ffffff)

### Example
```bash
qrcodegen generate --url "https://twofoxout.com/albums/berlin" --color "#f54b37" --output "berlin-qr.png"
```

## Web Server Usage

Start the web server:

```bash
qrcodegen server --addr ":8080"
```

Then open [http://localhost:8080](http://localhost:8080) in your browser.

- Enter the URL and pick a color.
- Click "Generate QR Code" to view and download your QR code.

## Building from Source

1. Clone the repository
2. Run `go build`
3. The executable will be created in the current directory

---

**Made with Go, [urfave/cli](https://github.com/urfave/cli), and [skip2/go-qrcode](https://github.com/skip2/go-qrcode).** 