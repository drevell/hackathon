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
	"fmt"
	"net/http"
	"os"
)

const (
	baseUrl string = "https://chat.googleapis.com/v1/spaces/%s/messages?key=%s&token=%s"
	space   string = "AAAAUJgrNvE"
	key     string = "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI"
	token   string = "VyyAuvWCmbwWkXtytXxNIc5R_xX41_fi4WxieHsSggs%3D"
)

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func realMain() error {
	/* TODO: use this flags example later
	f := flag.NewFlagSet("", flag.ExitOnError)
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\n", strings.TrimSpace(lintCommandHelp))
		f.PrintDefaults()
	}
	showVersion := f.Bool("version", false, "display version information")

	if err := f.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	if *showVersion {
		fmt.Fprintln(os.Stderr, terraformlinter.HumanVersion)
		return nil
	}
	*/

	url := fmt.Sprintf(baseUrl, space, key, token)
	fmt.Println("url: ", url)
	var jsonData = []byte(`{
		"text": "Hey <users/100449440289517201826>"
	}`)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("creating http request failed: %w", err)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("sending http request failed: %w", err)
	}
	fmt.Println("resp: ", resp)
	defer resp.Body.Close()

	return nil
}
