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
	var privateKey string
	var appID, installationID int64

	flag.Int64Var(&appID, "app_id", 164400, "github app id")
	flag.Int64Var(&installationID, "installation_id", 22221748, "github installation id")

	flag.StringVar(&privateKey, "private_key", os.Getenv("PRIVATE_KEY"), "")
	flag.Parse()

	log.Println("ok", privateKey[:10])
	itr, err := ghinstallation.NewKeyFromFile(
		http.DefaultTransport,
		appID, installationID,
		privateKey,
	)
	if err != nil {
		panic(err)
		// Handle error.
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
