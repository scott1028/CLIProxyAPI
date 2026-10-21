package util

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/png"
	"sync"
)

// whiteImageCache stores encoded placeholders by size; only a few sizes exist.
var whiteImageCache sync.Map

func CreateWhiteImageBase64(aspectRatio string) (string, error) {
	width := 1024
	height := 1024

	switch aspectRatio {
	case "1:1":
		width = 1024
		height = 1024
	case "2:3":
		width = 832
		height = 1248
	case "3:2":
		width = 1248
		height = 832
	case "3:4":
		width = 864
		height = 1184
	case "4:3":
		width = 1184
		height = 864
	case "4:5":
		width = 896
		height = 1152
	case "5:4":
		width = 1152
		height = 896
	case "9:16":
		width = 768
		height = 1344
	case "16:9":
		width = 1344
		height = 768
	case "21:9":
		width = 1536
		height = 672
	}

	size := image.Pt(width, height)
	if cached, ok := whiteImageCache.Load(size); ok {
		return cached.(string), nil
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

	var buf bytes.Buffer

	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())
	whiteImageCache.Store(size, base64String)
	return base64String, nil
}
