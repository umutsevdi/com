package util

import (
	"os"
	"os/exec"
	"server/syslog"
	"strings"
)

var Cache map[string][]byte = make(map[string][]byte)

func traverse() {
	exec.Command("/bin/bash", "min.sh").Run()
	_, err := os.ReadDir("content")
	if err != nil {
		syslog.Fatal("Cache directory was not found")
		return
	}
	mapFiles("")
	for k := range Cache {
		syslog.Info("cache:", k)
	}
}

func mapFiles(path string) {
	data, err := os.ReadDir("content/" + path)
	if err != nil {
		syslog.Fatal("An error occurred while caching:" + err.Error())
	}

	for _, v := range data {
		if v.Type().IsDir() {
			mapFiles(path + "/" + v.Name())
		} else {
			d, err := os.ReadFile("content/" + path + "/" + v.Name())
			if _, ok := Cache[path+"/"+v.Name()]; err == nil && !ok {
				if strings.Contains(v.Name(), "-min") {
					Cache[path+"/"+strings.Replace(v.Name(), "-min", "", 1)] = d
				} else {
					Cache[path+"/"+v.Name()] = d
				}
			}
		}
	}
}

func ListPages(path string, routes *[]string) {
	data, _ := os.ReadDir("content/" + path)
	for _, v := range data {
		if v.Type().IsDir() {
			ListPages(path+"/"+v.Name(), routes)
		} else if !strings.Contains(v.Name(), "-min") {
			*routes = append(*routes, path+"/"+v.Name())
		}
	}
}
