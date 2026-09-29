//go:build windows
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	app, _ := os.Executable()
	root := filepath.Dir(app)
	ffmpeg := filepath.Join(root, "ffmpeg", "ffmpeg.exe")
	if _, err := os.Stat(ffmpeg); err != nil {
		fmt.Println("Thieu ffmpeg\\ffmpeg.exe")
		fmt.Println("Hay dung goi Portable duoc GitHub Actions build.")
		return
	}

	fmt.Println("THINH AUDIO TOOL - WINDOWS")
	fmt.Println("1) Chuyen file local")
	fmt.Println("2) Chuyen URL audio truc tiep http/https")
	fmt.Print("Chon: ")
	in := bufio.NewReader(os.Stdin)
	choice, _ := in.ReadString('\n')
	choice = strings.TrimSpace(choice)

	fmt.Print("File/URL: ")
	src, _ := in.ReadString('\n')
	src = strings.TrimSpace(src)
	if src == "" { return }

	fmt.Println("1) MP3 320k  2) MP3 256k  3) WAV 44.1kHz  4) M4A AAC 256k")
	fmt.Print("Dinh dang: ")
	f, _ := in.ReadString('\n')
	f = strings.TrimSpace(f)

	outDir := filepath.Join(root, "output")
	_ = os.MkdirAll(outDir, 0755)

	base := filepath.Base(src)
	if choice == "2" { base = "audio_url" }
	base = strings.TrimSuffix(base, filepath.Ext(base))
	if base == "" || base == "." { base = "audio" }

	ext := ".mp3"
	args := []string{"-hide_banner", "-loglevel", "error", "-y"}
	if choice == "2" {
		if !strings.HasPrefix(src, "http://") && !strings.HasPrefix(src, "https://") {
			fmt.Println("URL phai bat dau bang http:// hoac https://")
			return
		}
	}
	args = append(args, "-i", src, "-vn", "-ar", "44100")

	switch f {
	case "2":
		args = append(args, "-c:a", "libmp3lame", "-b:a", "256k")
	case "3":
		ext = ".wav"
		args = append(args, "-c:a", "pcm_s16le")
	case "4":
		ext = ".m4a"
		args = append(args, "-c:a", "aac", "-b:a", "256k", "-movflags", "+faststart")
	default:
		args = append(args, "-c:a", "libmp3lame", "-b:a", "320k")
	}

	out := filepath.Join(outDir, base+ext)
	for i := 2; ; i++ {
		if _, err := os.Stat(out); os.IsNotExist(err) { break }
		out = filepath.Join(outDir, fmt.Sprintf("%s_%d%s", base, i, ext))
	}
	args = append(args, out)

	fmt.Println("Dang xu ly...")
	cmd := exec.Command(ffmpeg, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Loi:", err)
		return
	}
	fmt.Println("Xong:", out)
}
