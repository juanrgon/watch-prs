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
	for _, p := range pulls {
		url := p.GetHTMLURL()
		branch := p.GetHead()
		name := branch.GetRef()

		ciStatus := getPullRequestCombinedStatus(client, org, repo, branch)
		mergeable := p.GetMergeable()
		var overallStatus state
		if mergeable && state(*ciStatus.State) == success {
			overallStatus = success
		} else {
			overallStatus = failure
		}

		printPullBranch(name, url, overallStatus)
		printPullMergeable(mergeable)
		printCIStatus(ciStatus)
		fmt.Println()
	}
}

func printPullBranch(name string, url string, status state) {
	cn := coloredByState(name, status)
	fmt.Printf("%s %s\n", cn, prism.InMagenta(url))
}

func printPullMergeable(m bool) {
	s := "mergeable"
	if m {
		s = string(prism.InGreen(s))
	} else {
		s = string(prism.InRed(s))
	}
	fmt.Printf("    %s\n", s)
}

func printCIStatus(s *github.CombinedStatus) {
	header := "ci-status"
	combinedState := state(*s.State)

	if state(combinedState) == success {
		fmt.Printf("    %s\n", coloredByState(header, success))
	}else{
		fmt.Printf("    %s:\n", coloredByState(header, combinedState))
		for _, status := range s.Statuses {
			fmt.Printf("        %v: %v\n", coloredState(*status.State), prism.InCyan(*status.TargetURL))
		}
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
