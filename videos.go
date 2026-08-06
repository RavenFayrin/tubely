package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	var data struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}
	videoJSON := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var output bytes.Buffer
	videoJSON.Stdout = &output
	err := videoJSON.Run()
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(output.Bytes(), &data)
	if err != nil {
		return "", err
	}
	if len(data.Streams) < 1 {
		return "", errors.New("JSON unmarshaling failure")
	}

	width := data.Streams[0].Width
	height := data.Streams[0].Height

	if width == 16*height/9 {
		return "16:9", nil
	} else if height == 16*width/9 {
		return "9:16", nil
	}
	return "other", nil
}