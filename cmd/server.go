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
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>QR Code Generator</title>
<style>
body { 
    font-family: sans-serif; 
    background: #181818; 
    color: #fff; 
    display: flex; 
    flex-direction: column; 
    align-items: center;
    padding: 1rem;
    margin: 0;
}
form { 
    margin: 2em 0;
    width: 100%;
    max-width: 500px;
}
input[type=text] { 
    width: 100%;
    padding: 0.8em;
    margin: 0.5em 0;
    border-radius: 8px;
    border: none;
    box-sizing: border-box;
    font-size: 16px;
}
input[type=color] {
    width: 100%;
    padding: 0.5em;
    margin: 0.5em 0;
    border-radius: 8px;
    border: none;
    height: 50px;
}
button { 
    width: 100%;
    padding: 1em;
    margin-top: 1em;
    border-radius: 8px;
    border: none;
    background: #f54b37;
    color: #fff;
    font-weight: bold;
    cursor: pointer;
    font-size: 16px;
}
img { 
    margin-top: 2em;
    background: #000;
    padding: 1em;
    border-radius: 16px;
    max-width: 100%;
    height: auto;
}
label {
    display: block;
    margin: 1em 0 0.5em 0;
}
h1 {
    font-size: 1.8rem;
    text-align: center;
}
</style>
</head>
<body>
<h1>QR Code Generator</h1>
<form method="GET" action="/qrcode">
    <label for="url">URL:</label>
    <input type="text" id="url" name="url" value="https://twofoxout.com/albums/berlin" required>
    <label for="color">Color:</label>
    <input type="color" id="color" name="color" value="#f54b37">
    <label>Style: <select name="style">
        <option value="square">Square</option>
        <option value="rounded">Rounded</option>
        <option value="triangle">Triangle</option>
    </select></label>
    <button type="submit">Generate QR Code</button>
</form>
</body>
</html>
`)
}

func qrImageHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	colorHex := r.URL.Query().Get("color")
	style := r.URL.Query().Get("style")
	if url == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}
	if colorHex == "" {
		colorHex = "#f54b37"
	}
	if style == "" {
		style = "square"
	}
	img, err := GenerateQRCodeImage(url, colorHex, style)
	if err != nil {
		http.Error(w, "Failed to generate QR code: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, img)
	w.Write(buf.Bytes())
}
