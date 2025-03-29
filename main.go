package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// var config Config

var (
	viewArr = []string{"providers", "models", "conversations", "chatLog", "input"}
	active  = 0
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

// func handleEsc(g *gocui.Gui, v *gocui.View) error {
// 	b := make([]byte, 2)
// 	os.Stdin.Read(b)
// 	if string(b) == "[Z" {
// 		return prevView(g, v)
// 	}
// 	return nil
// }
//
// func prevView(g *gocui.Gui, v *gocui.View) error {
// 	nextIndex := active - 1
// 	if nextIndex < 0 {
// 		nextIndex = len(viewArr) - 1
// 	}
// 	name := viewArr[nextIndex]
//
// 	if _, err := setCurrentViewOnTop(g, name); err != nil {
// 		return err
// 	}
//
// 	if nextIndex == 4 {
// 		g.Cursor = true
// 	} else {
// 		g.Cursor = false
// 	}
//
// 	active = nextIndex
// 	return nil
// }

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 4 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

func keybindings(g *gocui.Gui) error {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})
	if err != nil {
		return err
	}

	// err = g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, handleEsc)
	// if err != nil {
	// 	return err
	// }

	err = g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '1', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := g.SetViewOnTop("conversations")
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Top-left "Models" view (fixed height).
	if v, err := g.SetView("providers", 0, 0, maxX/8-1, 10); err != nil {
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

		v.Title = "Providers"
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
		if _, err = setCurrentViewOnTop(g, "providers"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("models", maxX/8, 0, maxX/4-1, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Models"
		v.Clear()
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

	// Input box for Chat Log view
	if v, err := g.SetView("input", maxX/4, maxY-9, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Input"
		v.Editable = true
		v.Wrap = true
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

// func quit(g *gocui.Gui, v *gocui.View) error {
// 	return gocui.ErrQuit
// }

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.InputEsc = true

	g.SetManagerFunc(layout)

	// Keybinding for quitting: Ctrl+C
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

////////////////////////////////////////////////
///////////////////////////////////////////////
//// EXAMPLES TESTING AREA BELOW           ///
/////////////////////////////////////////////
////////////////////////////////////////////

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// package main
//
// import (
// 	"fmt"
// 	"log"
//
// 	"github.com/jroimartin/gocui"
// )
//
// var (
// 	viewArr = []string{"v1", "v2", "v3", "v4"}
// 	active  = 0
// )
//
// func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
// 	if _, err := g.SetCurrentView(name); err != nil {
// 		return nil, err
// 	}
// 	return g.SetViewOnTop(name)
// }
//
// func nextView(g *gocui.Gui, v *gocui.View) error {
// 	nextIndex := (active + 1) % len(viewArr)
// 	name := viewArr[nextIndex]
//
// 	out, err := g.View("v2")
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Fprintln(out, "Going from view "+v.Name()+" to "+name)
//
// 	if _, err := setCurrentViewOnTop(g, name); err != nil {
// 		return err
// 	}
//
// 	if nextIndex == 0 || nextIndex == 3 {
// 		g.Cursor = true
// 	} else {
// 		g.Cursor = false
// 	}
//
// 	active = nextIndex
// 	return nil
// }
//
// func layout(g *gocui.Gui) error {
// 	maxX, maxY := g.Size()
// 	if v, err := g.SetView("v1", 0, 0, maxX/2-1, maxY/2-1); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 		v.Title = "v1 (editable)"
// 		v.Editable = true
// 		v.Wrap = true
//
// 		if _, err = setCurrentViewOnTop(g, "v1"); err != nil {
// 			return err
// 		}
// 	}
//
// 	if v, err := g.SetView("v2", maxX/2-1, 0, maxX-1, maxY/2-1); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 		v.Title = "v2"
// 		v.Wrap = true
// 		v.Autoscroll = true
// 	}
// 	if v, err := g.SetView("v3", 0, maxY/2-1, maxX/2-1, maxY-1); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 		v.Title = "v3"
// 		v.Wrap = true
// 		v.Autoscroll = true
// 		fmt.Fprint(v, "Press TAB to change current view")
// 	}
// 	if v, err := g.SetView("v4", maxX/2, maxY/2, maxX-1, maxY-1); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 		v.Title = "v4 (editable)"
// 		v.Editable = true
// 	}
// 	return nil
// }
//
// func quit(g *gocui.Gui, v *gocui.View) error {
// 	return gocui.ErrQuit
// }
//
// func main() {
// 	g, err := gocui.NewGui(gocui.OutputNormal)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	defer g.Close()
//
// 	g.Highlight = true
// 	g.Cursor = true
// 	g.SelFgColor = gocui.ColorGreen
//
// 	g.SetManagerFunc(layout)
//
// 	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
// 		log.Panicln(err)
// 	}
// 	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
// 		log.Panicln(err)
// 	}
//
// 	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
// 		log.Panicln(err)
// 	}
// }
