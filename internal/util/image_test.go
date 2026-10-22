package util

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"testing"
)

func TestCreateWhiteImageBase64UsesAspectRatioSize(t *testing.T) {
	tests := []struct {
		aspectRatio string
		width       int
		height      int
	}{
		{aspectRatio: "16:9", width: 1344, height: 768},
		{aspectRatio: "2:3", width: 832, height: 1248},
		{aspectRatio: "unknown", width: 1024, height: 1024},
	}

	for _, tt := range tests {
		t.Run(tt.aspectRatio, func(t *testing.T) {
			encoded, err := CreateWhiteImageBase64(tt.aspectRatio)
			if err != nil {
				t.Fatalf("CreateWhiteImageBase64() error = %v", err)
			}
			raw, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				t.Fatalf("DecodeString() error = %v", err)
			}
			config, err := png.DecodeConfig(bytes.NewReader(raw))
			if err != nil {
				t.Fatalf("DecodeConfig() error = %v", err)
			}
			if config.Width != tt.width || config.Height != tt.height {
				t.Fatalf("size = %dx%d, want %dx%d", config.Width, config.Height, tt.width, tt.height)
			}
		})
	}
}

func TestCreateWhiteImageBase64ReturnsCachedValue(t *testing.T) {
	first, err := CreateWhiteImageBase64("21:9")
	if err != nil {
		t.Fatalf("first CreateWhiteImageBase64() error = %v", err)
	}
	second, err := CreateWhiteImageBase64("21:9")
	if err != nil {
		t.Fatalf("second CreateWhiteImageBase64() error = %v", err)
	}
	if first != second {
		t.Fatal("cached image differs from first encoded image")
	}
	if _, ok := whiteImageCache.Load(image.Pt(1536, 672)); !ok {
		t.Fatal("21:9 image was not cached")
	}
}
