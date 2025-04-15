package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// Convos represents a conversation with a title and chat history
type Convos struct {
	Title       string                         `json:"title"`
	ChatHistory []openai.ChatCompletionMessage `json:"chat_history"`
	Provider    string                         `json:"provider"`
	Model       string                         `json:"model"`
	CreatedAt   time.Time                      `json:"created_at"`
	UpdatedAt   time.Time                      `json:"updated_at"`
}

// NewConvos creates a new conversation with the given title, provider, and model
func NewConvos(title, provider, model string) *Convos {
	now := time.Now()
	return &Convos{
		Title:       title,
		ChatHistory: []openai.ChatCompletionMessage{},
		Provider:    provider,
		Model:       model,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddMessage adds a message to the conversation history
func (c *Convos) AddMessage(role, content string) {
	c.ChatHistory = append(c.ChatHistory, openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})
	c.UpdatedAt = time.Now()
}

// GetChatHistoryDir returns the directory path for storing chat history
func GetChatHistoryDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	chatHistoryDir := filepath.Join(homeDir, ".config", "atlas", "chat-history")
	return chatHistoryDir, nil
}

// GetProviderModelDir returns the directory path for a specific provider and model
func GetProviderModelDir(provider, model string) (string, error) {
	chatHistoryDir, err := GetChatHistoryDir()
	if err != nil {
		return "", err
	}

	providerModelDir := filepath.Join(chatHistoryDir, fmt.Sprintf("%s-%s", provider, model))
	return providerModelDir, nil
}

// Save saves the conversation to a file
func (c *Convos) Save() error {
	// Update the timestamp
	c.UpdatedAt = time.Now()

	// Get the directory path
	providerModelDir, err := GetProviderModelDir(c.Provider, c.Model)
	if err != nil {
		return err
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(providerModelDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create a filename based on the creation time and title
	timestamp := c.CreatedAt.Format("2006-01-02-150405")
	safeTitle := sanitizeFilename(c.Title)
	if safeTitle == "" {
		safeTitle = "untitled"
	}

	filename := fmt.Sprintf("%s-%s.json", timestamp, safeTitle)
	filePath := filepath.Join(providerModelDir, filename)

	// Marshal the conversation to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %w", err)
	}

	// Write the file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// LoadConvos loads a conversation from a file
func LoadConvos(filePath string) (*Convos, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal the JSON
	var convo Convos
	if err := json.Unmarshal(data, &convo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversation: %w", err)
	}

	return &convo, nil
}

// ListConversations returns a list of all saved conversations for a provider and model
func ListConversations(provider, model string) ([]*Convos, error) {
	providerModelDir, err := GetProviderModelDir(provider, model)
	if err != nil {
		return nil, err
	}

	// Check if the directory exists
	if _, err := os.Stat(providerModelDir); os.IsNotExist(err) {
		// Directory doesn't exist, return an empty list
		return []*Convos{}, nil
	}

	// Read all JSON files in the directory
	files, err := filepath.Glob(filepath.Join(providerModelDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	// Load each conversation
	conversations := make([]*Convos, 0, len(files))
	for _, file := range files {
		convo, err := LoadConvos(file)
		if err != nil {
			// Log the error but continue with other files
			fmt.Printf("Error loading conversation from %s: %v\n", file, err)
			continue
		}
		conversations = append(conversations, convo)
	}

	return conversations, nil
}

// sanitizeFilename removes characters that are not allowed in filenames
func sanitizeFilename(filename string) string {
	// Replace invalid characters with underscores
	invalidChars := []rune{'/', '\\', ':', '*', '?', '"', '<', '>', '|'}
	result := filename

	for range invalidChars {
		result = filepath.Clean(result)
		result = filepath.Base(result)
	}

	// Limit the length of the filename
	if len(result) > 50 {
		result = result[:50]
	}

	return result
}
