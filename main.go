package main

import (
	"fmt"
	"os"
)

func main() {
	configuration := LoadConfig()
	username := configuration.Github.Username

	githubClient := LoadGithubClient(configuration) // TODO: Support remote repositories besides github
	org, repo := GetCurrentRemoteName()             // TODO: Allow other repos besides the current working directory

	pullFilters := PullRequestFilters{Owner: username, Assignee: username}
	pulls, err := GetPullRequests(githubClient, org, repo, &pullFilters)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	printPullStatuses(githubClient, org, repo, pulls)
}
