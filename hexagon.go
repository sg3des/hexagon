package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	arg "github.com/alexflint/go-arg"
	"github.com/mattn/go-gtk/gtk"
)

var (
	ui *UI

	args struct {
		Files []string `arg:"positional"`
	}
)

func init() {
	log.SetFlags(log.Lshortfile)
	arg.MustParse(&args)
	if len(args.Files) == 0 {
		args.Files = append(args.Files, "")
	}
}

func main() {
	gtk.Init(nil)
	ui = CreateUI()

	for _, filename := range args.Files {
		NewTab(filename)
	}

	gtk.Main()
}

func byteToHex(data []byte) ([]string, string, []string) {
	var lines []string
	var nums []string
	var ascii []string

	reader := bytes.NewReader(data)

	var line = make([]byte, 16)
	for {
		n, err := reader.Read(line)
		if n == 0 {
			break
		}
		line = line[:n]

		lines = append(lines, fmt.Sprintf("% x", line))
		nums = append(nums, fmt.Sprintf(" %06x  ", len(nums)*16))
		ascii = append(ascii, byteToASCII(line))

		if err != nil {
			break
		}
	}

	return nums, strings.Join(lines, "\n"), ascii
}

func byteToASCII(data []byte) string {
	var s string
	for _, b := range data {
		if b < 32 || b > 126 {
			s += "."
		} else {
			s += string(b)
		}
	}

	return s
}

func tabsContains(filename string) bool {
	for n, t := range tabs {
		if t.filename == filename {
			ui.notebook.SetCurrentPage(n)
			return true
		}
	}
	return false
}
