package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// var config Config

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Top-left "Models" view (fixed height).
	if v, err := g.SetView("models", 0, 0, maxX/4-1, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Load config from default path
		config, err := LoadConfig()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		config.ActiveProvider = "openai" // will eventually pull from some file that save active config
		config.ActiveModel = "gpt-4o"    // will eventually pull from some file that save active config

		v.Title = "Active Model"
		v.Clear()

		// providers := config.GetAllProviders()
		// models, err := config.GetModelsForProvider(providers[0])
		// if err != nil {
		// 	log.Printf("Error getting models for %s: %v", providers[0], err)
		// }
		fmt.Fprintf(v, "%s -> %s", config.ActiveProvider, config.ActiveModel)
		// for _, provider := range providers {
		// 	fmt.Fprintln(v, provider)
		// 	for _, model := range models {
		// 		if model.Name == config.ActiveModel {
		// 			fmt.Fprintf(v, "  * %s\n", model.Name)
		// 		} else {
		// 			fmt.Fprintf(v, "  %s\n", model.Name)
		// 		}
		// 	}
		// }

	}

	// Left column "Conversations" view (below "Models").
	if v, err := g.SetView("conversations", 0, 11, maxX/4-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Conversations"
	}

	// Right-side "Chat Log" view.
	if v, err := g.SetView("chatLog", maxX/4, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Chat Log"
	}

	// Bottom "commandBar" view (spanning the full width).
	if v, err := g.SetView("commandBar", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Command"
		// v.Frame = false
		// v.Autoscroll = true
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	// Keybinding for quitting: Ctrl+C
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
