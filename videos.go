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

	aspect1, aspect2 := GetReducedRatio(data.Streams[0].Width, data.Streams[0].Height)
	if aspect1 == 16 && aspect2 == 9 {
		return "16:9", nil
	}
	if aspect1 == 9 && aspect2 == 16 {
		return "9:16", nil
	}
	return "other", nil
}

func GetReducedRatio(w, h int) (int, int) {
	d := greatestCommonDivisor(w, h)
	return w / d, h / d
}

func greatestCommonDivisor(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
