package cmd

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"

	"github.com/urfave/cli/v2"
)

func ServerCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start a web server to generate QR codes via a simple HTML form.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Usage: "Address to listen on",
				Value: ":8080",
			},
		},
		Action: func(c *cli.Context) error {
			addr := c.String("addr")
			http.HandleFunc("/", qrFormHandler)
			http.HandleFunc("/qrcode", qrImageHandler)
			fmt.Printf("Server running at http://localhost%s\n", addr)
			return http.ListenAndServe(addr, nil)
		},
	}
}

func qrFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>QR Code Generator</title>
<style>
body { font-family: sans-serif; background: #181818; color: #fff; display: flex; flex-direction: column; align-items: center; }
form { margin: 2em 0; }
input[type=text], input[type=color] { padding: 0.5em; margin: 0.5em; border-radius: 4px; border: none; }
button { padding: 0.5em 1em; border-radius: 4px; border: none; background: #f54b37; color: #fff; font-weight: bold; cursor: pointer; }
img { margin-top: 2em; background: #000; padding: 1em; border-radius: 16px; }
</style>
</head>
<body>
<h1>QR Code Generator</h1>
<form method="GET" action="/qrcode">
	<label>URL: <input type="text" name="url" value="https://twofoxout.com/albums/berlin" size="40" required></label><br>
	<label>Color: <input type="color" name="color" value="#f54b37"></label><br>
	<button type="submit">Generate QR Code</button>
</form>
</body>
</html>
`)
}

func qrImageHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	colorHex := r.URL.Query().Get("color")
	if url == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}
	if colorHex == "" {
		colorHex = "#f54b37"
	}
	img, err := GenerateQRCodeImage(url, colorHex)
	if err != nil {
		http.Error(w, "Failed to generate QR code: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, img)
	w.Write(buf.Bytes())
}
