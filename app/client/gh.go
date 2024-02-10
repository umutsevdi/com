package client

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/umutsevdi/site/config"
)

type Repository struct {
	Name        string
	Description string
	Url         string
	Language    []Language
	ImageUrl    string
	DoubleSize  bool
	Stars       int
	Forks       int
	License     string
	isPinned    bool
}

var pinned [6]Repository = [6]Repository{}
var repoList []Repository = []Repository{}

var lock bool

var m map[string]string

func init() {
	m = getLanguageMap()
	gh := config.Github()
	go GitHubBatch(config.User(), config.Token())
	ticker := time.NewTicker(time.Duration(*gh.PeriodMin) * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				GitHubBatch(config.User(), config.Token())
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func GetGh() []Repository {
	for lock {

	}
	return pinned[0:6]
}

func GitHubBatch(user, token string) {
	log.Println("Updating GitHub cache")
	lock = true
	sendRepoListQuery(user, token)
	sendPinnedRepositoryQuery(user, token)
	lock = false
	log.Println("GitHub cached successfully")
}

func sendRepoListQuery(user, token string) {
	buffer, err := sendRequest(user, PINNED_QUERY, token)
	if err != nil {
		log.Println("ERROR: Failed to fetch from Github")
		lock = false
		return
	}

	var repos pinnedRepositories
	err = json.Unmarshal(buffer, &repos)
	repos.setRepository(&pinned)

	if err != nil {
		log.Println("ERROR: Failed to map the result from Github.", err)
		lock = false
		return
	}
}

func sendPinnedRepositoryQuery(user, token string) {
	buffer, err := sendRequest(user, REPO_QUERY, token)
	if err != nil {
		log.Println("ERROR: Failed to fetch from Github")
		lock = false
		return
	}

	var repos repositoryList
	err = json.Unmarshal(buffer, &repos)
	repos.setRepository(&repoList)

	if err != nil {
		log.Println("ERROR: Failed to map the result from Github.", err)
		lock = false
		return
	}

}

// Sends a Graphql Request to the GitHub API and returns the response as string
// @param username username of the requested account
// @param query graphql query to fetch
// @return string response if successful
// @return error if response is not successful
func sendRequest(username, query, token string) ([]byte, error) {
	URL := "https://api.github.com/graphql"
	q := strings.ReplaceAll(strings.Replace(query, "__USER__", username, 1), "\n", " ")
	q = "{\"query\": \"query " + q + "\"}"
	req, err := http.NewRequest("POST", URL, strings.NewReader(q))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Print("client | ", err.Error())
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("client | ", err.Error())
		return nil, err
	}
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERROR: Error during read", err)
	}
	return bodyResp, nil
}
