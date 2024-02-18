package client

import (
	"slices"
	"strings"
)

func getLanguageMap() map[string]string {
	m = make(map[string]string)
	m["Shell"] = "/static/img/icon/lang/bash.png"
	m["C"] = "/static/img/icon/lang/c.png"
	m["CSS"] = "/static/img/icon/lang/css.png"
	m["C++"] = "/static/img/icon/lang/cpp.png"
	m["CMake"] = "/static/img/icon/lang/cmake.png"
	m["Dart"] = "/static/img/icon/lang/dart.png"
	m["Dockerfile"] = "/static/img/icon/lang/docker.png"
	m["Go"] = "/static/img/icon/lang/go.png"
	m["GDScript"] = "/static/img/icon/lang/godot.png"
	m["HTML"] = "/static/img/icon/lang/html.png"
	m["Java"] = "/static/img/icon/lang/java.png"
	m["JavaScript"] = "/static/img/icon/lang/js.png"
	m["Lua"] = "/static/img/icon/lang/lua.png"
	m["Makefile"] = "/static/img/icon/lang/make.png"
	m["Perl"] = "/static/img/icon/lang/perl.png"
	m["Python"] = "/static/img/icon/lang/python.png"
	m["Jupyter Notebook"] = "/static/img/icon/lang/python.png"
	m["PLpgSQL"] = "/static/img/icon/lang/sql.png"
	m["TypeScript"] = "/static/img/icon/lang/ts.png"
	m["Vim"] = "/static/img/icon/lang/vim.png"
	m["Vim Snippet"] = "/static/img/icon/lang/vim.png"
	return m
}

type Language struct {
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
          licenseInfo {
            name
            nickname
          }
          languages(first: 3, orderBy: {field: SIZE, direction: DESC}) {
            nodes {
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
					LicenseInfo struct {
						Nickname string `json:"nickname"`
						Name     string `json:"name"`
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
			License:     v.LicenseInfo.Nickname,
			Language:    make([]Language, len(v.Languages.Nodes)),
		}
		if v.LicenseInfo.Nickname == "" {
			r[i].License = v.LicenseInfo.Name
		}
		for j, v := range v.Languages.Nodes {
			r[i].Language[j] = Language{Src: m[v.Name], Name: v.Name}
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
      first: 25
      orderBy: {field: PUSHED_AT, direction: DESC}
    ) {
      nodes {
        name
        description
        url
        stargazerCount
        forkCount
        licenseInfo {
          name
          nickname
        }
        primaryLanguage {
          name
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
					Name        string `json:"name"`
					Description string `json:"description"`
					Url         string `json:"url"`
					Stars       int    `json:"stargazerCount"`
					Forks       int    `json:"forkCount"`
					LicenseInfo struct {
						Nickname string `json:"nickname"`
						Name     string `json:"name"`
					} `json:"licenseInfo"`
					Language struct {
						Name string `json:"name"`
					} `json:"primaryLanguage"`
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
			License:     v.LicenseInfo.Nickname,
			Stars:       v.Stars,
			Forks:       v.Forks,
			Language:    []Language{{Name: v.Language.Name, Src: m[v.Language.Name]}},
		}
		if v.LicenseInfo.Nickname == "" {
			repo.License = v.LicenseInfo.Name
		}

		names := []string{}
		for _, v := range strings.Split(v.Name, "-") {
			names = append(names, strings.ToUpper(v[:1])+v[1:])
		}
		isPinned := false
		repo.Name = strings.Join(names, " ")
		for _, p := range pinned {
			if p.Name == repo.Name {
				isPinned = true
			}
		}
		if !isPinned && len(repo.Description) > 0 && len(*r) < 16 {
			*r = append(*r, repo)
		}

	}
}
