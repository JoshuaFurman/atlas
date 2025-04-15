package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	openai "github.com/sashabaranov/go-openai"
)

type Chat struct {
	context context.Context
	client  *openai.Client
	request openai.ChatCompletionRequest
}

func formatCode(code, lang string) {
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("solarized-dark")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	err = formatter.Format(os.Stdout, style, iterator)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func processChar(char rune, state *int, buffer *string, code *string, lang *string) {
	switch *state {
	case 0:
		if char == '`' {
			*state = 1
		} else {
			fmt.Print(string(char))
		}
	case 1:
		if char == '`' {
			*state = 2
		} else {
			*state = 0
			fmt.Print("`", string(char))
		}
	case 2:
		if char == '`' {
			*state = 3
			*buffer = ""
		} else {
			*state = 0
			fmt.Print("``", string(char))
		}
	case 3:
		if char == '\n' {
			*state = 4
			*lang = *buffer
			*buffer = ""
		} else {
			*buffer += string(char)
		}
	case 4:
		if char == '`' {
			*state = 5
		} else {
			*code += string(char)
		}
	case 5:
		if char == '`' {
			*state = 6
		} else {
			*state = 4
			*code += "`" + string(char)
		}
	case 6:
		if char == '`' {
			formatCode(*code, *lang)
			*code = ""
			*lang = ""
			*state = 0
		} else {
			*state = 4
			*code += "``" + string(char)
		}
	}
}
