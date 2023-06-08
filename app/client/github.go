package client

import (
	"encoding/json"
	"fmt"
	"server/syslog"
	"server/util"
	"time"
)

var cache Cache
var m map[string]string

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Language    []Language
	ImageUrl    string `json:"openGraphImageUrl"`
}

type Cache struct {
	repositories      *[]Repository
	lastExecutionTime time.Time
}

func init() {
	m = make(map[string]string)
	m["Shell"] = ""
	m["C"] = ""
	m["CSS"] = ""
	m["C++"] = ""
	m["CMake"] = ""
	m["Dockerfile"] = ""
	m["Go"] = ""
	m["GDScript"] = ""
	m["HTML"] = ""
	m["Java"] = " "
	m["JavaScript"] = ""
	m["Lua"] = "󰢱"
	m["Makefile"] = ""
	m["Perl"] = " "
	m["Python"] = " "
	m["Jupyter Notebook"] = ""
	m["PLpgSQL"] = ""
	m["TypeScript"] = "󰛦"
	m["Vim"] = ""
	m["Vim Snippet"] = ""
	updateCache(GetPinnedRepositories(util.Config.User))
	syslog.Info(cache)
}

// Gets the pinned repositories of given username
// @param username username of the account
// @return *[]Repository Repositories of the account
// @return error on failure
func GetPinnedRepositories(username string) (*[]Repository, error) {
	if cache.repositories != nil && time.Now().
		Before(cache.lastExecutionTime.Add(time.Duration(util.Config.ApiCache)*time.Minute)) {
		syslog.Info(cache.lastExecutionTime.String(), " responding with cache")
		return cache.repositories, nil
	}

	syslog.Info("client Expired cache, connecting to API")
	resp, err := sendRequest(username)
	if err != nil {
		return nil, err
	}
	syslog.Info("client API synchronization is completed with www.github.com")
	r, err := toStruct(resp)
	updateCache(r, err)
	return r, err
}

// Returns a single repository as string
// @return string
func (r *Repository) String() string {
	return fmt.Sprintf("{name: %s, description: %s, url: %s, img: %s, language: %s}",
		r.Name,
		r.Description,
		r.Url,
		r.ImageUrl,
		r.Language[0].Name,
	)
}

// Converts the given string into an array of {@link Repository}
// @param responseBody JSON body to convert
// @return *[]Repository Array of repositories
// @return error on failure
func toStruct(responseBody string) (*[]Repository, error) {
	var r response
	err := json.Unmarshal([]byte(responseBody), &r)
	if err != nil {
		return nil, err
	}
	var pinnedItems []Repository = make([]Repository, len(r.Data.User.Pinned.Edges))

	for i, v := range r.Data.User.Pinned.Edges {
		d, _ := json.Marshal(v.Node)
		json.Unmarshal(d, &pinnedItems[i])
		pinnedItems[i].Language = make([]Language, len(v.Node.Languages.Edge))
		for j, v := range v.Node.Languages.Edge {
			pinnedItems[i].Language[j] = v.Node
		}
	}
	return &pinnedItems, nil
}

// Updates most recently stored repository with it's time
// @param *[]Repository slice of gh repositories
// @param err any error
func updateCache(r *[]Repository, err error) {
	if err == nil {
		cache.lastExecutionTime = time.Now()
		cache.repositories = r
	}
}

func UnicodeForLanguage(input string) string {
	if v, ok := m[input]; ok {
		return v
	}
	return input
}
