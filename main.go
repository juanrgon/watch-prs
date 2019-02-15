package main

import (
	"fmt"
	"os"
	"strconv"
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
		var ri int
		ri = 30
		for i, a := range os.Args[1:] {
			if a == "--only-created" {
				oc = true
			}
			if a == "--only-assigned" {
				oa = true
			}

			if a == "--refresh" {
				val, err := strconv.Atoi(os.Args[i + 2])
				if err != nil {
					fmt.Printf("Invalid --refresh value: %v", os.Args[i+2])
					os.Exit(1)
				}
				ri = val
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
		countDownTillNextRefresh(ri)
		separator = fmt.Sprintf("(%v) %v", time.Now().Format(time.Kitchen), "==========================================================================")
	}
}
