package main

import (
	"fmt"
	"bytes"
	"os/exec"
	"encoding/json"
)

type FFProbeOutput struct {
	Streams []struct {
		Width   int `json:"width"`
		Height  int `json:"height"`
	} `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filePath,
	)
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Internal server error", err)
	}

	data := buf.Bytes()
	var out FFProbeOutput
	err = json.Unmarshal(data, &out)
	if err != nil {
		return "", fmt.Errorf("couldn't parse ffprobe output: %w", err)
	}

	width := out.Streams[0].Width
	height := out.Streams[0].Height

	// 16:9-ish
	if width*9 >= height*16-1000 && width*9 <= height*16+1000 {
		return "16:9", nil
	}

	// 9:16-ish
	if height*9 >= width*16-1000 && height*9 <= width*16+1000 {
		return "9:16", nil
	}

	return "other", nil
	}
