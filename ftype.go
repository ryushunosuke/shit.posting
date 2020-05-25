package main

import (
	"bytes"
	"log"
	"os/exec"
)

// GetFileType returns the file type of the probed content.
// Doesn't return png, webm, jpg or anything like that. Do your own matching for your case.
func GetFileType(Path string) (Type string) {
	Cmd := exec.Command("ffprobe", "-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "format=format_long_name",
		"-of", "default=nokey=1:noprint_wrappers=1",
		Path)
	var OutBuf bytes.Buffer
	Cmd.Stdout = &OutBuf
	err := Cmd.Run()
	if err != nil {

		log.Println("File ", Path, "lead to ", err.Error())
		return ""
	}
	s := OutBuf.Bytes()
	Type = string(s[:])
	Cmd.Stdout = nil

	Cmd.Process.Kill()
	return Type

}

// Convert is used to converting what ffprobe shits out.
func Convert(s string) (Type string) {
	switch s {
	case "Matroska / WebM\n":
		Type = "webm"
	case "image2 sequence\n":
		Type = "jpg"
	case "piped png sequence\n":
		Type = "png"
	case "QuickTime / MOV\n":
		Type = "mp4"
	case "Animated Computer Image Graphic (GIF)\n":
		Type = "gif"
	default:
		Type = ""
	}
	return
}
