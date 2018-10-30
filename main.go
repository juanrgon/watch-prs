package main

import (
	"fmt"
	"time"
)

func main() {
	configuration := loadConfig()
	username := configuration.Github.Username

	githubClient := loadGithubClient(configuration) // TODO: Support remote repositories besides github
	org, repo := getCurrentRemoteName()             // TODO: Allow other repos besides the current working directory

	separator := ""
	for {
		pullFilters := pullRequestFilters{Owner: username, Assignee: username}
		pulls := getPullRequests(githubClient, org, repo, &pullFilters)

		// TODO: See if this can be wrapped in a decorator
		fmt.Println(separator)
		fmt.Println(clearTerminalSequence)
		printPullStatuses(githubClient, org, repo, pulls)
		countDownTillNextRefresh(30)
		separator = fmt.Sprintf("(%v) %v", time.Now().Format(time.Kitchen), "==========================================================================")
	}
}
