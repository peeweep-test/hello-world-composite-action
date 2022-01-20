package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

func main() {
	var appID, installationID int64

	flag.Int64Var(&appID, "app_id", 164400, "github app id")
	flag.Int64Var(&installationID, "installation_id", 22221748, "github installation id")

	flag.Parse()

	itr, err := ghinstallation.New(http.DefaultTransport, appID, installationID, []byte(os.Getenv("PRIVATE_KEY")))
	if err != nil {
		panic(err)
	}

	client := github.NewClient(&http.Client{Transport: itr})
	client.Repositories.UpdateFile(context.Background(),
		"peeweep-test", "test-action", "hello",
		&github.RepositoryContentFileOptions{})
	repos, _, err := client.Apps.ListRepos(context.Background(), &github.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Println(repos)
}
