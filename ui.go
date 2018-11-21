package main

import (
	"fmt"
	"time"

	"github.com/google/go-github/github"
	"github.com/juanrgon/prism"
)

const clearTerminalSequence = "\033[H\033[J"
const clearAllAfterCursorSequence = "\033[K"

func overwriteLine(text string) {
	fmt.Print(clearAllAfterCursorSequence, text, "\r")
}

type state string

const (
	success state = state("success")
	pending state = state("pending")
	failure state = state("failure")
)

func printPullStatuses(client *github.Client, org string, repo string, pulls []*github.PullRequest) {
	var pendingPulls []*github.PullRequest
	var failingPulls []*github.PullRequest
	for _, p := range pulls {
		branch := p.GetHead()

		ciStatus := getPullRequestCombinedStatus(client, org, repo, branch)
		ciState := state(ciStatus.GetState())
		if ciState == success {
			printPull(p, client, org, repo)
		} else if ciState == pending {
			pendingPulls = append(pendingPulls, p)
		} else {
			failingPulls = append(failingPulls, p)
		}
	}

	for _, p := range pendingPulls {
		printPull(p, client, org, repo)
	}

	for _, p := range failingPulls {
		printPull(p, client, org, repo)
	}
}

func printPull(p *github.PullRequest, client *github.Client, org string, repo string) {
	url := p.GetHTMLURL()
	branch := p.GetHead()
	name := branch.GetRef()

	ciStatus := getPullRequestCombinedStatus(client, org, repo, branch)
	ciState := state(ciStatus.GetState())

	printPullBranch(name, url, ciState)
	if state(ciStatus.GetState()) != success {
		printCIStatus(ciStatus)
	}
	fmt.Println()
}

func printPullBranch(name string, url string, ci state) {
	var ciStatus prism.DecoratedString
	if ci == success {
		ciStatus = coloredByState("passing", success)
	} else if ci == pending {
		ciStatus = coloredByState("pending", pending)
	} else {
		ciStatus = coloredByState("failing", failure)
	}
	fmt.Println(prism.InMagenta(url))
	fmt.Printf("%s: %s \n", ciStatus, name)
}

func printPullMergeable(m bool) {
	if !m {
		s := string(prism.InRed("merge conflict or failed review"))
		fmt.Printf("    %s\n", s)
	}
}

func printCIStatus(s *github.CombinedStatus) {
	for _, status := range s.Statuses {
		fmt.Printf("    %v: %v\n", coloredState(*status.State), prism.InCyan(*status.TargetURL))
	}
}

func coloredByState(t string, s state) (ds prism.DecoratedString) {
	switch s {
	case failure:
		ds = prism.InRed(t)
	case success:
		ds = prism.InGreen(t)
	case pending:
		ds = prism.InYellow(t)
	}
	return
}

func coloredState(s string) (ds prism.DecoratedString) {
	return coloredByState(s, state(s))
}

func countDownTillNextRefresh(s int) {
	for i := s; i > 0; i-- {
		message := fmt.Sprintf("Refreshing in %d", i)
		overwriteLine(message)
		time.Sleep(1000 * time.Millisecond)
	}
	overwriteLine("Refreshing")
}
