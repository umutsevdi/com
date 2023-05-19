package client

import (
	"io/ioutil"
	"net/http"
	"server/syslog"
	"server/util"
	"strings"
)

type response struct {
	Data responseData `json:"data"`
}

type responseData struct {
	User user `json:"user"`
}

type user struct {
	Pinned pinnedItems `json:"pinnedItems"`
}

type pinnedItems struct {
	Edges []item `json:"edges"`
}

type item struct {
	Node repository `json:"node"`
}

type repository struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Url         string       `json:"url"`
	Languages   LanguageEdge `json:"languages"`
	ImageUrl    string       `json:"openGraphImageUrl"`
}
type LanguageEdge struct {
	Edge []LanguageNode `json:"edges"`
}

type LanguageNode struct {
	Node Language `json:"node"`
}

type Language struct {
	Hex  string `json:"color"`
	Name string `json:"name"`
}

// Sends a Graphql Request to the GitHub API and returns the response as string
// @param username username of the requested account
// @return string response if successful
// @return error if response is not successful
func sendRequest(username string) (string, error) {
	URL := "https://api.github.com/graphql"
	q := "{" +
		"  \"query\":" +
		"    \"query MyQuery {" +
		"  user(login: \\\"" + util.Config.User + "\\\"){" +
		"    pinnedItems(first: 6) {" +
		"      edges {" +
		"        node {" +
		"          ... on Repository {" +
		"          name " +
		"          url " +
		"          languages(first: 10) { " +
		"            edges {" +
		"              node { " +
		"                color " +
		"                name" +
		"               }" +
		"             }" +
		"            }" +
		"            description" +
		"            openGraphImageUrl" +
		"           }" +
		"          } " +
		"        } " +
		"       } " +
		"      }" +
		"     }\"" +
		"}"
	req, err := http.NewRequest("POST", URL, strings.NewReader(q))
	req.Header.Add("Authorization", "Bearer "+util.Config.Token)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		syslog.Warn("client | ", err.Error())
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		syslog.Warn("client | ", err.Error())
		return "", err
	}
	bodyResp, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return string(bodyResp), nil
}
