package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

// GitHub access token
var githubToken = "YOUR_GITHUB_APP_ACCESS_TOKEN"

func createCheckRun(repo string, sha string) int64 {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	ownerRepo, _ := splitOwnerAndRepo(repo)

	opts := github.CreateCheckRunOptions{
		Name:    "Commit Linter",
		HeadSHA: sha,
		Status:  github.String("in_progress"),
	}

	checkRun, _, err := client.Checks.CreateCheckRun(ctx, ownerRepo[0], ownerRepo[1], opts)
	if err != nil {
		log.Println("Error creating check run:", err)
		return 0
	}
	return checkRun.GetID()
}

func updateCheckRun(repo string, checkRunID int64, status, result string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	ownerRepo, _ := splitOwnerAndRepo(repo)

	conclusion := "success"
	if status == "failure" {
		conclusion = "failure"
	}

	opts := github.UpdateCheckRunOptions{
		Status:     github.String(status),
		Conclusion: github.String(conclusion),
		Output: &github.CheckRunOutput{
			Title:   github.String("Commit Lint Results"),
			Summary: github.String(result),
		},
	}

	_, _, err := client.Checks.UpdateCheckRun(ctx, ownerRepo[0], ownerRepo[1], checkRunID, opts)
	if err != nil {
		log.Println("Error updating check run:", err)
	}
}

// Helper function to split "owner/repo"
func splitOwnerAndRepo(repo string) ([]string, error) {
	chunk := strings.Split(repo, "/")
	if len(chunk) != 2 {
		return nil, fmt.Errorf("invalid repository name: '%s'", repo)
	}
	return chunk, nil
}
