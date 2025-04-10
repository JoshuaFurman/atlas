package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/sashabaranov/go-openai"
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
	currentChat      []openai.ChatCompletionMessage
)

// Process the input text when Enter is pressed
func processInput(g *gocui.Gui, v *gocui.View) error {
	inputText := v.Buffer()

	// Trim the input text to remove trailing newline
	// inputText = strings.TrimSpace(inputText)

	// Skip empty messages
	if inputText == "" {
		return nil
	}

	// Add the user message to the chat log (right-aligned)
	chatLogView, err := g.View("chatLog")
	if err != nil {
		return err
	}

	// Format user message (right-aligned with padding)
	addUserMessage(chatLogView, inputText)

	// Clear the input view after processing
	v.Clear()
	v.SetCursor(0, 0)

	currentProvider, err := config.GetProviderConfig(providers[activeProvider])
	if err != nil {
		log.Fatalf("couldnt get api-key: %v\n", err)
	}

	var client *openai.Client
	if providers[activeProvider] != "openai" {
		clientConfig := openai.DefaultConfig(currentProvider.APIKey)
		clientConfig.BaseURL = currentProvider.Endpoint
		client = openai.NewClientWithConfig(clientConfig)
	} else {
		client = openai.NewClient(currentProvider.APIKey)
	}
	context := context.Background()
	currentChat = append(currentChat, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: inputText,
	})

	go func() {
		request := openai.ChatCompletionRequest{
			Model:       models[activeModel].Name,
			Temperature: models[activeModel].Temperature,
			Messages:    currentChat,
		}
		response, err := client.CreateChatCompletion(context, request)
		if err != nil {
			log.Fatalf("ChatCompletionStream error: %v\n", err)
			return
		}

		g.Update(func(g *gocui.Gui) error {
			addAIResponse(chatLogView, response.Choices[0].Message.Content)
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
		padding := max(width-len(line)-2, 0)
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

	// add AI repsonse back to the chat history
	currentChat = append(currentChat, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: message,
	})
}

// Format a message with word wrapping
func formatMessage(message string, maxWidth int, isUser bool) string {
	words := strings.Fields(message)
	if len(words) == 0 {
		return ""
	}

	// Add a prefix to indicate who's speaking
	prefix := models[activeModel].Name + ": "
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

// Insert a new line in the input view when alt+Enter is pressed
func insertNewLine(g *gocui.Gui, v *gocui.View) error {
	v.EditNewLine()
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

	// ensures that openai is always the first selected provider
	for i, provider := range providers {
		if provider == "openai" {
			activeProvider = i
		}
	}
	config.ActiveProvider = providers[activeProvider]
	config.ActiveModel = models[activeModel].Name
	active = 4 // sets active pane to input view

	// setup initial currentChat
	currentChat = append(currentChat, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: models[activeModel].SystemPrompt,
	},
	)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.Mouse = true
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
