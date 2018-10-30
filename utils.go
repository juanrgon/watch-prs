package main

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
)

func GetCurrentRemoteName() (org string, repo string) {
	localRepo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	origin, err := localRepo.Remote("origin")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	return parseRemoteURL(origin.Config().URLs[0])
}

func parseRemoteURL(url string) (org string, repo string) {
	re := regexp.MustCompile("git@github.com:(?P<org>[a-zA-Z0-9_-]+)/(?P<repo>[a-zA-Z0-9_-]+)")
	matches := re.FindStringSubmatch(url)
	return matches[1], matches[2]
}

func LoadGithubClient(c config) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Github.OauthToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

type PullRequestFilters struct {
	Assignee string
	Owner    string
}

func GetPullRequests(gh *github.Client, org string, repo string, filters *PullRequestFilters) ([]*github.PullRequest, error) {
	pulls, _, err := gh.PullRequests.List(context.Background(), org, repo, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	filteredPulls := []*github.PullRequest{}
	for _, pull := range pulls {
		if pull.GetAssignee().GetLogin() == filters.Assignee {
			filteredPulls = append(filteredPulls, pull)
		} else if pull.GetUser().GetLogin() == filters.Owner {
			filteredPulls = append(filteredPulls, pull)
		}

	}
	return filteredPulls, nil
}

func GetPullRequestCombinedStatus(client *github.Client, org string, repo string, branch *github.PullRequestBranch) *github.CombinedStatus {
	noOpts := github.ListOptions{}
	want, _, err := client.Repositories.GetCombinedStatus(context.Background(), org, repo, branch.GetSHA(), &noOpts)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	return want
}
