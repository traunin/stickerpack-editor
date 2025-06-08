package resize

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	// "io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Traunin/stickerpack-editor/apps/api/emote"
	"github.com/disintegration/imaging"
)

func FitEmote(emote *emote.EmoteData) error {
	if emote.Animated {
		resizedWebm, err := fitGif(emote.File)
		if err != nil {
			return fmt.Errorf("Error resizing emote: %w", err)
		}
		emote.File = resizedWebm
		return nil
	}

	resizedPng, err := fitPng(emote.File)
	if err != nil {
		return fmt.Errorf("Error resizing emote: %w", err)
	}
	emote.File = resizedPng

	return nil
}

func fitPng(input []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	srcBounds := img.Bounds()
	width := srcBounds.Dx()
	height := srcBounds.Dy()

	var newImg *image.NRGBA
	if width >= height {
		newImg = imaging.Resize(img, 512, 0, imaging.Lanczos)
	} else {
		newImg = imaging.Resize(img, 0, 512, imaging.Lanczos)
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, newImg)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func fitGif(input []byte) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "gifconv")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.gif")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, fmt.Errorf("failed to write input file: %w", err)
	}

	outputPath := filepath.Join(tmpDir, "output.webm")

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-c:v", "libvpx-vp9",
		"-pix_fmt", "yuva420p",
		"-f", "webm",
		"-an",          // No audio
		"-row-mt", "1", // Multi-threading
		"-crf", "40", // Quality, increase if Durov complains about file size
		"-b:v", "0", // Variable bitrate
		"-t", "3", // Max duration 3 seconds
		"-vf", "scale=512:-1", // Simplified scaling (width=512, maintain aspect)
		"-auto-alt-ref", "0", // Better for transparent videos
		"-quality", "good", // Quality preset
		"-cpu-used", "4", // Faster encoding
		outputPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %v\n%s", err, stderr.String())
	}

	output, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read output file: %w", err)
	}

	return output, nil
}
