package config

import (
	"encoding/json"
	"log"
	"os"
)

var (
	DEFAULT_URI           string = ""
	DEFAULT_DIR           string = "content"
	DEFAULT_PORT          uint64 = 8080
	DEFAULT_PERIOD_MINUTE int    = 30
	DEFAULT_WATCH_SEC     uint64 = 1
)

var instance Config

type GitHubConfig struct {
	User      *string `json:"user,omitempty"`
	Token     *string `json:"token,omitempty"`
	PeriodMin *uint64 `json:"periodMin,omitempty"`
}

type Config struct {
	URI *string `json:"uri,omitempty"`
	// Port to launch the application
	Port *uint64 `json:"port,omitempty"`
	// Path to pages static and components directories
	ContentDir *string       `json:"directory,omitempty"`
	Github     *GitHubConfig `json:"github,omitempty"`
	WatchSec   *uint64       `json:"watchSec"`
	Subdomains []string      `json:"subdomains"`
	Exclude    []string      `json:exclude`
}

// Parses the configuration file at $WEBWATCH_CONFIG or config.json file
// whichever is available.
func init() {
	log.Println("initializing config")
	path, found := os.LookupEnv("WEBWATCH_CONFIG")
	if !found || path == "" {
		path = "config.json"
	}
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("ERROR: config.json is not found at \"", path,
			"\". Either create a config file called \"config.json\", or define a valid",
			"$WEBWATCH_CONFIG.")
	}

	err = json.Unmarshal(file, &instance)
	if err != nil {
		log.Fatal("ERROR: JSON parse", file, ".", err.Error())
	}

	fillEmptyFields()
	log.Println("Server has been started with following configurations:",
		"\n- Target:           ", *instance.URI, ":", *instance.Port,
		"\n- Subdomains:      ", *&instance.Subdomains,
		"\n- HotReload:        ", *instance.WatchSec, "sec",
		"\n- GitHub:           {User   : ", *instance.Github.User, ", Cache: ", *instance.Github.PeriodMin, "min }",
	)

}

// Sanitizes the invalid inputs from the configuration file.
func fillEmptyFields() {
	if instance.URI == nil {
		log.Println("WARN: \"uri\" is not defined at configuration. Continuing with", DEFAULT_URI)
		instance.URI = &DEFAULT_URI
	}
	if instance.Port == nil {
		log.Println("WARN : \"port\" is not defined. Continuing with", DEFAULT_PORT)
		instance.Port = &DEFAULT_PORT
	}
	if instance.ContentDir == nil {
		log.Println("WARN : \"directory\" is not defined. Defaulting to", DEFAULT_DIR)
		*instance.ContentDir = "content"
	}
	if instance.Github == nil || instance.Github.Token == nil || instance.Github.User == nil {
		log.Fatal("ERROR : \"github\" is not properly defined.")
	}
	if instance.Github.PeriodMin == nil {
		log.Println("WARN : \"github.periodMin\" is not defined. Defaulting to", DEFAULT_PERIOD_MINUTE)
		*instance.Github.PeriodMin = uint64(DEFAULT_PERIOD_MINUTE)
	}
	if instance.WatchSec == nil {
		log.Println("WARN : \"WatchSec\" is not defined. Defaulting to", DEFAULT_WATCH_SEC)
		instance.WatchSec = &DEFAULT_WATCH_SEC
	}
}

func ContentDirectory() string { return *instance.ContentDir }
func Port() int                { return int(*instance.Port) }
func URI() string              { return *instance.URI }
func Github() GitHubConfig     { return *instance.Github }
func User() string             { return *instance.Github.User }
func Token() string            { return *instance.Github.Token }
func Period() int              { return int(*instance.Github.PeriodMin) }
func WatchSecond() int         { return int(*instance.WatchSec) }
func Subdomains() *[]string    { return &instance.Subdomains }
func Exclude() *[]string       { return &instance.Exclude }
