package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
	gsv "github.com/mattn/go-gtk/gtksourceview"
)

var (
	newtabiter = 0
	tabs       []*Tab
)

type Tab struct {
	filename string
	label    *gtk.Label

	// hbox        *gtk.HBox
	lines       *gtk.TextView
	linesbuffer *gtk.TextBuffer

	// swin         *gtk.ScrolledWindow
	// scrollbar    *gtk.Adjustment
	source       *gsv.SourceView
	sourcebuffer *gsv.SourceBuffer

	ascii       *gtk.TextView
	asciibuffer *gtk.TextBuffer
}

func NewTab(filename string) {
	if tabsContains(filename) {
		return
	}

	var newfile bool
	if filename == "" {
		filename = fmt.Sprintf("new-%d", newtabiter)
		newtabiter++
		newfile = true
	}

	t := &Tab{
		filename: filename,
	}

	t.lines = gtk.NewTextView()
	t.lines.SetEditable(false)
	t.lines.SetCursorVisible(false)
	t.lines.SetState(gtk.STATE_INSENSITIVE)
	t.lines.ModifyFontEasy("LiberationMono 8")
	t.lines.ModifyText(gtk.STATE_NORMAL, gdk.NewColor("grey"))
	t.linesbuffer = t.lines.GetBuffer()

	t.sourcebuffer = gsv.NewSourceBufferWithLanguage(gsv.SourceLanguageManagerGetDefault().GetLanguage("hex"))
	t.source = gsv.NewSourceViewWithBuffer(t.sourcebuffer)
	t.source.SetHighlightCurrentLine(true)
	t.source.ModifyFontEasy("LiberationMono 8")

	t.ascii = gtk.NewTextView()
	t.ascii.SetEditable(false)
	t.ascii.SetCursorVisible(true)
	t.ascii.ModifyFontEasy("LiberationMono 8")
	t.ascii.ModifyText(gtk.STATE_NORMAL, gdk.NewColor("grey"))
	t.asciibuffer = t.ascii.GetBuffer()

	if !newfile {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println(err)
			return
		}

		linenums, text, ascii := byteToHex(data)

		t.sourcebuffer.BeginNotUndoableAction()
		t.sourcebuffer.SetText(text)
		t.sourcebuffer.EndNotUndoableAction()

		t.SetLineNumbers(linenums)
		t.SetASCII(ascii)
	}

	scrollSource := gtk.NewScrolledWindow(nil, nil)
	scrollSource.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_NEVER)
	scrollSource.SetShadowType(gtk.SHADOW_IN)
	scrollSource.Add(t.source)

	scrollLines := gtk.NewScrolledWindow(nil, scrollSource.GetVAdjustment())
	scrollLines.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_NEVER)
	scrollLines.SetShadowType(gtk.SHADOW_IN)
	scrollLines.Add(t.lines)

	scrollASCII := gtk.NewScrolledWindow(nil, scrollSource.GetVAdjustment())
	scrollASCII.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_NEVER)
	scrollASCII.SetShadowType(gtk.SHADOW_IN)
	scrollASCII.Add(t.ascii)

	hbox := gtk.NewHBox(false, 0)
	hbox.PackStart(scrollLines, false, false, 0)
	hbox.PackStart(scrollSource, true, true, 0)
	hbox.PackEnd(scrollASCII, false, false, 0)

	swin := gtk.NewScrolledWindow(nil, scrollSource.GetVAdjustment())
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.Add(hbox)

	t.label = gtk.NewLabel(path.Base(filename))
	t.label.SetTooltipText(filename)

	n := ui.notebook.AppendPage(swin, t.label)
	ui.notebook.ShowAll()
	ui.notebook.SetCurrentPage(n)
	t.source.GrabFocus()

	tabs = append(tabs, t)

	t.ascii.Connect("button-press-event", t.FocusASCII)
	t.ascii.Connect("move-cursor", t.FocusASCII)

}

func (t *Tab) SetLineNumbers(linenums []string) {
	t.linesbuffer.SetText(strings.Join(linenums, "\n"))
}

func (t *Tab) SetASCII(text []string) {
	t.asciibuffer.SetText(strings.Join(text, "\n"))
}

func (t *Tab) FocusASCII() {
	var iter gtk.TextIter
	mark := t.asciibuffer.GetInsert()
	t.asciibuffer.GetIterAtMark(&iter, mark)
	log.Println(iter.GetOffset())
}
