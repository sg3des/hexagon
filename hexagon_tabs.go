package main

import (
	"encoding/hex"
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
	sourcetag    *gtk.TextTag

	ascii       *gtk.TextView
	asciibuffer *gtk.TextBuffer
	asciitag    *gtk.TextTag
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
	t.asciibuffer.Connect("mark-set", t.FocusASCII)
	t.sourcebuffer.Connect("mark-set", t.FocusSource)
	t.sourcebuffer.Connect("changed", t.ChangedSource)
}

func (t *Tab) SetLineNumbers(linenums []string) {
	t.linesbuffer.SetText(strings.Join(linenums, "\n"))
}

func (t *Tab) SetASCII(text []string) {
	t.asciibuffer.SetText(strings.Join(text, "\n"))
}

func (t *Tab) FocusASCII() {
	var iter gtk.TextIter
	t.asciibuffer.GetIterAtMark(&iter, t.asciibuffer.GetInsert())

	if t.asciitag == nil {
		t.asciitag = t.asciibuffer.CreateTag("selected", map[string]string{"background": "#666", "foreground": "#fff"})
		t.sourcetag = t.sourcebuffer.CreateTag("selected", map[string]string{"background": "#666", "foreground": "#fff"})
	}

	t.RemoveTag("selected")

	var start, end gtk.TextIter
	t.asciibuffer.GetIterAtOffset(&end, iter.GetOffset()+1)
	t.asciibuffer.ApplyTag(t.asciitag, &iter, &end)

	row := iter.GetLine()
	col := iter.GetLineOffset() * 3
	if col+2 > 16*3 {
		return
	}
	t.sourcebuffer.GetIterAtLineOffset(&start, row, col)
	t.sourcebuffer.GetIterAtLineOffset(&end, row, col+2)
	t.sourcebuffer.ApplyTag(t.sourcetag, &start, &end)
}

func (t *Tab) FocusSource() {
	var iter gtk.TextIter
	t.sourcebuffer.GetIterAtMark(&iter, t.sourcebuffer.GetInsert())

	if t.sourcetag == nil {
		t.asciitag = t.asciibuffer.CreateTag("selected", map[string]string{"background": "#666", "foreground": "#fff"})
		t.sourcetag = t.sourcebuffer.CreateTag("selected", map[string]string{"background": "#666", "foreground": "#fff"})
	}

	t.RemoveTag("selected")
	offset := iter.GetOffset() - iter.GetOffset()%3

	var start, end gtk.TextIter
	t.sourcebuffer.GetIterAtOffset(&start, offset)
	t.sourcebuffer.GetIterAtOffset(&end, offset+2)
	t.sourcebuffer.ApplyTag(t.sourcetag, &start, &end)

	// row := start.GetLine()
	// col := start.GetLineOffset() / 3
	// t.asciibuffer.GetIterAtLineOffset(&start, row, col)
	// t.asciibuffer.GetIterAtLineOffset(&end, row, col+1)
	// t.asciibuffer.ApplyTag(t.asciitag, &start, &end)
}

func (t *Tab) RemoveTag(name string) {
	var start, end gtk.TextIter
	t.asciibuffer.GetIterAtOffset(&start, 0)
	t.asciibuffer.GetIterAtOffset(&end, t.asciibuffer.GetCharCount())
	t.asciibuffer.RemoveTagByName(name, &start, &end)

	t.sourcebuffer.GetIterAtOffset(&start, 0)
	t.sourcebuffer.GetIterAtOffset(&end, t.sourcebuffer.GetCharCount())
	t.sourcebuffer.RemoveTagByName(name, &start, &end)
}

func (t *Tab) ChangedSource() {
	log.Println("changed")
	text := t.GetText(false)

	lines := strings.Split(text, "\n")
	var ascii []string
	var linenums []string
	var lineoff int
	for _, line := range lines {
		line = strings.Replace(line, " ", "", -1)
		data, err := hex.DecodeString(line)
		if err != nil {
			log.Println("failed convert hex to data,", err)
			return
		}

		// asciiBytes := byteToASCII(data)

		ascii = append(ascii, byteToASCII(data))

		linenum := fmt.Sprintf(" %06x  ", lineoff)
		if len(data) == 0 {
			linenum = ""
		}
		linenums = append(linenums, linenum)
		lineoff += len(data)
	}

	// t.linesbuffer.SetText(linenums)
	t.SetLineNumbers(linenums)
	t.SetASCII(ascii)
	// t.asciibuffer.SetText(strings.Join(ascii, "\n"))
}

func (t *Tab) GetText(hiddenChars bool) string {
	var start gtk.TextIter
	var end gtk.TextIter

	t.sourcebuffer.GetStartIter(&start)
	t.sourcebuffer.GetEndIter(&end)
	return t.sourcebuffer.GetText(&start, &end, hiddenChars)
}
