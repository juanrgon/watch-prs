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
	b, err := ioutil.ReadFile(configFilePath())
	var m error
	switch et := err.(type) {
	case nil:
	case *os.PathError:
		m = fmt.Errorf("\n%v: %v", prism.InRed("Could not open config file"), et.Path)
	default:
		m = fmt.Errorf("\n%v: %v", prism.InRed("Unexpected error reading config file"), err.Error())
	}
	if err != nil {
		fmt.Println(m)
		os.Exit(1)
	}

	err = json.Unmarshal(b, &c)
	switch err.(type) {
	case nil:
	default:
		m = fmt.Errorf("\n%v: %v", prism.InRed("Unexpected error reading config file"), err.Error())
	}
	if err != nil {
		m := fmt.Sprintf("\nCould not open config file: %v", err.Error())
		panic(m)
	}
	return
}

func configFilePath() string {
	currentUser, _ := user.Current()
	homeDir := currentUser.HomeDir
	return homeDir + "/.config/watch-prs"
}
