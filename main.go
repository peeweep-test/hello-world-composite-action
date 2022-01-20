package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

func main() {
	privateKey := []byte(os.Getenv("PRIVATE_KEY"))
	var appID, installationID int64
	var src, dest, message string
	flag.Int64Var(&appID, "app_id", 0, "*github app id")
	flag.Int64Var(&installationID, "installation_id", 0, "*github installation id")
	flag.StringVar(&message, "message", "chore: Sync by "+os.Getenv("GITHUB_REPOSITORY"), "*commit message")
	flag.StringVar(&src, "src", "", "*src path")
	flag.StringVar(&dest, "dest", "", "*dest path")
	flag.Parse()
	if appID == 0 || installationID == 0 || len(src) == 0 || len(dest) == 0 {
		flag.PrintDefaults()
		return
	}
	itr, err := ghinstallation.New(http.DefaultTransport, appID, installationID, []byte(privateKey))
	if err != nil {
		panic(err)
	}
	client := github.NewClient(&http.Client{Transport: itr})
	ctx := context.Background()

	arr := strings.SplitN(dest, "/", 3)
	if len(arr) != 3 {
		log.Fatal("wrong dist. example: username/repo/file")
	}
	owner := arr[0]
	repo := arr[1]
	path := arr[2]

	fileContent, _, resp, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		if resp.StatusCode != http.StatusNotFound {
			panic(err)
		}
	}
	var sha string
	if fileContent != nil {
		sha = fileContent.GetSHA()
	}
	content := []byte(time.Now().String())
	_, _, err = client.Repositories.UpdateFile(
		ctx, owner, repo, path,
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
