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

func printPullStatuses(client *github.Client, org string, repo string, pulls []*github.PullRequest) {
	for {
		fmt.Println(clearTerminalSequence)
		for _, pull := range pulls {
			branch := pull.GetHead()
			status := GetPullRequestCombinedStatus(client, org, repo, branch)
			success := "success"
			fmt.Printf("%v: %s %s\n", coloredState(*status.State), branch.GetRef(), prism.InMagenta(pull.GetHTMLURL()))
			if status.State != &success {
				for _, status := range status.Statuses {
					fmt.Printf("    %v: %v\n", coloredState(*status.State), prism.InCyan(*status.TargetURL))
				}
			}
			fmt.Println()
		}
		countDownTillNextRefresh(30)
		fmt.Println("=========================================")
	}
}

func coloredState(s string) (ds prism.DecoratedString) {
	switch s {
	case "failure":
		ds = prism.InRed(s)
	case "success":
		ds = prism.InGreen(s)
	case "pending":
		ds = prism.InYellow(s)
	}
	return
}

func countDownTillNextRefresh(s int) {
	for i := s; i > 0; i-- {
		message := fmt.Sprintf("Refreshing in %d", i)
		overwriteLine(message)
		time.Sleep(1000 * time.Millisecond)
	}
}
