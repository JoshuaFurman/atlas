package main

import (
	"log"

	"github.com/jroimartin/gocui"
	openai "github.com/sashabaranov/go-openai"
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

// Navigate up in the convos list
func moveConvoUp(g *gocui.Gui, v *gocui.View) error {
	if selectedConvo > 0 {
		selectedConvo--
		updateConvosView(g)
	}
	return nil
}

// Navigate down in the conversations list
func moveConvoDown(g *gocui.Gui, v *gocui.View) error {
	if selectedConvo < len(conversations)-1 {
		selectedConvo++
		updateConvosView(g)
	}
	return nil
}

// Select the currently highlighted model as the active model
func selectModel(g *gocui.Gui, v *gocui.View) error {
	activeModel = selectedModel
	updateModelsView(g)
	config.ActiveModel = models[activeModel].Name

	// Save the current conversation first
	if currentConvo != nil && len(currentConvo.ChatHistory) > 1 {
		if err := currentConvo.Save(); err != nil {
			log.Printf("Failed to save current conversation: %v", err)
		}
	}

	// Create a new conversation with the selected model
	currentConvo = NewConvos("New Chat", providers[activeProvider], models[activeModel].Name)
	currentConvo.AddMessage(openai.ChatMessageRoleSystem, models[activeModel].SystemPrompt)

	// Load conversations for the new provider/model combination
	loadConversations(providers[activeProvider], models[activeModel].Name)
	updateConvosView(g)

	chatLogView, err := g.View("chatLog")
	if err != nil {
		return err
	}
	chatLogView.Clear()
	chatLogView.SetCursor(0, 0)

	inputView, err := g.View("input")
	if err != nil {
		return err
	}
	inputView.Clear()
	inputView.SetCursor(0, 0)

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
	config.ActiveProvider = providers[activeProvider]

	// Save the current conversation first
	if currentConvo != nil && len(currentConvo.ChatHistory) > 1 {
		if err := currentConvo.Save(); err != nil {
			log.Printf("Failed to save current conversation: %v", err)
		}
	}

	// Create a new conversation with the selected provider and model
	currentConvo = NewConvos("New Chat", providers[activeProvider], models[activeModel].Name)
	currentConvo.AddMessage(openai.ChatMessageRoleSystem, models[activeModel].SystemPrompt)

	// Load conversations for the new provider/model combination
	loadConversations(providers[activeProvider], models[activeModel].Name)
	updateConvosView(g)

	chatLogView, err := g.View("chatLog")
	if err != nil {
		return err
	}
	chatLogView.Clear()
	chatLogView.SetCursor(0, 0)

	inputView, err := g.View("input")
	if err != nil {
		return err
	}
	inputView.Clear()
	inputView.SetCursor(0, 0)

	return nil
}

// Select the currently highlighted conversation as the active conversation
func selectConvo(g *gocui.Gui, v *gocui.View) error {
	if selectedConvo < 0 || selectedConvo >= len(conversations) {
		return nil
	}

	// Save the current conversation first
	if currentConvo != nil && len(currentConvo.ChatHistory) > 1 {
		if err := currentConvo.Save(); err != nil {
			log.Printf("Failed to save current conversation: %v", err)
		}
	}

	// Load the selected conversation
	loadConversation(selectedConvo)
	updateConvosView(g)

	// Update the chat log view with the conversation history
	chatLogView, err := g.View("chatLog")
	if err != nil {
		return err
	}
	chatLogView.Clear()

	// Display all messages except the system prompt
	for i, msg := range currentConvo.ChatHistory {
		if i == 0 && msg.Role == openai.ChatMessageRoleSystem {
			// Skip system prompt
			continue
		}

		switch msg.Role {
		case openai.ChatMessageRoleUser:
			addUserMessage(chatLogView, msg.Content)
		case openai.ChatMessageRoleAssistant:
			addAIResponse(chatLogView, msg.Content)
		}
	}

	return nil
}
