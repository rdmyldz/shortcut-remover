package main

import (
	"log"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func getLongText(char string, num int) string {
	var s strings.Builder
	for i := 0; i < num; i++ {
		s.WriteString(char)
	}
	return s.String()
}

func writeTitle(p *widgets.Paragraph) {
	// p.Text = "Hello World!"
	p.Text = getLongText("s", 25)
	p.PaddingLeft = 20
	p.SetRect(0, 0, 120, 65)
}

// func useChan(p *widgets.Paragraph) {
// 	ch := make(chan *widgets.Paragraph)

// }

func renderDrawables(changed <-chan bool, drawables ...ui.Drawable) {

	for range changed {
		ui.Render(drawables...)
	}
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	changed := make(chan bool)
	defer close(changed)
	go renderDrawables(changed)
	title := widgets.NewParagraph()

	title.Title = "another title"
	title.Text = "Hello World!"
	title.Text = getLongText("s", 25)
	title.PaddingLeft = 20
	title.SetRect(0, 0, 120, 65)

	// ui.Render(title)
	changed <- true

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}
}
