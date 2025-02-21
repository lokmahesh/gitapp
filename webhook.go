package main

import (
	"bytes"
	// "context"
	// "encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/go-github/v50/github"
)

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(os.Getenv("GITHUB_WEBHOOK_SECRET")))
	if err != nil {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		http.Error(w, "Error parsing webhook", http.StatusBadRequest)
		return
	}

	switch e := event.(type) {
	case *github.PushEvent:
		repo := e.Repo.GetFullName()
		commitSHA := e.HeadCommit.GetID()
		go processCommit(repo, commitSHA)
	}

	w.WriteHeader(http.StatusOK)
}

func processCommit(repo string, sha string) {
	log.Printf("Processing commit %s in %s\n", sha, repo)

	// Create a Check Run
	checkRunID := createCheckRun(repo, sha)

	// Clone the repository
	cmd := exec.Command("git", "clone", "https://github.com/"+repo+".git")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("Error cloning repo:", err)
		updateCheckRun(repo, checkRunID, "failure", "Failed to clone repository")
		return
	}

	// Run commit linting
	result := runCommitLint(repo)

	// Update GitHub Check Run
	updateCheckRun(repo, checkRunID, "completed", result)
}
