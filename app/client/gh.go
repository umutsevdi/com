package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Language    []Language
	ImageUrl    string `json:"openGraphImageUrl"`
	DoubleSize  bool
}

var repositories [6]Repository
var lock bool

var m map[string]string

func init() {
	m = make(map[string]string)
	m["Shell"] = "/static/img/icon/langbash.png"
	m["C"] = "/static/img/icon/langc.png"
	m["CSS"] = "/static/img/icon/langcss.png"
	m["C++"] = "/static/img/icon/langcpp.png"
	m["CMake"] = "/static/img/icon/langcmake.png"
	m["Dockerfile"] = "/static/img/icon/langdocker.png"
	m["Go"] = "/static/img/icon/langdocker.png"
	m["GDScript"] = "/static/img/icon/langgodot.png"
	m["HTML"] = "/static/img/icon/langhtml.png"
	m["Java"] = "/static/img/icon/langjava.png"
	m["JavaScript"] = "/static/img/icon/langjs.png"
	m["Lua"] = "/static/img/icon/langlua.png"
	m["Makefile"] = "/static/img/icon/langmake.png"
	m["Perl"] = "/static/img/icon/langperl.png"
	m["Python"] = "/static/img/icon/langpython.png"
	m["Jupyter Notebook"] = "/static/img/icon/langpython.png"
	m["PLpgSQL"] = "/static/img/icon/langsql.png"
	m["TypeScript"] = "/static/img/icon/langts.png"
	m["Vim"] = "/static/img/icon/langvim.png"
	m["Vim Snippet"] = "/static/img/icon/langvim.png"
}

func GetGh() []Repository {
	for lock {

	}
	return repositories[0:6]
}

func FetchGh(user, token string) {
	lock = true

	log.Println("Github Fetching")
	buffer, err := sendRequest(user, token)
	if err != nil {
		log.Println("Error: Failed to fetch from Github")
		lock = false
		return
	}
	err = mapResponseToStruct(buffer)
	if err != nil {
		log.Println("Error: Failed to map the result from Github")
		lock = false
		return
	}
	for i, v := range repositories {
		for j := range v.Language {
			v.Language[j].Src = m[v.Language[j].Name]
		}
		repositories[i].DoubleSize = i == 1 || i == 2 || i == 5
	}
	log.Println("Github Done")
	lock = false
}

func mapResponseToStruct(responseBody string) error {
	var r response
	err := json.Unmarshal([]byte(responseBody), &r)
	if err != nil {
		return err
	}

	for i, v := range r.Data.User.Pinned.Edges {
		d, _ := json.Marshal(v.Node)
		json.Unmarshal(d, &repositories[i])
		repositories[i].Language = make([]Language, len(v.Node.Languages.Edge))
		for j, v := range v.Node.Languages.Edge {
			repositories[i].Language[j] = v.Node
		}
	}
	return nil
}

// Sends a Graphql Request to the GitHub API and returns the response as string
// @param username username of the requested account
// @return string response if successful
// @return error if response is not successful
func sendRequest(username, token string) (string, error) {
	URL := "https://api.github.com/graphql"
	q := "{" +
		"  \"query\":" +
		"    \"query MyQuery {" +
		"  user(login: \\\"" + username + "\\\"){" +
		"    pinnedItems(first: 6) {" +
		"      edges {" +
		"        node {" +
		"          ... on Repository {" +
		"          name " +
		"          url " +
		"          languages(first: 3) { " +
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
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Print("client | ", err.Error())
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("client | ", err.Error())
		return "", err
	}
	bodyResp, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return string(bodyResp), nil
}

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
	Src  string
}
