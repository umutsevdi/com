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
	//page := "//Title: An example title\n" +
	//	"// Author: Umut Sevdi\n" +
	//	"# Hello everyone\n" +
	//	"## This is a subtitle\n" +
	//	"###This is a joined subtitle\n" +
	//	"This is a raw string\n" +
	//	"=>www.google.com\n" +
	//	"=> www.duckduckgo.com\n" +
	//	"=> www.umutsevdi.com My Personal Site\n" +
	//	"> Some quote\n" +
	//	"> Some another quote\n" +
	//	"```c\n" +
	//	"int main(void) {\n" +
	//	"    printf(\"Hello world\\n\");\n" +
	//	"}\n" +
	//	"```" +
	//	"---" +
	//	"Bye"
	//var g index.Gemtext
	//g.Parse(page)

	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	go pages.Serve(e)
	str := ":" + strconv.Itoa(int(*config.C.Port))
	log.Fatal(e.Start(str))

}
