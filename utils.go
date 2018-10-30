package main

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/juanrgon/prism"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
)

func getCurrentRemoteName() (org string, repo string) {
	localRepo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Println(prism.InYellow("Current directory is not a git repo."))
		os.Exit(1)
	}
	origin, err := localRepo.Remote("origin")
	if err != nil {
		fmt.Println(prism.InYellow("No remote set for this repo."))
		os.Exit(1)
	}
	return parseRemoteURL(origin.Config().URLs[0])
}

func parseRemoteURL(url string) (org string, repo string) {
	re := regexp.MustCompile("git@github.com:(?P<org>[a-zA-Z0-9_-]+)/(?P<repo>[a-zA-Z0-9_-]+)")
	matches := re.FindStringSubmatch(url)
	return matches[1], matches[2]
}

func loadGithubClient(c config) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Github.OauthToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

type pullRequestFilters struct {
	Assignee string
	Owner    string
}

func getPullRequests(gh *github.Client, org string, repo string, filters *pullRequestFilters) []*github.PullRequest {
	var filteredPulls []*github.PullRequest
	opts := &github.PullRequestListOptions{ListOptions: github.ListOptions{PerPage: 20}}
	for {
		pulls, resp, err := gh.PullRequests.List(context.Background(), org, repo, opts)
		if err != nil {
			fmt.Printf("%v: (%T) %v", prism.InRed("Error getting pull requests statuses from github"), err, err.Error())
			fmt.Printf("\n\n%v: %v", "Please review instructions on creating config file:", prism.InCyan("https://github.com/juanrgon/watch-prs#4-create-a-config-file"))
			os.Exit(1)
		}

		for _, pull := range pulls {
			if filters.Assignee != "" && pull.GetAssignee().GetLogin() == filters.Assignee {
				filteredPulls = append(filteredPulls, pull)
			} else if filters.Owner != "" && pull.GetUser().GetLogin() == filters.Owner {
				filteredPulls = append(filteredPulls, pull)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return filteredPulls
}

func getPullRequestCombinedStatus(client *github.Client, org string, repo string, branch *github.PullRequestBranch) *github.CombinedStatus {
	noOpts := github.ListOptions{}
	want, _, err := client.Repositories.GetCombinedStatus(context.Background(), org, repo, branch.GetSHA(), &noOpts)
	if err != nil {
		fmt.Printf("%v: (%T) %v", prism.InRed("Error getting combined status of "+branch.GetRef()), err, err.Error())
		os.Exit(1)
	}
	return want
}
