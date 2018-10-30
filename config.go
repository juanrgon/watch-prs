package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
)

type config struct {
	Github githubConfig `json:"github.com"`
}

type githubConfig struct {
	Username   string `json:"user"`
	OauthToken string `json:"oauth_token"`
}

func LoadConfig() config {
	b, fileErr := ioutil.ReadFile(configFilePath())
	if fileErr != nil {
		fmt.Print(fileErr)
	}

	var configuration config
	jsonErr := json.Unmarshal(b, &configuration)
	if jsonErr != nil {
		fmt.Println(jsonErr)
	}
	return configuration
}

func configFilePath() string {
	currentUser, _ := user.Current()
	homeDir := currentUser.HomeDir
	return homeDir + "/.config/watch-prs"
}
