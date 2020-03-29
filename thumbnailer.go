package main

import (
	"fmt"
	"os/exec"
)

// ThumbnailFile creates a file and puts it into ./thumbnail/
func ThumbnailFile(path string, hash string) {
	fmt.Printf("Creating thumbnail for: %s\n", path)
	ftype, err := GetFileType(path)
	if err != nil {
		return
	}
	ftype = Convert(ftype)
	switch ftype {
	case "png":
		if err := exec.Command(
			"ffmpeg",
			"-i", path,
			"-vf", "scale=w=150:h=150:force_original_aspect_ratio=decrease",
			config.ThumbnailFolder+hash+".jpg",
		).Start(); err != nil {
			fmt.Printf("Error" + err.Error())
			return
		}
	default:
		if err := exec.Command(
			"ffmpeg",
			"-i", path,
			"-vf", "scale=w=150:h=150:force_original_aspect_ratio=decrease",
			config.ThumbnailFolder+hash+".jpg",
		).Start(); err != nil {
			fmt.Printf("Error" + err.Error())
			return
		}

	}

}
