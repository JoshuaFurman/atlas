package main

import "github.com/jroimartin/gocui"

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Top-left "Providers" view (fixed height).
	if v, err := g.SetView("providers", 0, 0, maxX/8-1, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "[1]-Providers"
		v.SelBgColor = gocui.ColorBlack
		v.SelFgColor = gocui.ColorGreen
		v.Wrap = true

		// Initialize the providers view with the list of providers
		updateProvidersView(g)
	}

	if v, err := g.SetView("models", maxX/8, 0, maxX/4-1, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "[2]-Models"
		v.SelBgColor = gocui.ColorBlack
		v.SelFgColor = gocui.ColorGreen
		v.Wrap = true

		updateModelsView(g)
	}

	// Left column "Conversations" view (below "Models").
	if v, err := g.SetView("conversations", 0, 11, maxX/4-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "[3]-Conversations"

		updateConvosView(g)
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
	}
	return nil
}
