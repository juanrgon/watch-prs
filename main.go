package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	configuration := loadConfig()
	username := configuration.Github.Username

	githubClient := loadGithubClient(configuration) // TODO: Support remote repositories besides github
	org, repo := getCurrentRemoteName()             // TODO: Allow other repos besides the current working directory

	separator := ""
	for {
		var oc bool
		var oa bool
		for _, a := range os.Args[1:] {
			if a == "--only-created" {
				oc = true
			}
			if a == "--only-assigned" {
				oa = true
			}
		}

		var pf pullRequestFilters
		if oc {
			pf.Owner = username
		} else if oa {
			pf.Assignee = username
		} else {
			pf.Owner = username
			pf.Assignee = username
		}

		pulls := getPullRequests(githubClient, org, repo, &pf)

		// TODO: See if this can be wrapped in a decorator
		fmt.Println(separator)
		fmt.Println(clearTerminalSequence)
		printPullStatuses(githubClient, org, repo, pulls)
		countDownTillNextRefresh(30)
		separator = fmt.Sprintf("(%v) %v", time.Now().Format(time.Kitchen), "==========================================================================")
	}
}
