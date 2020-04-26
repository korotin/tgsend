package main

import (
	"log"
)

func main() {
	log.SetPrefix("tgsend: ")

	var (
		config *config
		input  *userInput
		err    error
	)

	config, err = getConfig()
	if err != nil {
		log.Fatal(err)
	}

	input, err = readInput(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s → %s...", input.bot, input.chat)

	err = sendMessage(config, input)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Message sent")
}
