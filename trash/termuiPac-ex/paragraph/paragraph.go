package main

import (
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func renderWidgets(p *widgets.Paragraph) {
	w, _ := ui.TerminalDimensions()
	centerx1 := (w / 2) - (len(p.Text) / 2)
	centerx2 := (w / 2) + (len(p.Text) / 2) + 3
	p.SetRect(centerx1, 0, centerx2, 5)
	time.Sleep(time.Millisecond)
	ui.Clear()
	ui.Render(p)
}

func showNumbers(p *widgets.Paragraph) {
	for i := 0; i < 100000; i++ {
		p.Text = fmt.Sprint(i)
		time.Sleep(time.Millisecond)
		renderWidgets(p)
	}

}

type input struct {
	scan rune
}

func trimFunc() {

}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// var userInput *input
	w, _ := ui.TerminalDimensions()

	text := "lorem Lorem ipsum dolor sit amet. Ab sapiente molestiae qui porro debitis qui suscipit dolores."

	centerx1 := (w / 2) - (len(text) / 2)
	centerx2 := (w / 2) + (len(text) / 2) + 3
	p0 := widgets.NewParagraph()
	p0.Text = text
	// p0.SetRect(0, 0, 20, 5)
	p0.SetRect(centerx1, 0, centerx2, 5)
	p0.Border = true
	// p0.PaddingLeft = w / 2

	// p1 := widgets.NewParagraph()
	// p1.Title = "标签"
	// p1.Text = "你好，世界。"
	// p1.SetRect(20, 0, 35, 5)

	// p2 := widgets.NewParagraph()
	// p2.Title = "Multiline"
	// p2.Text = "Simple colored text\nwith label. It [can be](fg:red) multilined with \\n or [break automatically](fg:red,fg:bold)"
	// p2.SetRect(0, 5, 35, 10)
	// p2.BorderStyle.Fg = ui.ColorYellow

	// p3 := widgets.NewParagraph()
	// p3.Title = "Auto Trim"
	// p3.Text = "Long text with label and it is auto trimmed."
	// p3.SetRect(0, 10, 40, 15)

	// p4 := widgets.NewParagraph()
	// p4.Title = "Text Box with Wrapping"
	// p4.Text = "Press q to QUIT THE DEMO. [There](fg:blue,mod:bold) are other things [that](fg:red) are going to fit in here I think. What do you think? Now is the time for all good [men to](bg:blue) come to the aid of their country. [This is going to be one really really really long line](fg:green) that is going to go together and stuffs and things. Let's see how this thing renders out.\n    Here is a new paragraph and stuffs and things. There should be a tab indent at the beginning of the paragraph. Let's see if that worked as well."
	// p4.SetRect(40, 0, 70, 20)
	// p4.BorderStyle.Fg = ui.ColorBlue

	// ui.Render(p0, p1, p2, p3, p4)
	// ch := make(chan *widgets.Paragraph)

	// go func() {
	// 	for {
	// 		wigPar := <-ch
	// 		ui.Render(wigPar)

	// 	}
	// }()

	ui.Render(p0)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "s":
			// if userInput == nil {
			// 	userInput = &input{scan: rune('e')}
			// 	log.Printf("userInput: %v\n", userInput)
			// 	showNumbers(p0)
			// }
			showNumbers(p0)
		case "<Resize>":
			// w, _ = ui.TerminalDimensions()
			// centerx1 = (w / 2) - (len(text) / 2)
			// centerx2 = (w / 2) + (len(text) / 2) + 3
			// p0.SetRect(centerx1, 0, centerx2, 5)
			// ui.Clear()
			// ui.Render(p0)
			renderWidgets(p0)
		}
	}
}
