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
	User       string `json:"user"`
	OauthToken string `json:"oauth_token"`
}

func loadConfig() config {
	b, fileErr := ioutil.ReadFile(configFilePath())
	if fileErr != nil {
		fmt.Print(fileErr)
	}

	fmt.Println(b)

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
	return homeDir + "/.config/watch-repos"
}
