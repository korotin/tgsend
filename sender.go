package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	sendMessageUrlMask = "https://api.telegram.org/bot%s/sendMessage"
)

type (
	endpointResponse struct {
		Ok bool
	}
)

func getParseMode(format messageFormat) string {
	switch format {
	case html:
		return "html"
	case markdown:
		return "markdown"
	}

	return ""
}

func getUrlValues(config *config, input *userInput) url.Values {
	values := url.Values{}
	values.Set("chat_id", string(config.GetChatId(input.chat)))
	values.Set("text", *input.message)
	if input.silent {
		values.Set("disable_notification", "1")
	}
	if input.format != raw {
		values.Set("parse_mode", getParseMode(input.format))
	}

	return values
}

func getUrl(config *config, input *userInput) string {
	return fmt.Sprintf(sendMessageUrlMask, config.GetBotId(input.bot))
}

func sendMessage(config *config, input *userInput) error {
	endpointUrl := getUrl(config, input)
	endpointParams := getUrlValues(config, input)
	response, err := http.PostForm(endpointUrl, endpointParams)
	if err != nil {
		return err
	}

	rawBody, _ := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()

	endpointResponse := endpointResponse{false}
	if err := json.Unmarshal(rawBody, &endpointResponse); err != nil {
		return err
	}

	if !endpointResponse.Ok {
		return fmt.Errorf("failed to send message: %s", string(rawBody))
	}

	return nil
}
