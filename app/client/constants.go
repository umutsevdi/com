package client

import (
	"slices"
	"strings"
)

func getLanguageMap() map[string]string {
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
	return m
}

type Language struct {
	Hex  string `json:"color"`
	Name string `json:"name"`
	Src  string
}

const PINNED_QUERY = `{
  user(login: \"__USER__\") {
    pinnedItems(first: 6) {
      nodes {
        ... on Repository {
          name
          description
          url
          stargazerCount
          forkCount
          openGraphImageUrl
          licenseInfo {
            nickname
          }
          languages(first: 3, orderBy: {field: SIZE, direction: DESC}) {
            nodes {
              color
              name
            }
          }
        }
      }
    }
  }
}`

type pinnedRepositories struct {
	Data struct {
		User struct {
			Pinned struct {
				Nodes []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					Url         string `json:"url"`
					Stars       int    `json:"stargazerCount"`
					Forks       int    `json:"forkCount"`
					Languages   struct {
						Nodes []Language `json:"nodes"`
					} `json:"languages"`
					ImageUrl    string `json:"openGraphImageUrl"`
					LicenseInfo struct {
						Nickname string `json:"nickname"`
					} `json:"licenseInfo"`
				} `json:"nodes"`
			} `json:"pinnedItems"`
		} `json:"user"`
	} `json:"data"`
}

func (p *pinnedRepositories) setRepository(r *[6]Repository) {
	for i, v := range p.Data.User.Pinned.Nodes {
		r[i] = Repository{
			Description: v.Description,
			Url:         v.Url,
			Stars:       v.Stars,
			Forks:       v.Forks,
			ImageUrl:    v.ImageUrl,
			License:     v.LicenseInfo.Nickname,
			Language:    make([]Language, len(v.Languages.Nodes)),
		}
		for j, v := range v.Languages.Nodes {
			r[i].Language[j] = Language{
				Src:  m[v.Name],
				Name: v.Name,
				Hex:  v.Hex,
			}
			r[i].DoubleSize = i == 1 || i == 2 || i == 5
		}
		names := []string{}
		for _, v := range strings.Split(v.Name, "-") {
			names = append(names, strings.ToUpper(v[:1])+v[1:])
		}
		r[i].Name = strings.Join(names, " ")
	}
}

const REPO_QUERY = `{
  user(login: \"__USER__\") {
    repositories(
      privacy: PUBLIC
      ownerAffiliations: OWNER
      first: 20
      orderBy: {field: CREATED_AT, direction: DESC}
    ) {
      nodes {
        name
        description
        url
        stargazerCount
        forkCount
        primaryLanguage {
          name
          color
        }
        licenseInfo {
          nickname
        }
      }
    }
  }
}`

type repositoryList struct {
	Data struct {
		User struct {
			Repository struct {
				Nodes []struct {
					Name        string   `json:"name"`
					Description string   `json:"description"`
					Url         string   `json:"url"`
					Stars       int      `json:"stargazerCount"`
					Forks       int      `json:"forkCount"`
					Langauge    Language `json:"primaryLanguage"`
					LicenseInfo struct {
						Nickname string `json:"nickname"`
					} `json:"licenseInfo"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"user"`
	} `json:"data"`
}

func (p *repositoryList) setRepository(r *[]Repository) {
	*r = slices.Delete(*r, 0, len(*r))
	for _, v := range p.Data.User.Repository.Nodes {
		repo := Repository{
			Description: v.Description,
			Url:         v.Url,
			Stars:       v.Stars,
			Forks:       v.Forks,
			License:     v.LicenseInfo.Nickname,
			Language:    []Language{v.Langauge},
		}
		names := []string{}
		for _, v := range strings.Split(v.Name, "-") {
			names = append(names, strings.ToUpper(v[:1])+v[1:])
		}
		repo.Name = strings.Join(names, " ")
		*r = append(*r, repo)

	}
}
