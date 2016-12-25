package main

import (
	"log"

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
