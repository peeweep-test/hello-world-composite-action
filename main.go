package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

func main() {
	privateKey := []byte(os.Getenv("PRIVATE_KEY"))
	repositor := os.Getenv("GITHUB_REPOSITORY")
	var appID, installationID int64
	flag.Int64Var(&appID, "app_id", 164400, "github app id")
	flag.Int64Var(&installationID, "installation_id", 22221748, "github installation id")
	flag.Parse()

	itr, err := ghinstallation.New(http.DefaultTransport, appID, installationID, []byte(privateKey))
	if err != nil {
		panic(err)
	}
	client := github.NewClient(&http.Client{Transport: itr})
	ctx := context.Background()

	owner := "peeweep-test"
	repo := "test-action"
	path := "hello"

	fileContent, _, resp, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		log.Println(err)
		if resp.StatusCode != http.StatusNotFound {
			panic(err)
		}
	}
	var sha string
	if fileContent != nil {
		sha = fileContent.GetSHA()
	}
	log.Println("exists file:", fileContent)
	message := "chore: Sync by " + repositor
	content := []byte(time.Now().String())
	_, _, err = client.Repositories.UpdateFile(ctx,
		owner, repo, path,
		&github.RepositoryContentFileOptions{
			Message: &message,
			Content: content,
			SHA:     &sha,
		},
	)
	if err != nil {
		panic(err)
	}
}
