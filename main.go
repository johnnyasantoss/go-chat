package main

import (
	"log"

	gc "github.com/rthornton128/goncurses"
)

func main() {
	stdscr, err := gc.Init()
	if err != nil {
		log.Fatal("Init window: ", err)
	}
	defer stdscr.Delete()

	stdscr.Keypad(true)
	stdscr.Timeout(0)

	maxY, maxX := stdscr.MaxYX()

	chatHeight, chatWidth, chatY, chatX := 5, maxX, maxY-5, 0
	win_chat, err := gc.NewWindow(chatHeight, chatWidth, chatY, chatX)
	if err != nil {
		log.Fatal("NewWindow:", err)
	}
	defer win_chat.Delete()
	err = win_chat.Border('|', '|', '-', '-', '+', '+', '+', '+')
	if err != nil {
		log.Fatal("border: ", err)
	}
	win_chat.MovePrint(1, 1, "Chat")
	win_chat.ScrollOk(true)
	win_chat.Refresh()

	height := maxY - 5
	win_input, err := gc.NewWindow(height, maxX, 0, 0)
	if err != nil {
		log.Fatal("NewWindow:", err)
	}
	defer win_input.Delete()
	win_input.Border('|', '|', '-', '-', '+', '+', '+', '+')
	win_input.MovePrint(1, 1, "Input")
	win_input.Refresh()

	curY, _ := win_chat.CursorYX()

	for {
		line := getLine(win_input)

		if line == "/exit" {
			return
		}

		for i := range chatHeight {
			win_chat.MovePrintln(i, 2, line)
		}

		win_input.Move(0, 0)
		win_input.MovePrint(3, 1, "> ")
		win_input.Clear()
		win_input.Move(curY+1, 3)
		win_input.Refresh()
	}
}

func getLine(win *gc.Window) string {
	line := ""

	for {
		key := win.GetChar()

		if key == 0 {
			// timeout
			continue
		}

		if key == '\n' {
			return line
		}

		line += gc.KeyString(key)
	}
}
