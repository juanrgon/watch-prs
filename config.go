package main

import (
	"encoding/json"
	"fmt"
	"github.com/juanrgon/prism"
	"io/ioutil"
	"os"
	"os/user"
)

type config struct {
	Github githubConfig `json:"github.com"`
}

type githubConfig struct {
	Username   string `json:"user"`
	OauthToken string `json:"oauth_token"`
}

func loadConfig() (c config) {
	p := configFilePath()
	b, err := ioutil.ReadFile(p)
	var m error
	switch err.(type) {
	case nil:
	case *os.PathError:
		m = fmt.Errorf("\n%v: %v", prism.InRed("Could not open config file"), p)
	default:
		m = fmt.Errorf("\n%v: (%T) %v", prism.InRed("Unexpected error reading config file " + p), err, err.Error())
	}
	if err != nil {
		fmt.Println(m)
		os.Exit(1)
	}

	err = json.Unmarshal(b, &c)
	switch err.(type) {
	case nil:
	case *json.SyntaxError:
		m = fmt.Errorf("\n%v: %v", prism.InRed("Invalid JSON in config file " + p), err.Error())
	default:
		m = fmt.Errorf("\n%v: (%T) %v", prism.InRed("Unexpected error parsing config file " + p), err, err.Error())
	}
	if err != nil {
		fmt.Println(m)
		fmt.Printf("\n%v: %v", "Please review instructions on creating config file:", prism.InCyan("https://github.com/juanrgon/watch-prs#create-a-config-file"))
		os.Exit(1)
	}
	return
}

func configFilePath() string {
	currentUser, _ := user.Current()
	homeDir := currentUser.HomeDir
	return homeDir + "/.config/watch-prs"
}
