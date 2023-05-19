package util

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"server/syslog"
	"strconv"
	"strings"
)

type Configuration struct {
	User     string
	Ip       string
	Port     int
	LogLevel syslog.LOG_LEVEL
	Cache    bool
	ApiCache int
	Token    string
	Site     string
}

var Config Configuration

// Reads the app.config and initializes the Configuration struct
func init() {
	file, err := ioutil.ReadFile("app.config")
	if err != nil {
		log.Fatal("Error: app.config was not found")
	}
	log.Println("Reading app.config")
	lines := strings.Split(string(file), "\n")
	for i := range lines {
		line := strings.Split(lines[i], "=")
		if len(line) <= 1 {
			continue
		}
		switch line[0] {
		case "cache-time":
			Config.ApiCache, err = strconv.Atoi(line[1])
			if err != nil {
				log.Fatal("Error:  must be an integer")
			}
		case "token":
			Config.Token = line[1]
		case "user":
			Config.User = line[1]
		case "port":
			Config.Port, err = strconv.Atoi(line[1])
			if err != nil {
				log.Fatal("Error: port must be an integer")
			}
		case "ip":
			Config.Ip = line[1]
		case "site":
			Config.Site = line[1]
		case "cache":
			Config.Cache, err = strconv.ParseBool(line[1])
		case "log-level":
			v, _ := strconv.Atoi(line[1])
			Config.LogLevel = syslog.LOG_LEVEL(v)
		}
	}
	if Config.User == "" {
		log.Fatal("Error: GitHub account is not defined")
	}
	if Config.Token == "" {
		log.Fatal("Error: Missing API token")
	}
	if Config.Ip == "" {
		log.Println("Warning: Ip was not found in the configuration file, fallback to localhost")
	}
	if Config.Site == "" {
		log.Fatal("Error: site is not defined")
	}
	if Config.Port == 0 {
		log.Println("Warning: Port was not found in configuration file, fallback to :8080")
		Config.Port = 8080
	}
	syslog.SetLogLevel(Config.LogLevel)
	if Config.Cache {
		log.Println("File-system caching enabled. Caching directory")
		traverse()
	}
	log.Println("Configuration file was parsed successfully")
}

// Returns the requested file in the /content/ directory
//
// - If file-system caching is enabled, returns from the cache
//
// @param URL string path to file
// @return []byte content of the file, if file exist, nil otherwise
// @return nil if the file exists, error type otherwise
func ReadByteFrom(url string) ([]byte, error) {
	if Config.Cache {
		if v, ok := Cache[url]; ok {
			return v, nil
		}
		return nil, errors.New("not found")
	}
	file, err := os.ReadFile("content/" + url)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Parses received URL and extracts it's extension
// @param URL path to file
// @return string corresponding extension type in the format of .type
func Ext(url string) string {
	p := strings.Split(url, "/")
	fname := strings.Split(p[len(p)-1], ".")
	return "." + fname[len(fname)-1]
}

// Parses received URL and extracts it's name
// @param URL path to file
// @return file name
func File(url string) string {
	return strings.Split(url, ".")[0]
}
