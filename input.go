package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	raw messageFormat = iota
	markdown
	html
)

type (
	messageFormat int

	userInput struct {
		bot     aliasName
		chat    aliasName
		format  messageFormat
		silent  bool
		message *string
	}
)

func (i *userInput) Init(config *config) {
	i.bot = config.defaultBot
	i.chat = config.defaultChat
	i.format = raw
	i.silent = false
	i.message = new(string)
}

func readStdin() *string {
	bytesData, _ := ioutil.ReadAll(os.Stdin)
	stringData := string(bytesData)

	return &stringData
}

func readInput(config *config) (*userInput, error) {
	input := userInput{}
	input.Init(config)

	var isMarkdown, isHtml bool

	flag.Var(&input.bot, "bot", "bot alias")
	flag.Var(&input.chat, "chat", "chat alias")
	flag.BoolVar(&input.silent, "silent", false, "send message in silent mode")
	flag.StringVar(input.message, "msg", "", "message to send")
	flag.BoolVar(&isMarkdown, "md", false, "use markdown for formatting")
	flag.BoolVar(&isHtml, "html", false, "use html for formatting")

	flag.Parse()

	if _, ok := config.bots[input.bot]; !ok {
		return nil, fmt.Errorf("unknown bot alias: %s", input.bot)
	}

	if _, ok := config.chats[input.chat]; !ok {
		return nil, fmt.Errorf("unknown chat alias: %s", input.bot)
	}

	if isMarkdown && isHtml {
		return nil, errors.New("md and html flags are mutually exclusive")
	}

	stdinMessage := readStdin()
	if *input.message == "" && *stdinMessage == "" {
		return nil, errors.New("message is empty")
	}
	if *input.message != "" && *stdinMessage != "" {
		return nil, errors.New("message should be passed via stdin or via argument but not both")
	}
	if *stdinMessage != "" {
		input.message = stdinMessage
	}

	switch true {
	case isMarkdown:
		input.format = markdown
	case isHtml:
		input.format = html
	}

	return &input, nil
}
