package main

import "github.com/jroimartin/gocui"

func keybindings(g *gocui.Gui) error {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})
	if err != nil {
		return err
	}

	// Arrow keys for navigating the models list when models view is active
	err = g.SetKeybinding("models", gocui.KeyArrowUp, gocui.ModNone, moveModelUp)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("models", gocui.KeyArrowDown, gocui.ModNone, moveModelDown)
	if err != nil {
		return err
	}

	// Add vim-style navigation with 'j' and 'k' keys for models
	err = g.SetKeybinding("models", 'k', gocui.ModNone, moveModelUp)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("models", 'j', gocui.ModNone, moveModelDown)
	if err != nil {
		return err
	}

	// Add Enter key binding to select the current model
	err = g.SetKeybinding("models", gocui.KeyEnter, gocui.ModNone, selectModel)
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

	// Arrow keys for navigating the convos list when convos view is active
	err = g.SetKeybinding("convos", gocui.KeyArrowUp, gocui.ModNone, moveConvoUp)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("convos", gocui.KeyArrowDown, gocui.ModNone, moveConvoDown)
	if err != nil {
		return err
	}

	// Add vim-style navigation with 'j' and 'k' keys for convos
	err = g.SetKeybinding("convos", 'k', gocui.ModNone, moveConvoUp)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("convos", 'j', gocui.ModNone, moveConvoDown)
	if err != nil {
		return err
	}

	// Add Enter key binding to select the current convo
	err = g.SetKeybinding("convos", gocui.KeyEnter, gocui.ModNone, selectConvo)
	if err != nil {
		return err
	}

	// Add Enter key binding to process input text
	err = g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, processInput)
	if err != nil {
		return err
	}

	// Add Shift+Enter key binding to insert a new line
	err = g.SetKeybinding("input", gocui.KeyEnter, gocui.ModAlt, insertNewLine)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '1', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := setCurrentViewOnTop(g, "providers")
		g.Cursor = false
		active = 0
		return err
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '2', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := setCurrentViewOnTop(g, "models")
		g.Cursor = false
		active = 1
		return err
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '3', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := setCurrentViewOnTop(g, "conversations")
		g.Cursor = false
		active = 2
		return err
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '4', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := setCurrentViewOnTop(g, "chatLog")
		g.Cursor = false
		active = 3
		return err
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '5', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := setCurrentViewOnTop(g, "input")
		g.Cursor = true
		active = 4
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
