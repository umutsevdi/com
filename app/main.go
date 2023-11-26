package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"strconv"
	"umutsevdi/com/config"
	_ "umutsevdi/com/index"
	"umutsevdi/com/pages"
)

func main() {
	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	go pages.Serve(e)
	str := ":" + strconv.Itoa(int(*config.C.Port))
	log.Fatal(e.Start(str))

}
