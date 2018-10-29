package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	configuration := LoadConfig()
	username := configuration.Github.Username

	githubClient := LoadGithubClient(configuration) // TODO: Support remote repositories besides github
	org, repo := GetCurrentRemoteName()             // TODO: Allow other repos besides the current working directory

	separator := ""
	for {
		pullFilters := PullRequestFilters{Owner: username, Assignee: username}
		pulls, err := GetPullRequests(githubClient, org, repo, &pullFilters)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// TODO: See if this can be wrapped in a decorator
		fmt.Println(separator)
		fmt.Println(clearTerminalSequence)
		printPullStatuses(githubClient, org, repo, pulls)
		countDownTillNextRefresh(30)
		separator = fmt.Sprintf("(%v) %v", time.Now().Format(time.Kitchen), "==========================================================================")
	}
}
