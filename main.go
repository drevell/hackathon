// Copyright 2023 The Authors (see AUTHORS file)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/abcxyz/pkg/cli"
)

const (
	envVarWithJsonInput string = "GITHUB_CONTEXT"
	baseUrl             string = "https://chat.googleapis.com/v1/spaces/%s/messages?key=%s&token=%s"
)

type ChatCommand struct {
	cli.BaseCommand

	flagSpace string
	flagKey   string
	flagToken string
}

func (c *ChatCommand) Desc() string {
	return "Send a message to a Chat space"
}

// TODO: nest a "workflowfailed" command under this.
func (c *ChatCommand) Help() string {
	return `
Usage: {{ COMMAND }} [options]

  The chat command sends messages to Chat spaces.
`
}

func (c *ChatCommand) Flags() *cli.FlagSet {
	set := cli.NewFlagSet()

	f := set.NewSection("Chat space options")

	f.StringVar(&cli.StringVar{
		Name:    "space",
		Aliases: []string{"s"},
		Example: "AAAAUJgrNvE",
		Default: "",
		Target:  &c.flagSpace,
		Usage:   "Identifier for chat space to send message to.",
	})

	f.StringVar(&cli.StringVar{
		Name:    "key",
		Aliases: []string{"k"},
		Example: "BFzaSyDdI0hCZ...",
		Default: "",
		Target:  &c.flagKey,
		Usage:   `"key" from chat webhook url.`,
	})

	f.StringVar(&cli.StringVar{
		Name:    "token",
		Aliases: []string{"t"},
		Example: "VyzWCmy_DdI0hCZ...",
		Default: "",
		Target:  &c.flagToken,
		Usage:   `"token" from chat webhook url.`,
	})

	return set
}

func (c *ChatCommand) Run(ctx context.Context, args []string) error {
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	args = f.Args()
	if len(args) != 0 {
		return fmt.Errorf("expected 0 arguments, got %q", args)
	}

	ghJson := os.Getenv(envVarWithJsonInput)
	if ghJson == "" {
		fmt.Printf("warning: %s not set, will use demo values", envVarWithJsonInput)
	}
	fmt.Println("ghJson: ", ghJson)

	b, err := messageBody(ghJson)
	if err != nil {
		return fmt.Errorf("failed to generate message body: %w", err)
	}

	url := fmt.Sprintf(baseUrl, c.flagSpace, c.flagKey, c.flagToken)
	fmt.Println("url: ", url)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("creating http request failed: %w", err)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("sending http request failed: %w", err)
	}
	fmt.Println("resp: ", resp)
	defer resp.Body.Close()

	testUserID, err := userNameTOUserId("Rui Zhang")
	if err != nil {
		return fmt.Errorf("failed to get userID: %w", err)
	}
	fmt.Println("testUserID: ", testUserID)

	return nil
}

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func realMain() error {
	// Create the command.
	rootCmd := func() cli.Command {
		return &cli.RootCommand{
			Name:    "hackathon",
			Version: "0.1",
			Commands: map[string]cli.CommandFactory{
				"chat": func() cli.Command {
					return &ChatCommand{}
				},
			},
		}
	}

	cmd := rootCmd()

	// Help output is written to stderr by default. Redirect to stdout so the
	// "Output" assertion works.
	cmd.SetStderr(os.Stdout)

	ctx := context.Background()
	err := cmd.Run(ctx, os.Args[1:])
	if err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}

	/*
		ghJson := os.Getenv(envVarWithJsonInput)
		if ghJson == "" {
			fmt.Println("warning: ", envVarWithJsonInput, " not set, will use demo values")
		}
		fmt.Println("ghJson: ", ghJson)

		b, err := messageBody(ghJson)
		if err != nil {
			return fmt.Errorf("failed to generate message body: %w", err)
		}

		url := fmt.Sprintf(baseUrl, Space, Key, Token)
		fmt.Println("url: ", url)
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
		if err != nil {
			return fmt.Errorf("creating http request failed: %w", err)
		}
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return fmt.Errorf("sending http request failed: %w", err)
		}
		fmt.Println("resp: ", resp)
		defer resp.Body.Close()

		testUserID, err := userNameTOUserId("Rui Zhang")
		if err != nil {
			return fmt.Errorf("failed to get userID: %w", err)
		}
		fmt.Println("testUserID: ", testUserID)
	*/

	return nil
}

