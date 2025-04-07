package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

var (
	viewArr = []string{"providers", "models", "conversations", "chatLog", "input"}
	active  = 0

	// List of available providers
	providers        []string
	models           []ModelConfig
	selectedProvider = 0 // Index of the currently selected provider
	activeProvider   = 0 // Index of the active (confirmed) provider
	selectedModel    = 0 // Index of the currently selected model
	activeModel      = 0 // Index of the active (confirmed) model
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

// Process the input text when Enter is pressed
func processInput(g *gocui.Gui, v *gocui.View) error {
	inputText := v.Buffer()

	// Trim the input text to remove trailing newline
	inputText = strings.TrimSpace(inputText)

	// Skip empty messages
	if inputText == "" {
		return nil
	}

	// Add the user message to the chat log (right-aligned)
	chatLogView, err := g.View("chatLog")
	if err != nil {
		return err
	}

	// Get the width of the chatLog view
	// _, maxY := chatLogView.Size()

	// Format user message (right-aligned with padding)
	addUserMessage(chatLogView, inputText)

	// Clear the input view after processing
	v.Clear()
	v.SetCursor(0, 0)

	// Simulate an AI response after a short delay
	go func() {
		time.Sleep(500 * time.Millisecond) // Simulate thinking time
		g.Update(func(g *gocui.Gui) error {
			// Generate a mock response
			response := generateResponse(inputText)

			// Add the AI response to the chat log (left-aligned)
			addAIResponse(chatLogView, response)
			return nil
		})
	}()

	return nil
}

// Add a user message to the chat log (right-aligned)
func addUserMessage(v *gocui.View, message string) {
	width, _ := v.Size()

	// Format the message with word wrapping
	formattedMsg := formatMessage(message, width-10, true) // -10 for padding

	// Add a separator line
	fmt.Fprintln(v)

	// Print the formatted message lines (right-aligned)
	for _, line := range strings.Split(formattedMsg, "\n") {
		padding := width - len(line) - 2
		if padding < 0 {
			padding = 0
		}
		fmt.Fprintf(v, "%s\033[32m%s\033[0m\n", strings.Repeat(" ", padding), line)
	}

	// Auto-scroll to the bottom
	v.Autoscroll = true
}

// Add an AI response to the chat log (left-aligned)
func addAIResponse(v *gocui.View, message string) {
	width, _ := v.Size()

	// Format the message with word wrapping
	formattedMsg := formatMessage(message, width-10, false) // -10 for padding

	// Add a separator line
	fmt.Fprintln(v)

	// Print the formatted message lines (left-aligned with padding)
	for _, line := range strings.Split(formattedMsg, "\n") {
		fmt.Fprintf(v, "  \033[36m%s\033[0m\n", line)
	}

	// Auto-scroll to the bottom
	v.Autoscroll = true
}

// Format a message with word wrapping
func formatMessage(message string, maxWidth int, isUser bool) string {
	words := strings.Fields(message)
	if len(words) == 0 {
		return ""
	}

	// Add a prefix to indicate who's speaking
	prefix := "AI: "
	if isUser {
		prefix = "You: "
	}

	var lines []string
	currentLine := prefix

	for _, word := range words {
		// Check if adding this word would exceed the max width
		if len(currentLine)+len(word)+1 > maxWidth {
			// Start a new line
			lines = append(lines, currentLine)
			currentLine = "    " + word // Indent continuation lines
		} else {
			if currentLine == prefix {
				currentLine += word
			} else {
				currentLine += " " + word
			}
		}
	}

	// Add the last line
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// Generate a mock response based on the user input
func generateResponse(input string) string {
	// For now, just echo back a simple response
	responses := []string{
		"I understand you're saying: " + input,
		"That's an interesting point about: " + input,
		"Let me think about: " + input,
		"I'm processing your message: " + input,
		"Thanks for sharing your thoughts on: " + input,
	}

	return responses[0]
}

// Insert a new line in the input view when alt+Enter is pressed
func insertNewLine(g *gocui.Gui, v *gocui.View) error {
	v.EditNewLine()
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

// Navigate up in the models list
func moveModelUp(g *gocui.Gui, v *gocui.View) error {
	if selectedModel > 0 {
		selectedModel--
		updateModelsView(g)
	}
	return nil
}

// Navigate down in the models list
func moveModelDown(g *gocui.Gui, v *gocui.View) error {
	if selectedModel < len(models)-1 {
		selectedModel++
		updateModelsView(g)
	}
	return nil
}

// Select the currently highlighted model as the active model
func selectModel(g *gocui.Gui, v *gocui.View) error {
	activeModel = selectedModel
	updateModelsView(g)

	// Here you would typically update your configuration or application state
	// to use the newly selected model

	return nil
}

// Update the models view with the current selection
func updateModelsView(g *gocui.Gui) error {
	v, err := g.View("models")
	if err != nil {
		return err
	}

	v.Clear()

	// Display the list of models with the selected one in green
	// and the active one with an asterisk
	for i, model := range models {
		prefix := "  " // Default prefix (two spaces)
		if i == activeModel {
			prefix = "* " // Asterisk for active model
		}

		if i == selectedModel {
			fmt.Fprintf(v, "\033[32m%s%s\033[0m\n", prefix, model.Name) // Green color for selected
		} else {
			fmt.Fprintf(v, "%s%s\n", prefix, model.Name)
		}
	}

	// Update the active model in the config
	// This would normally update your config, but for now we'll just display it
	fmt.Fprintf(v, "\033[32m\n\nACTIVE:\n%s -> %s\033[0m", providers[activeProvider], models[activeModel].Name)

	return nil
}

// Select the currently highlighted provider as the active provider
func selectProvider(g *gocui.Gui, v *gocui.View) error {
	activeProvider = selectedProvider
	// Get models for the selected provider from config
	var err error
	models, err = config.GetModelsForProvider(providers[activeProvider])
	if err != nil {
		log.Fatalf("Failed to get models: %v", err)
	}
	selectedModel = 0
	activeModel = 0

	updateProvidersView(g)
	updateModelsView(g)

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
	fmt.Fprintf(v, "\033[32m\n\nACTIVE:\n%s -> %s\033[0m", providers[activeProvider], models[activeModel].Name)

	return nil
}

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

	///////////////////////// VIEWS BEGIN

	// Top-left "Providers" view (fixed height).
	if v, err := g.SetView("providers", 0, 0, maxX/8-1, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

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
		v.SelBgColor = gocui.ColorBlack
		v.SelFgColor = gocui.ColorGreen

		updateModelsView(g)
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
	}
	return nil
}

func main() {
	// Load config from default path
	var err error
	config, err = LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Get providers from config
	providers = config.GetAllProviders()

	// Get models for the selected provider from config
	models, err = config.GetModelsForProvider(providers[selectedProvider])
	if err != nil {
		log.Fatalf("Failed to get models: %v", err)
	}

	active = 4 // sets active pane to input view

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	// Keybinding for quitting: Ctrl+C
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
