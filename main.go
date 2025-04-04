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

	// List of available providers
	providers        []string
	selectedProvider = 0 // Index of the currently selected provider
	activeProvider   = 0 // Index of the active (confirmed) provider
	config           *Config
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

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

// Navigate up in the providers list
func moveProviderUp(g *gocui.Gui, v *gocui.View) error {
	if selectedProvider > 0 {
		selectedProvider--
		updateProvidersView(g)
	}
	return nil
}

// Navigate down in the providers list
func moveProviderDown(g *gocui.Gui, v *gocui.View) error {
	if selectedProvider < len(providers)-1 {
		selectedProvider++
		updateProvidersView(g)
	}
	return nil
}

// Select the currently highlighted provider as the active provider
func selectProvider(g *gocui.Gui, v *gocui.View) error {
	activeProvider = selectedProvider
	updateProvidersView(g)

	// Here you would typically update your configuration or application state
	// to use the newly selected provider

	return nil
}

// Update the providers view with the current selection
func updateProvidersView(g *gocui.Gui) error {
	v, err := g.View("providers")
	if err != nil {
		return err
	}

	v.Clear()

	// Display the list of providers with the selected one in green
	// and the active one with an asterisk
	for i, provider := range providers {
		prefix := "  " // Default prefix (two spaces)
		if i == activeProvider {
			prefix = "* " // Asterisk for active provider
		}

		if i == selectedProvider {
			fmt.Fprintf(v, "\033[32m%s%s\033[0m\n", prefix, provider) // Green color for selected
		} else {
			fmt.Fprintf(v, "%s%s\n", prefix, provider)
		}
	}

	// Update the active provider in the config
	// This would normally update your config, but for now we'll just display it
	fmt.Fprintf(v, "\n%s -> %s", providers[activeProvider], "gpt-4o")

	return nil
}

func keybindings(g *gocui.Gui) error {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})
	if err != nil {
		return err
	}

	// Arrow keys for navigating the providers list when providers view is active
	err = g.SetKeybinding("providers", gocui.KeyArrowUp, gocui.ModNone, moveProviderUp)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("providers", gocui.KeyArrowDown, gocui.ModNone, moveProviderDown)
	if err != nil {
		return err
	}

	// Add vim-style navigation with 'j' and 'k' keys
	err = g.SetKeybinding("providers", 'k', gocui.ModNone, moveProviderUp)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("providers", 'j', gocui.ModNone, moveProviderDown)
	if err != nil {
		return err
	}

	// Add Enter key binding to select the current provider
	err = g.SetKeybinding("providers", gocui.KeyEnter, gocui.ModNone, selectProvider)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '1', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := setCurrentViewOnTop(g, "providers")
		g.Cursor = false
		return err
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '2', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := setCurrentViewOnTop(g, "models")
		g.Cursor = false
		return err
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '3', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := setCurrentViewOnTop(g, "conversations")
		g.Cursor = false
		return err
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '4', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := setCurrentViewOnTop(g, "chatLog")
		g.Cursor = false
		return err
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '5', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := setCurrentViewOnTop(g, "input")
		g.Cursor = true
		return err
	})
	if err != nil {
		return err
	}

	// Tab to move to next view
	err = g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView)
	if err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Top-left "Providers" view (fixed height).
	if v, err := g.SetView("providers", 0, 0, maxX/8-1, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Load config from default path
		// var err error
		config, err = LoadConfig()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		// Get providers from config
		providers = config.GetAllProviders()

		v.Title = "[1]-Providers"
		v.SelBgColor = gocui.ColorBlack
		v.SelFgColor = gocui.ColorGreen

		// Initialize the providers view with the list of providers
		updateProvidersView(g)
	}

	if v, err := g.SetView("models", maxX/8, 0, maxX/4-1, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "[2]-Models"
		v.Clear()
	}

	// Left column "Conversations" view (below "Models").
	if v, err := g.SetView("conversations", 0, 11, maxX/4-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "[3]-Conversations"
	}

	// Right-side "Chat Log" view.
	if v, err := g.SetView("chatLog", maxX/4, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "[4]-Chat Log"
	}

	// Input box for Chat Log view
	if v, err := g.SetView("input", maxX/4, maxY-9, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "[5]-Input"
		v.Editable = true
		v.Wrap = true

		if _, err = setCurrentViewOnTop(g, "input"); err != nil {
			return err
		}

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
