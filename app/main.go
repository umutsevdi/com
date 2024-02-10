package main

import (
	"log"

	_ "github.com/umutsevdi/site/client"
	_ "github.com/umutsevdi/site/sync"
	"github.com/umutsevdi/site/pages"
)

func main() {
	log.Println("Main started")
	pages.Start()
}