// ghJson: JSON blob from Github workflow
func messageBody(ghJson string) ([]byte, error) {
	// TODO: convert ghJson string -> bytes[], then pass to json.Unmarshal() to
	// get a JSON message (https://pkg.go.dev/encoding/json#example-Unmarshal),
	// then retrieve the values we want to put into the chat message.
	parsedGhJson := map[string]any{}
	err := json.Unmarshal([]byte(ghJson), &parsedGhJson)
	if err != nil {
		return nil, fmt.Errorf("JSON unmarshal error: %w:", err)
	}

	// example of the parserdGhJson
	// "ref": "refs/heads/main",
	//   "sha": "541bf7efdce0ec8d874c61e12a65560a51e7f6be",
	//   "repository": "drevell/hackathon",
	//   "repository_owner": "drevell",
	//   "run_id": "4866865265",
	//   "run_number": "3",
	//   "run_attempt": "1",
	//   "triggering_actor": "drevell",
	//   "workflow": "Revell testing",

	text, err := messageText()
	if err != nil {
		return nil, fmt.Errorf("failed to generate message text: %w", err)
	}
	//var jsonData = []byte(fmt.Sprintf(
	//	`{
	//		"text": %s
	//	}`, text))
	// Example from https://developers.google.com/chat/api/guides/crudl/messages#create_a_card_message
	// We can also add buttons and other widgets
	jsonData := map[string]any{
		"cardsV2": map[string]any{
			"cardId": "createCardMessage",
			"card": map[string]any{
				"header": map[string]any{
					"title":     "This is the title",
					"subtitle":  text,
					"imageUrl":  "https://developers.google.com/chat/images/chat-product-icon.png",
					"imageType": "CIRCLE",
				},
				"sections": []any{
					map[string]any{
						"header":                    "This is the section header",
						"collapsible":               true,
						"uncollapsibleWidgetsCount": 1,
						"widgets": []map[string]any{
							{
								"decoratedText": map[string]any{
									"startIcon": map[string]any{
										"knownIcon": "DESCRIPTION",
									},
									"text": fmt.Sprintf("GitHub workflow failed: <pre>%s</pre>", parsedGhJson["workflow"]),
								},
							},
							{
								"decoratedText": map[string]any{
									"startIcon": map[string]any{
										"knownIcon": "DESCRIPTION",
									},
									"text": fmt.Sprintf("Run by: <pre>%s</pre>", parsedGhJson["triggering_actor"]),
								},
							},
							{
								"buttonList": map[string]any{
									"buttons": []any{
										map[string]any{
											"text": "Open Failing Run",
											"onClick": map[string]any{
												"openLink": map[string]any{
													"url": fmt.Sprintf("https://github.com/%s/actions/runs/%s",
														parsedGhJson["repository"], parsedGhJson["run_id"]),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return json.Marshal(jsonData)
}

func messageText() (string, error) {
	//t := "Hey <users/100449440289517201826>"
	t := "Hey from cli"
	return t, nil
}

func htmlTextForParsedGhJson(parsedGhJson map[string]any) (string, error) {
	var htmlText = "<ul>"
	for key, val := range parsedGhJson {
		itemList := fmt.Sprintf("<li> %s: %s </li>", key, val)
		htmlText += itemList
	}
	htmlText += "</ul>"
	return htmlText, nil
}

func userNameTOUserId(userName string) (string, error) {
	userIdLoopUp := make(map[string]string)

	// data-hovercard-id=peterhornyack@google.com
	// data-name=Peter Hornyack , data-member-id= 'user/human/100976612597299399360',
	userIdLoopUp["Dave Revell"] = "100449440289517201826"
	userIdLoopUp["Peter Hornyack"] = "100976612597299399360"
	userIdLoopUp["Rui Zhang"] = "100505938603736143694"
	userIdLoopUp["Qinhang Li"] = "103739887526947033760"
	userIdLoopUp["Jonathan Hong"] = "111198084449821418477"

	return userIdLoopUp[userName], nil
}
