package main

import (
	"fmt"
	"os/exec"
)

// ThumbnailFile creates a file and puts it into ./thumbnail/
func ThumbnailFile(path string, hash string) {
	fmt.Printf("Creating thumbnail for: %s\n", path)
	ftype := GetFileType(path)
	ftype = Convert(ftype)
	switch ftype {
	case "png":
		Cmd := exec.Command(
			"ffmpeg",
			"-i", path,
			"-vf", "scale=w=150:h=150:force_original_aspect_ratio=decrease",
			config.ThumbnailFolder+hash+".jpg",
		)
		err := Cmd.Run()
		Cmd.Process.Kill()
		if err != nil {
			fmt.Printf("Error" + err.Error())
			return
		}
	default:
		Cmd := exec.Command(
			"ffmpeg",
			"-i", path,
			"-vf", "scale=w=150:h=150:force_original_aspect_ratio=decrease",
			config.ThumbnailFolder+hash+".jpg",
		)
		err := Cmd.Run()
		Cmd.Process.Kill()
		if err != nil {
			fmt.Printf("Error" + err.Error())
			return
		}

	}

}
