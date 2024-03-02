package main

import (
	"io"
	"log"
	"strings"

	"johnnyasantos.com/chat/shared"

	gc "github.com/rthornton128/goncurses"
)

func main() {
	log.Default().SetOutput(io.Discard)

	stdscr, err := gc.Init()
	if err != nil {
		log.Fatal("Init window: ", err)
	}
	defer gc.End()
	defer stdscr.Delete()

	stdscr.Keypad(true)
	gc.Raw(false)
	stdscr.Timeout(0)

	maxY, maxX := stdscr.MaxYX()

	chatHeight, chatWidth, chatY, chatX := 5, maxX, maxY-5, 0
	winChat, err := gc.NewWindow(chatHeight, chatWidth, chatY, chatX)
	if err != nil {
		log.Fatal("NewWindow:", err)
	}
	defer winChat.Delete()
	err = winChat.Border('|', '|', '-', '-', '+', '+', '+', '+')
	if err != nil {
		log.Fatal("border: ", err)
	}
	winChat.MovePrint(1, 1, "Chat")
	winChat.ScrollOk(true)
	winChat.Refresh()

	height := maxY - 5
	winInput, err := gc.NewWindow(height, maxX, 0, 0)
	if err != nil {
		log.Fatal("NewWindow:", err)
	}
	defer winInput.Delete()
	winInput.Border('|', '|', '-', '-', '+', '+', '+', '+')
	winInput.MovePrint(1, 1, "Input")
	winInput.Refresh()

	linesChan := make(chan string, 10)
	go readLine(winInput, linesChan)

	closing := make(chan bool, 1)
	shared.HandleSignals(closing)

	line := ""

	for {
		select {
		case line = <-linesChan:
		case <-closing:
			return
		}

		if strings.TrimSpace(line) == "" {
			continue
		}

		if line == "/exit" {
			return
		}

		winChat.MovePrintln(chatHeight-5, 2, line)
		winChat.Refresh()

		winInput.MovePrint(maxY-5, 1, "> ")
		winInput.Refresh()
	}
}

func readLine(win *gc.Window, lines chan<- string) {
	line := ""

	for {
		key := win.GetChar()

		if key == 0 {
			// timeout
			continue
		}

		if key == '\n' {
			lines <- line
			line = ""
			continue
		}

		line += gc.KeyString(key)
	}
}
