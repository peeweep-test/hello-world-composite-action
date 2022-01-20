package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v42/github"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type syncProject struct {
	File   string   `json:"Name"`
	Target []string `json:"Repos"`
}

func main() {
	privateKey := []byte(os.Getenv("PRIVATE_KEY"))
	var appID, installationID int64
	var src, dest, message string
	flag.Int64Var(&appID, "app_id", 0, "*github app id")
	flag.Int64Var(&installationID, "installation_id", 0, "*github installation id")
	flag.StringVar(&message, "message", "chore: Sync by "+os.Getenv("GITHUB_REPOSITORY"), "*commit message")
	flag.Parse()
	if appID == 0 || installationID == 0 {
		flag.PrintDefaults()
		return
	}
	//itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appID, installationID, "123.pem")
	itr, err := ghinstallation.New(http.DefaultTransport, appID, installationID, []byte(privateKey))
	if err != nil {
		panic(err)
	}
	var s []syncProject
	client := github.NewClient(&http.Client{Transport: itr})
	ctx := context.Background()
	//_ = github.NewClient(&http.Client{Transport: itr})
	//_ = context.Background()

	jsonFile, err := os.Open("sync.json")
	if err != nil {
		log.Fatal(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &s)

	for i := 0; i < len(s); i++ {
		src = s[i].File
		for j := 0; j < len(s[i].Target); j++ {
			dest = s[i].Target[j]
			log.Println(src + " to " + dest + "/" + src)

			arr := strings.SplitN(dest, "/", 2)

			if len(arr) != 2 {
				log.Fatal("wrong dist. example: username/repo")
			}

			owner := arr[0]
			repo := arr[1]
			path := src

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

			f, err := os.Open(src)
			if err != nil {
				log.Fatal(err)
			}
			byteValue, _ := ioutil.ReadAll(f)
			content := byteValue

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
	}
}
