package resize

import (
	"bytes"
	"fmt"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/emote"
	"github.com/disintegration/imaging"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	maxVideoSize = 256 * 1024 // 256 KB
	maxFPS       = 30
	maxDuration  = 3.0
)

var numCPUs = runtime.NumCPU()

func FitEmote(emote *emote.EmoteData) error {
	if emote.Animated {
		resizedWebm, err := fitGIF(emote.File)
		if err != nil {
			return fmt.Errorf("error resizing emote: %w", err)
		}
		emote.File = resizedWebm
		return nil
	}

	resizedPng, err := fitPNG(emote.File)
	if err != nil {
		return fmt.Errorf("error resizing emote: %w", err)
	}
	emote.File = resizedPng

	return nil
}

func fitPNG(input []byte) ([]byte, error) {
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

func fitGIF(input []byte) ([]byte, error) {
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

	info, err := getVideoInfo(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get video info: %w", err)
	}

	fps := capFPS(info.FPS)
	duration := capDuration(info.Duration)

	encodingAttempts := encodingAttempts(duration)
	for i, config := range encodingAttempts {
		if err := runFFMPEG(fps, duration, inputPath, outputPath, config); err != nil {
			if i == len(encodingAttempts)-1 {
				return nil, fmt.Errorf("all ffmpeg attempts failed, last error: %w", err)
			}
			continue
		}

		output, err := os.ReadFile(outputPath)
		if err != nil {
			continue
		}
		if len(output) <= maxVideoSize {
			return output, nil
		}

		if i == len(encodingAttempts)-1 {
			return nil, fmt.Errorf("output exceeds 256KB after all attempts: %d bytes", len(output))
		}
	}

	return nil, fmt.Errorf("failed to create valid output")
}

type videoInfo struct {
	FPS      float64
	Duration float64
	Width    int
	Height   int
}

type encodingConfig struct {
	bitrate string
	crf     int
	cpuUsed int
}

func capFPS(fps float64) float64 {
	if fps > maxFPS || fps == 0 {
		return maxFPS
	}
	return fps
}

func capDuration(duration float64) float64 {
	if duration > maxDuration {
		return maxDuration
	}
	return duration
}

func encodingAttempts(duration float64) []encodingConfig {
	return []encodingConfig{
		{bitrate: calculateTargetBitrate(duration, 0.85), crf: 32, cpuUsed: 1},
		{bitrate: calculateTargetBitrate(duration, 0.75), crf: 35, cpuUsed: 2},
		{bitrate: calculateTargetBitrate(duration, 0.65), crf: 40, cpuUsed: 3},
		{bitrate: calculateTargetBitrate(duration, 0.55), crf: 45, cpuUsed: 4},
	}
}

func getVideoInfo(inputPath string) (*videoInfo, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-select_streams", "v:0",
		"-show_entries", "stream=r_frame_rate,width,height,duration",
		"-of", "csv=p=0",
		inputPath,
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	parts := strings.Split(strings.TrimSpace(string(out)), ",")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid ffprobe output: %s", out)
	}

	info := &videoInfo{}

	// FPS is a fraction
	fpsParts := strings.Split(parts[0], "/")
	if len(fpsParts) == 2 {
		num, _ := strconv.ParseFloat(fpsParts[0], 64)
		den, _ := strconv.ParseFloat(fpsParts[1], 64)
		if den != 0 {
			info.FPS = num / den
		}
	}

	info.Width, _ = strconv.Atoi(parts[1])
	info.Height, _ = strconv.Atoi(parts[2])

	info.Duration, _ = strconv.ParseFloat(parts[3], 64)

	return info, nil
}

func calculateTargetBitrate(duration float64, efficiency float64) string {
	targetBits := float64(maxVideoSize) * efficiency * 8
	bitrate := int(targetBits / duration)
	return fmt.Sprintf("%dk", bitrate/1000)
}

func runFFMPEG(fps, duration float64, inputPath, outputPath string, config encodingConfig) error {
	threads := fmt.Sprintf("%d", min(numCPUs, 4)) // cap at 4 threads

	cmd := exec.Command("ffmpeg",
		"-y",
		"-i", inputPath,
		"-t", fmt.Sprintf("%.2f", duration),
		"-r", fmt.Sprintf("%.0f", fps),
		"-vf", "scale='if(gt(a,1),512,-1)':'if(gt(a,1),-1,512)'",
		"-c:v", "libvpx-vp9",
		"-pix_fmt", "yuva420p",
		"-crf", fmt.Sprintf("%d", config.crf),
		"-b:v", config.bitrate,
		"-maxrate", config.bitrate,
		"-bufsize", fmt.Sprintf("%dk", parseBitrate(config.bitrate)*2),
		"-an", // No audio
		"-threads", threads,
		"-row-mt", "1",
		"-tile-columns", "2",
		"-quality", "good",
		"-cpu-used", fmt.Sprintf("%d", config.cpuUsed),
		"-static-thresh", "0",
		"-auto-alt-ref", "1",
		"-lag-in-frames", "16",
		"-f", "webm",
		outputPath,
	)

	return cmd.Run()
}

func parseBitrate(bitrateStr string) int {
	numStr := strings.TrimSuffix(bitrateStr, "k")
	val, _ := strconv.Atoi(numStr)
	return val
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
