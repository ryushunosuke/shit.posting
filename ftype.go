package main

import (
	"bytes"
	"os/exec"
)

// GetFileType returns the file type of the probed content.
// Doesn't return png, webm, jpg or anything like that. Do your own matching for your case.
func GetFileType(Path string) (Type string, err error) {
	Cmd := exec.Command("ffprobe", "-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "format=format_long_name",
		"-of", "default=nokey=1:noprint_wrappers=1",
		Path)
	var OutBuf bytes.Buffer
	Cmd.Stdout = &OutBuf
	err = Cmd.Start()
	if err != nil {
		panic(err)
	}
	Cmd.Wait()
	if err != nil {
		panic(err)
	}
	s := OutBuf.Bytes()
	Type = string(s[:])
	return Type, err

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

	}
	return
}
