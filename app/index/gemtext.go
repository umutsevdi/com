package index

import (
	"fmt"
	"strings"
	"time"
)

type GEMTEXT int

var GEMTEXT_STR []string = []string{
	"GEMTEXT_TEXT ",
	"GEMTEXT_LINK ",
	"GEMTEXT_LIST ",
	"GEMTEXT_QUOTE",
	"GEMTEXT_H3   ",
	"GEMTEXT_H2   ",
	"GEMTEXT_H1   ",
	"GEMTEXT_PRE  ",
}

const (
	GEMTEXT_TEXT = GEMTEXT(iota)
	GEMTEXT_LINK
	GEMTEXT_LIST
	GEMTEXT_QUOTE
	GEMTEXT_H3
	GEMTEXT_H2
	GEMTEXT_H1
	GEMTEXT_PRE
)

type Gemtext struct {
	Title   string
	Author  string
	Created time.Time
	Updated time.Time
	Lines   *[]Line
}

type Line struct {
	Type GEMTEXT
	Text string
	Arg  string
}

func (g *Gemtext) Parse(text string) {
	linesRaw := strings.Split(text, "\n")

	lines := make([]Line, 0, len(linesRaw))
	var pre bool = false
	var preLastHint string = ""
	for _, v := range linesRaw {
		if pre {
			if strings.HasPrefix(v, "```") {
				pre = false
				preLastHint = ""
				continue
			}
			lines = append(lines, Line{Type: GEMTEXT_PRE, Arg: preLastHint, Text: v})
			continue
		}
		if strings.HasPrefix(v, "//") {
			g.parseHeader(v)
		} else if strings.HasPrefix(v, "```") {
			pre = true
			preLastHint = strings.Split(v, " ")[0][3:]
		} else if strings.HasPrefix(v, ">") {
			lines = append(lines, parse(v[1:], ">"))
		} else if strings.HasPrefix(v, "*") {
			lines = append(lines, parse(v[1:], "*"))
		} else if strings.HasPrefix(v, "=>") {
			lines = append(lines, parseLink(v))
		} else if strings.HasPrefix(v, "#") {
			lines = append(lines, parseH(v))
		} else {
			lines = append(lines, Line{Text: strings.TrimSpace(v), Type: GEMTEXT_TEXT})
		}
	}
	for _, v := range lines {
		if v.Type != GEMTEXT_PRE {
			v.Text = strings.TrimSpace(v.Text)
			v.Arg = strings.TrimSpace(v.Arg)
		}
	}
	g.Lines = &lines

	for _, v := range *g.Lines {
		fmt.Println("{LineType:"+GEMTEXT_STR[v.Type], "\tText:"+v.Text, "\tArg:"+v.Arg, "}")
	}

}

func (g *Gemtext) parseHeader(line string) {
	s := (line)[2:]
	sArray := strings.Split(s, "=")
	if len(sArray) != 2 {
		return
	}
	if sArray[0] == "Title" {
		g.Title = sArray[1]
	} else if sArray[0] == "Author" {
		g.Author = sArray[1]
	}
}

func parseH(line string) Line {
	l := Line{}
	jCount := 0
	for _, j := range line {
		if j == '#' {
			jCount++
		}
		if jCount > 2 {
			break
		}
	}
	switch jCount {
	case 1:
		l.Type = GEMTEXT_H1
	case 2:
		l.Type = GEMTEXT_H2
	case 3:
		l.Type = GEMTEXT_H3
	}
	l.Text = strings.TrimSpace(line[jCount:])
	return l
}

func parseLink(line string) Line {
	l := Line{Type: GEMTEXT_LINK}
	p := strings.Split(strings.TrimSpace(line[2:]), " ")
	if len(p) > 1 {
		l.Arg = strings.Join(p[1:], " ")
		l.Text = p[0]
	} else {
		l.Arg = p[0]
		l.Text = p[0]
	}
	return l
}

func parse(line, char string) Line {
	l := Line{}
	if char == "*" {
		l.Type = GEMTEXT_LIST
	} else {
		l.Type = GEMTEXT_QUOTE
	}
	l.Text = strings.TrimSpace(line[1:])
	return l
}
